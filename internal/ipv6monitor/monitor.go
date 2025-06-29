package ipv6monitor

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"Multi-IPSIX/internal/config"
)

// MonitorInterface memulai pemantauan alamat IPv6 untuk antarmuka yang diberikan.
func MonitorInterface(iface config.InterfaceConfig) {
	log.Printf("Memulai pemantauan IPv6 untuk antarmuka: %s", iface.Name)
	ticker := time.NewTicker(10 * time.Second) // Polling setiap 10 detik
	defer ticker.Stop()

	for range ticker.C {
		currentAddresses, err := getIPv6Addresses(iface.Name)
		if err != nil {
			log.Printf("Gagal mendapatkan alamat IPv6 untuk %s: %v", iface.Name, err)
			continue
		}

		// Konversi allowedAddresses ke map untuk pencarian cepat
		allowedMap := make(map[string]struct{})
		for _, addr := range iface.Addresses {
			allowedMap[normalizeIPv6Address(addr)] = struct{}{}
		}

		for _, currentAddr := range currentAddresses {
			normalizedCurrentAddr := normalizeIPv6Address(currentAddr)
			if _, found := allowedMap[normalizedCurrentAddr]; !found {
				log.Printf("Alamat IPv6 tidak sah terdeteksi pada %s: %s. Menghapus...", iface.Name, currentAddr)
				err := deleteIPv6Address(currentAddr, iface.Name)
				if err != nil {
					log.Printf("Gagal menghapus alamat IPv6 %s dari %s: %v", currentAddr, iface.Name, err)
				} else {
					log.Printf("Berhasil menghapus alamat IPv6 %s dari %s", currentAddr, iface.Name)
				}
			}
		}
	}
}

// getIPv6Addresses mendapatkan daftar alamat IPv6 saat ini untuk antarmuka tertentu.
func getIPv6Addresses(ifaceName string) ([]string, error) {
	cmd := exec.Command("ip", "-6", "addr", "show", "dev", ifaceName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("perintah 'ip' gagal: %v, stderr: %s", err, stderr.String())
	}

	var addresses []string
	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "inet6") && !strings.Contains(line, "scope host") && !strings.Contains(line, "scope link") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "inet6" && i+1 < len(parts) {
					addrWithCIDR := parts[i+1]
					// Hapus bagian CIDR jika ada, karena allowedMap tidak memiliki CIDR
					addr := strings.Split(addrWithCIDR, "/")[0]
					addresses = append(addresses, addr)
					break
				}
			}
		}
	}
	return addresses, nil
}

// deleteIPv6Address menghapus alamat IPv6 dari antarmuka tertentu.
func deleteIPv6Address(address, ifaceName string) error {
	cmd := exec.Command("ip", "-6", "addr", "del", address, "dev", ifaceName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("gagal menghapus alamat IPv6 %s dari %s: %v, stderr: %s", address, ifaceName, err, stderr.String())
	}
	return nil
}

// normalizeIPv6Address menormalkan alamat IPv6 untuk perbandingan.
// Ini akan menghapus bagian CIDR dan memastikan format yang konsisten.
func normalizeIPv6Address(addr string) string {
	// Hapus bagian CIDR jika ada
	parts := strings.Split(addr, "/")
	normalized := parts[0]

	// TODO: Mungkin perlu normalisasi lebih lanjut (misalnya, ekspansi nol, huruf kecil)
	// untuk memastikan perbandingan yang akurat jika alamat dalam config tidak dinormalisasi.
	return normalized
}
