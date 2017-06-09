package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rai-project/config"
)

func usingSTS(opts *SessionOptions, roleSessionName, account, role string) error {
	conf := &aws.Config{
		Region: aws.String(opts.Region),
	}
	if config.IsVerbose {
		conf = conf.WithCredentialsChainVerboseErrors(true)
	}

	conf = conf.WithCredentials(credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     opts.AccessKey,
			SecretAccessKey: opts.SecretKey,
		}},
	}))

	sess := session.New()

	svc := sts.New(sess, conf)

	rolearn := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, role)
	output, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(rolearn),
		RoleSessionName: aws.String(roleSessionName),
		DurationSeconds: aws.Int64(int64(opts.stsRoleDurationSeconds.Seconds())),
	})
	if err != nil {
		log.Errorf("Unable to assume role: %v", err.Error())
		return err
	}

	accessKey := aws.StringValue(output.Credentials.AccessKeyId)
	secretKey := aws.StringValue(output.Credentials.SecretAccessKey)
	sessionToken := aws.StringValue(output.Credentials.SessionToken)

	os.Setenv("AWS_ACCESS_KEY_ID", accessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)
	os.Setenv("AWS_SESSION_TOKEN", sessionToken)

	opts.AccessKey = accessKey
	opts.SecretKey = secretKey
	opts.SessionToken = sessionToken

	return nil
}
