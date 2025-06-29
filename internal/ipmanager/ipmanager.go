package ipmanager

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// AddIPv6Address adds a given IPv6 address (with CIDR) to a specified network interface.
// It executes the `ip addr add` command.
// If the IP address already exists, it gracefully ignores the error and prints a confirmation.
func AddIPv6Address(iface, ipAddr string) error {
	cmd := exec.Command("ip", "addr", "add", ipAddr, "dev", iface)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Ignore "file exists" errors, which are expected if the IP is already configured.
		if !strings.Contains(string(output), "File exists") {
			return fmt.Errorf("error adding IP %s to %s: %v, output: %s", ipAddr, iface, err, string(output))
		}
	}
	fmt.Printf("Ensured IP %s exists on interface %s\n", ipAddr, iface)
	return nil
}

// DeprecateIPv6Address sets the preferred lifetime of a given IPv6 address to 0.
// This effectively marks the address as deprecated, preventing it from being used for new outgoing connections.
// The command requires the IP without the CIDR mask, so this function parses it first.
func DeprecateIPv6Address(iface, ipAddr string) error {
	// The 'ip addr change' command requires the IP without the CIDR mask.
	ip, _, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return fmt.Errorf("error parsing CIDR %s: %v", ipAddr, err)
	}

	cmd := exec.Command("ip", "addr", "change", ip.String(), "dev", iface, "preferred_lft", "0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error deprecating IP %s on %s: %v, output: %s", ipAddr, iface, err, string(output))
	}
	fmt.Printf("Deprecated IP %s on interface %s\n", ipAddr, iface)
	return nil
}
