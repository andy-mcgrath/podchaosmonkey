package main

import (
	"context"
	"log"
	"podchaosmonkey/pkg/chaos"
	"podchaosmonkey/pkg/config"
	"podchaosmonkey/pkg/kubeclientset"
)

func main() {
	// Load config
	log.Println("Initialising Configuration...")
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	// Create Kubernetes clientset
	log.Println("Creating Kubernetes ClientSet...")
	client, err := kubeclientset.NewClientSet(cfg)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err)
	}

	ctx := context.Background()

	// Start main event loop
	chaos.MurderLoop(ctx, cfg, client)
}


