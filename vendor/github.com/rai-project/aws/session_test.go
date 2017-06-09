package aws

import (
	_ "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SessionTestSuite struct {
	suite.Suite
}

func NewSessionTestSuite() *SessionTestSuite {
	return &SessionTestSuite{}
}

func (suite *SessionTestSuite) TestSts() {
	session, err := NewSession(
		Region(AWSRegionUSEast1),
		Sts(),
	)
	assert.NoError(suite.T(), err, "must be able to connect")
	assert.NotNil(suite.T(), session, "sesion must not be nil")
}

func TestSessionConfig(t *testing.T) {
	suite.Run(t, NewSessionTestSuite())
}
