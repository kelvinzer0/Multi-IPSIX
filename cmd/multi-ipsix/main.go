package main

import (
	"fmt"
	"log"
	"sync"

	"Multi-IPSIX/internal/config"
	"Multi-IPSIX/internal/ipmanager"
	"Multi-IPSIX/internal/ipv6monitor"
)

func main() {
	// Define the path to the configuration file.
	configFile := "/etc/multi-ipsix/atpajah.yaml"

	// Load configuration from the YAML file.
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading configuration from %s: %v", configFile, err)
	}

	fmt.Println("Successfully loaded configuration. Starting IP management...")

	var wg sync.WaitGroup

	// Process each interface defined in the configuration.
	for _, ifaceConfig := range cfg.Interfaces {
		fmt.Printf("--- Processing interface: %s ---\n", ifaceConfig.Name)

		// 1. Add all configured addresses to ensure they are present on the interface.
		for _, ipAddr := range ifaceConfig.Addresses {
			if err := ipmanager.AddIPv6Address(ifaceConfig.Name, ipAddr); err != nil {
				// Log non-fatal errors and continue.
				log.Printf("Warning on adding IP: %v\n", err)
			}
		}

		// 2. Deprecate all IPv6 addresses on the interface, except the priority IP.
        if err := ipmanager.DeprecateAllIPv6Addresses(ifaceConfig.Name, ifaceConfig.PriorityIP); err != nil {
            log.Printf("Warning on deprecating all IPv6 addresses: %v\n", err)
        }

		// 3. Ensure the priority IP is undeprecated.
		if ifaceConfig.PriorityIP != "" {
			if err := ipmanager.UndeprecateIPv6Address(ifaceConfig.Name, ifaceConfig.PriorityIP); err != nil {
				log.Printf("Warning on undeprecating priority IP %s: %v\n", ifaceConfig.PriorityIP, err)
			}
		}
		fmt.Printf("--- Finished processing interface: %s ---\n\n", ifaceConfig.Name)

		// Start IPv6 monitoring if enabled for this interface
		if ifaceConfig.MonitorIPv6 {
			wg.Add(1)
			go func(cfg config.InterfaceConfig) {
				defer wg.Done()
				ipv6monitor.MonitorInterface(cfg)
			}(ifaceConfig)
		}
	}

	fmt.Println("Multi-IPSIX configuration applied successfully. Starting IPv6 monitoring (if enabled)...")

	// Keep the main goroutine alive indefinitely to allow background goroutines to run.
	wg.Wait() // Wait for all monitoring goroutines to finish (they won't, so this keeps main alive)
	select {} // Block forever
}
