# Multi-IPSIX: Advanced IPv6 Address Management for Linux

![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-Linux-lightgrey.svg)

## Overview

**Multi-IPSIX** is a powerful, configuration-driven tool for managing multiple IPv6 addresses on Linux network interfaces. It is designed for advanced networking scenarios where precise control over source address selection is critical. The core feature of Multi-IPSIX is its ability to assign a "priority" IPv6 address to an interface while gracefully "deprecating" all other IPv6 addresses on that same interface. 

This ensures that only the priority IP is used for initiating new outbound connections (e.g., `curl`, `ssh`, `git`), while the deprecated addresses can still receive incoming traffic. This project was born from the need to control outgoing traffic source IPs in multi-homed or complex IPv6 environments.

## Key Features

- **Declarative Configuration**: Define all your interfaces, IPv6 addresses, and priority IPs in a simple, human-readable YAML file (`atpajah.yaml`).
- **Priority IP**: Assign a single, stable source IP for all outgoing traffic from an interface.
- **Graceful Deprecation**: Non-priority IPs are deprecated by setting their `preferred_lft` to `0`. They are not removed and can still be used for incoming connections.
- **Idempotent**: The tool can be run safely multiple times. It checks for existing IPs and avoids errors.
- **Extensible**: Built with a clean, modular structure in Go, making it easy to extend and maintain.

## How It Works

The tool operates in a straightforward sequence:

1.  **Load Configuration**: It reads the `atpajah.yaml` file to get the list of interfaces and their desired IP configurations.
2.  **Ensure Addresses**: For each interface, it iterates through the list of IPv6 addresses and ensures they are all present using the `ip addr add` command. If an address already exists, it continues without error.
3.  **Deprecate Non-Priority IPs**: It then re-iterates through the list. If an IP address does **not** match the designated `priority_ip` for that interface, it uses the `ip addr change` command to set its `preferred_lft` (preferred lifetime) to `0`. 

This `preferred_lft 0` setting is the key. The Linux kernel will not select a deprecated address as the source for new outgoing connections unless an application explicitly binds to it.

## Getting Started

### Prerequisites

- **Go**: Version 1.18 or higher.
- **Linux**: A Linux distribution with the `iproute2` toolset (which provides the `ip` command).
- **Root Privileges**: The tool must be run as `root` or with `sudo` as it modifies network interface configurations.

### Installation & Usage

1.  **Clone the repository (or download the files):**
    ```bash
    git clone https://github.com/your-username/Multi-IPSIX.git
    cd Multi-IPSIX
    ```

2.  **Configure your interfaces:**
    Open `atpajah.yaml` and customize it to your needs. Add your interface names, the full list of IPv6 addresses (with CIDR masks), and specify which one should be the `priority_ip`.

    **Example `atpajah.yaml`:**
    ```yaml
    interfaces:
      - name: "eth0"
        priority_ip: "2001:db8:1::1"
        addresses:
          - "2001:db8:1::1/64"
          - "2001:db8:1::2/64"
          - "2001:db8:1::3/64"
      - name: "wireguard-vpn"
        priority_ip: "2001:db8:2::aaaa"
        addresses:
          - "2001:db8:2::aaaa/128"
          - "2001:db8:2::bbbb/128"
    ```

3.  **Run the application:**
    Execute the program from the root of the project directory with `sudo`.

    ```bash
    sudo go run cmd/multi-ipsix/main.go
    ```

    The tool will print the actions it is taking for each interface.

## Building a Binary

For production use, you can compile the application into a single binary:

```bash
# This will create a binary named 'multi-ipsix' in the project root
go build -o multi-ipsix cmd/multi-ipsix/main.go
```

Then you can run the binary directly:

```bash
sudo ./multi-ipsix
```

## Future & Contribution

This tool is a foundational step towards more dynamic and intelligent IPv6 management. We welcome contributions! Feel free to open an issue or submit a pull request.

**Disclaimer**: Network configuration is sensitive. Always test in a safe environment before deploying to production systems.
