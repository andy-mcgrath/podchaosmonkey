package chaos

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
)

type (
  mockConfig struct {}
)

func (c *mockConfig) GetEnvironment() string  { return "DEV" }
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

func creatPods(c kubernetes.Interface, numOfPods int) error {
	for i := 0; i < numOfPods; i++ {
		pod := &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("test-pod-%d", i),
				Namespace: "test",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:            fmt.Sprintf("nginx-%d", i),
						Image:           "nginx",
						ImagePullPolicy: "Always",
					},
				},
			},
		}
		_, err := c.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {	return err }
	}
	return nil
}

func TestGetPodList(t *testing.T) {
	c := testclient.NewSimpleClientset()
	if err := creatPods(c, 3); err != nil {
		t.Errorf("Error creating test pods: %s", err)
	}
	p, err := getPodList(context.TODO(), &mockConfig{}, c)
	if err != nil {
		t.Errorf("Error getting pod list: %s", err)
	}
	assert.Len(t, p.Items, 3, "Pod list should have 3 pods")
}


func TestFilterPodList(t *testing.T) {
	c := testclient.NewSimpleClientset()
	if err := creatPods(c, 3); err != nil {
		t.Errorf("Error creating test pods: %s", err)
	}
	p, err := getPodList(context.TODO(), &mockConfig{}, c)
	if err != nil {
		t.Errorf("Error getting pod list: %s", err)
	}
	pf, err := filterPodList(&mockConfig{}, p)
	if err != nil {
		t.Errorf("Error getting pod list: %s", err)
	}
	assert.Len(t, pf.Items, 3, "Filtered Pod list should have 3 pods")
}

func TestMurderPod(t *testing.T) {
	c := testclient.NewSimpleClientset()
	if err := creatPods(c, 3); err != nil {
		t.Errorf("Error creating test pods: %s", err)
	}
	p, err := getPodList(context.TODO(), &mockConfig{}, c)
	if err != nil {
		t.Errorf("Error getting pod list: %s", err)
	}
	if err := murderPod(context.TODO(), &mockConfig{}, c, p); err != nil {
		t.Errorf("Error killing pods: %s", err)
	}
	p, err = getPodList(context.TODO(), &mockConfig{}, c)
	if err != nil {
		t.Errorf("Error getting pod list: %s", err)
	}
	assert.Len(t, p.Items, 2, "Filtered Pod list should have 2 pods")
}