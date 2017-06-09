package vipertags

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testYamlFile = "config.yaml"

var testJsonFile = "config.json"

var testYaml = `test:
    hostname: "127.0.0.1"
    port: 3126
    overridden: ""
    duration: "1h"
    aslice: ["elem0", "elem1"]
`

func TestMain(m *testing.M) {
	os.Setenv("CONF_TEST_OVERRIDDEN", "test")
	bytes := []byte(testYaml)
	err := ioutil.WriteFile(testYamlFile, bytes, 0777)
	if err != nil {
		panic("cannot create temp yaml config file needed for testing")
	}
	defer os.Remove(testYamlFile)
	FromFile(testYamlFile, "")
	m.Run()
}

func ExampleGetConfig() {
	// ./config.yaml
	// test:
	//    hostname: "127.0.0.1"
	//    port: 3126
	type StringConfig struct {
		SomeHostname string `config:"test.hostname"`
		Port         int    `config:"test.port"`
		Overriden    string `config:"test.overridden"`
	}
	type InvalidConfig struct {
		Invalid string `config:"some.invalid.config"`
	}
	type InvalidMapConfig struct{}
	c := StringConfig{}
	Setup("yaml", "CONF") // Or Setup("json")
	Fill(&c)
	fmt.Println(c.SomeHostname)
	fmt.Println(c.Port)
	fmt.Println(c.Overriden)
	i := InvalidConfig{}
	Fill(&i)
	fmt.Printf(i.Invalid)
	// Output:
	// 127.0.0.1
	// 3126
	// test
}

func TestDefaults(t *testing.T) {
	type StringConfig struct {
		FromEnv1 string `config:"name" default:"foo"`
	}
	c := StringConfig{}
	Setup("yaml", "CONF") // Or Setup("json")
	SetDefault(&c)
	assert.Equal(t, "foo", c.FromEnv1, "")
}

func TestEnvironment(t *testing.T) {
	os.Setenv("CONF_TEST1_FROMENV", "foo")
	os.Setenv("CONF_FROMENV", "bar")
	type StringConfig struct {
		FromEnv1 string `config:"test1.fromenv"`
		FromEnv2 string `config:"fromenv"`
	}
	c := StringConfig{}
	Setup("yaml", "CONF") // Or Setup("json")
	Fill(&c)
	assert.Equal(t, c.FromEnv1, "foo", "")
	assert.Equal(t, c.FromEnv2, "bar", "")
}

func ExampleDuration() {
	type DurConfig struct {
		Duration time.Duration `config:"test.duration"`
	}
	c := DurConfig{}
	Setup("yaml", "CONF")
	Fill(&c)
	fmt.Println(c.Duration)

	// Output:
	// 1h0m0s
}

func TestSlice(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	type SliceConfig struct {
		Slice []string `config:"test.aslice"`
	}
	c := SliceConfig{}
	Setup("yaml", "CONF") // Or Setup("json")
	SetDefaults(&c)
	Fill(&c)
	assert.Equal(t, []string{"elem0", "elem1"}, c.Slice, "")
}

func TestDefaultSlice(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	type SliceConfig struct {
		Slice []string `default:"default value"`
	}
	c := SliceConfig{}
	//Setup("yaml", "CONF") // Or Setup("json")
	SetDefaults(&c)
	assert.Equal(t, []string{"default", "value"}, c.Slice, "")
}

// func TestMemoryLimit(t *testing.T) {

// }
