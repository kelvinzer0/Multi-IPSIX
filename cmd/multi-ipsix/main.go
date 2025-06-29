package main

import (
	"fmt"
	"log"
	"net"

	"Multi-IPSIX/internal/config"
	"Multi-IPSIX/internal/ipmanager"
)

func main() {
	// Define the path to the configuration file.
	configFile := "atpajah.yaml"

	// Load configuration from the YAML file.
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading configuration from %s: %v", configFile, err)
	}

	fmt.Println("Successfully loaded configuration. Starting IP management...")

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

		// 2. Deprecate non-priority addresses.
		for _, ipAddr := range ifaceConfig.Addresses {
			// Parse the IP address to compare it without the CIDR suffix.
			ip, _, err := net.ParseCIDR(ipAddr)
			if err != nil {
				log.Printf("Warning: Could not parse CIDR %s: %v\n", ipAddr, err)
				continue
			}

			// If the current IP is the priority IP, skip it.
			if ip.String() == ifaceConfig.PriorityIP {
				fmt.Printf("Skipping deprecation for priority IP %s on interface %s\n", ipAddr, ifaceConfig.Name)
				continue
			}

			// Deprecate the address.
			if err := ipmanager.DeprecateIPv6Address(ifaceConfig.Name, ipAddr); err != nil {
				log.Printf("Warning on deprecating IP: %v\n", err)
			}
		}
		fmt.Printf("--- Finished processing interface: %s ---\n\n", ifaceConfig.Name)
	}

	fmt.Println("Multi-IPSIX configuration applied successfully.")
}
