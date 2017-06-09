package aws

import (
	_ "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AWSTestSuite struct {
	suite.Suite
}

func (suite *AWSTestSuite) SetupTest() {
}

func (suite *AWSTestSuite) TestLoad() {
	assert.NotNil(suite.T(), Config)
}

func (suite *AWSTestSuite) TestPrintable() {
	assert.NotEqual(suite.T(), "", Config.String())
}

func (suite *AWSTestSuite) TestRegion() {
	assert.Equal(suite.T(), AWSRegionUSEast1, Config.Region)
}

func TestAWSConfig(t *testing.T) {
	suite.Run(t, new(AWSTestSuite))
}
