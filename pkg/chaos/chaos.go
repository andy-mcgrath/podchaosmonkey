package chaos

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"podchaosmonkey/pkg/config"
	"regexp"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// MurderLoop - main pod murder loop
func MurderLoop(ctx context.Context, cfg config.IConfig, client kubernetes.Interface) {
	log.Println("PodChaosMonkey starting...")
	go func() {
		for {
			log.Printf("...next pod murder in %d seconds...", cfg.GetKillTimeDelay()/1000000000)
			time.Sleep(cfg.GetKillTimeDelay())

			podlist, err := getPodList(ctx, cfg, client)
			if err != nil {
				log.Fatalf("Error getting pod list: %s", err)
			}

			podlistFiltered, err := filterPodList(cfg, podlist)
			if err != nil {
				log.Fatalf("Error filtering pod list: %s", err)
			}

			if err := murderPod(ctx, cfg, client, podlistFiltered); err != nil {
				log.Fatalf("Error killing pod: %s", err)
			}
		}
	}()

	exitChannel := waitToExit()
	<-exitChannel
}

// Retruns a list of Kubernetes pods from a single namespace, set in configuration
func getPodList(ctx context.Context, cfg config.IConfig, client kubernetes.Interface) (*v1.PodList, error) {
	ctxTimeout, ctxCancel := context.WithTimeout(ctx, cfg.GetConnectionTimeout())
	defer ctxCancel()

	return client.CoreV1().Pods(cfg.GetNamespace()).List(ctxTimeout, metav1.ListOptions{})
}

// Takes a list of pods and randomly selects one to delete aka `murder`
func murderPod(ctx context.Context, cfg config.IConfig, client kubernetes.Interface, podlist *v1.PodList) error {
	if len(podlist.Items) == 0 {
		log.Printf("No pods to murder")
		return nil
	}
	ctxTimeout, ctxCancel := context.WithTimeout(ctx, cfg.GetConnectionTimeout())
	defer ctxCancel()

	randPodIndex := rand.Intn(len(podlist.Items))
	pod := &podlist.Items[randPodIndex]

	log.Printf("... !RedRuM pod '%s' pod MuRdeR!\n", pod.Name)
	return client.CoreV1().Pods(cfg.GetNamespace()).Delete(ctxTimeout, pod.Name, metav1.DeleteOptions{})
}

// Returns a Kubernetes core.v1.PodList with only pod whos Pod.Name matches the configured regex filter
func filterPodList(cfg config.IConfig, podlist *v1.PodList) (*v1.PodList, error) {
	filteredPods := []v1.Pod{}
	r, err := regexp.Compile(cfg.GetPodFilter())
	if err != nil {
		return nil, err
	}

	for _, pod := range podlist.Items {
		if r.MatchString(pod.Name) {
			filteredPods = append(filteredPods, pod)
		}
	}

	return &v1.PodList{
		TypeMeta: podlist.TypeMeta,
		ListMeta: podlist.ListMeta,
		Items:    filteredPods,
	}, nil
}

// WaitToExit - Returns a blocking channel which unblocks and closes on an OS Interrupt signal
func waitToExit() <-chan struct{} {
	runC := make(chan struct{}, 1)

	sc := make(chan os.Signal, 1)

	signal.Notify(sc, os.Interrupt)

	go func() {
		defer close(runC)
		<-sc
		log.Println("Exiting...")
		time.Sleep(500 * time.Millisecond)
	}()

	return runC
}
