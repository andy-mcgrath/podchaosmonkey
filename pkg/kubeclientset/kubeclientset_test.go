package kubeclientset

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	mockConfig struct {
		env string
	}
)

func (c *mockConfig) GetEnvironment() string  { return c.env }
func (c *mockConfig) GetLogLevel() string     { return "info" }
func (c *mockConfig) GetLogOutput() io.Writer { return os.Stdout }
func (c *mockConfig) GetNamespace() string    { return "test" }
func (c *mockConfig) GetPodFilter() string    { return ".*" }
func (c *mockConfig) GetKillTimeDelay() time.Duration {
	return time.Duration(20) * time.Second
}
func (c *mockConfig) GetConnectionTimeout() time.Duration {
	return time.Duration(5) * time.Second
}

func TestDevEnvironmentNewClientSet(t *testing.T) {
	c, err := NewClientSet(&mockConfig{
		env: "DEV",
	})
	if err != nil {
		t.Errorf("Error creating Kubernetes client: %s", err)
	}
	assert.NotNil(t, c, "ClientSet should not be nil")
}
