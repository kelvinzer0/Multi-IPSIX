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

// DeprecateIPv6Address sets the preferred lifetime of a given IPv6 address to 0.
// It correctly handles addresses with or without CIDR notation.
func DeprecateIPv6Address(iface, ipAddr string) error {
	var ip net.IP
	// Handle addresses with or without CIDR notation.
	if strings.Contains(ipAddr, "/") {
		parsedIP, _, err := net.ParseCIDR(ipAddr)
		if err != nil {
			return fmt.Errorf("error parsing CIDR %s: %v", ipAddr, err)
		}
		ip = parsedIP
	} else {
		ip = net.ParseIP(ipAddr)
		if ip == nil {
			return fmt.Errorf("invalid IP address format: %s", ipAddr)
		}
	}

	cmd := exec.Command("ip", "addr", "change", ip.String(), "dev", iface, "preferred_lft", "0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error deprecating IP %s on %s: %v, output: %s", ipAddr, iface, err, string(output))
	}
	fmt.Printf("Deprecated IP %s on interface %s\n", ipAddr, iface)
	return nil
}