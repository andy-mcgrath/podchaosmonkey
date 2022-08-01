package config

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigGetEnvironment(t *testing.T) {
	c, _ := Init()
	assert.Equal(t, "DEV", c.GetEnvironment(), "GetEnvironment to return `DEV`")
}

func TestConfigGetLogLevel(t *testing.T) {
	c, _ := Init()
	assert.Equal(t, "info", c.GetLogLevel(), "GetAppConfigSubID to return `info`")
}

func TestConfigGetNamespace(t *testing.T) {
  c, _ := Init()
  assert.Equal(t, "chaos", c.GetNamespace(), "GetNamespace to return `chaos`")
}

func TestConfigGetPodFilter(t *testing.T) {
  c, _ := Init()
  _, err := regexp.Compile(c.GetPodFilter())
  if err != nil {
    t.Errorf("GetPodFilter to return a valid regexp") // handle error
  }
  assert.Equal(t, ".*", c.GetPodFilter(), "GetPodFilter to return `.*`")
}

func TestConfigGetKillTimeDelay(t *testing.T) {
  c, _ := Init()
  assert.Equal(t, time.Duration(120)*time.Second, c.GetKillTimeDelay(), "GetKillTimeDelay to return `120 Seconds`")
}

func TestConfGetConnectionTimeout(t *testing.T) {
  c, _ := Init()
  assert.Equal(t, time.Duration(5)*time.Second, c.GetConnectionTimeout(), "GetConnectionTimeout to return `5 Seconds`")
}