package ipmanager

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// AddIPv6Address adds a given IPv6 address to a specified network interface.
// It ensures the address has a /64 CIDR suffix if it's missing.
func AddIPv6Address(iface, ipAddr string) error {
	// Ensure the address has a CIDR suffix, default to /64 for IPv6 if missing.
	addrToAdd := ipAddr
	if !strings.Contains(addrToAdd, "/") {
		addrToAdd = fmt.Sprintf("%s/64", ipAddr) // Append default IPv6 prefix
	}

	cmd := exec.Command("ip", "addr", "add", addrToAdd, "dev", iface)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Ignore "file exists" errors, which are expected if the IP is already configured.
		if !strings.Contains(string(output), "File exists") {
			return fmt.Errorf("error adding IP %s to %s: %v, output: %s", addrToAdd, iface, err, string(output))
		}
	}
	fmt.Printf("Ensured IP %s exists on interface %s\n", addrToAdd, iface)
	return nil
}

// DeprecateAllIPv6Addresses sets the preferred lifetime of all IPv6 addresses on a given interface to 0.
func DeprecateAllIPv6Addresses(iface string) error {
	// Get all IPv6 addresses on the interface
	cmd := exec.Command("ip", "-6", "addr", "show", "dev", iface)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error listing IPv6 addresses on %s: %v, output: %s", iface, err, string(output))
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "inet6") {
			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) > 1 {
				ipWithCIDR := fields[1]
				ip, _, err := net.ParseCIDR(ipWithCIDR)
				if err != nil {
					fmt.Printf("Warning: Could not parse IP %s: %v\n", ipWithCIDR, err)
					continue
				}

				// Deprecate each IPv6 address
				deprecateCmd := exec.Command("ip", "addr", "change", ip.String(), "dev", iface, "preferred_lft", "0")
				deprecateOutput, deprecateErr := deprecateCmd.CombinedOutput()
				if deprecateErr != nil {
					fmt.Printf("Warning: Error deprecating IP %s on %s: %v, output: %s\n", ip.String(), iface, deprecateErr, string(deprecateOutput))
				} else {
					fmt.Printf("Deprecated IP %s on interface %s\n", ip.String(), iface)
				}
			}
		}
	}
	return nil
}