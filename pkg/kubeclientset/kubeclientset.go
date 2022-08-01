package kubeclientset

import (
	"flag"
	"path/filepath"
	"podchaosmonkey/pkg/config"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Creates and returns Kubernetes clientset from local kubeconfig file, for an off-cluster clentset
func kubeClientset() (clientset *kubernetes.Clientset, err error) {
	k := kubeconfigPath()
	config, err := clientcmd.BuildConfigFromFlags("", *k)
	if err != nil {
		return
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	return
}

// Returns the path to local kubeconfig file
func kubeconfigPath() (kp *string) {
	home := homedir.HomeDir()
	if home != "" {
		home = filepath.Join(home, ".kube", "config")
	}
	kp = flag.String("kubeconfig", home, "absolute path to the kubeconfig file")
	flag.Parse()
	return
}

// NewClientSet returns a Kubernetes clientset, environment retrieved from configuration `cfg`
// - environment == `DEV` return kubeconfig file clientset
// - environment != `DEV` return in cluster config clientset (assumes running in a cluster)
func NewClientSet(cfg config.IConfig) (clientset *kubernetes.Clientset, err error) {
	if cfg.GetEnvironment() == "DEV" {
		return kubeClientset()
	}

	conf, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(conf)
}
