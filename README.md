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
- **Systemd Integration**: Can be installed as a `systemd` service to run automatically on boot.

## Getting Started

### Prerequisites

- **Go**: Version 1.18 or higher (for building from source).
- **Linux**: A Linux distribution with `systemd` and the `iproute2` toolset.
- **Root Privileges**: The installation script must be run with `sudo`.

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/Multi-IPSIX.git
    cd Multi-IPSIX
    ```

2.  **Customize the configuration:**
    Before installing, open `atpajah.yaml` and configure your interfaces, addresses, and priority IPs according to your needs.

3.  **Run the installer:**
    Execute the installation script with `sudo`. It will build the binary, move it to `/usr/local/bin`, copy the configuration, and set up the `systemd` service.
    ```bash
    sudo bash install.sh
    ```

The service is now enabled and will start automatically on boot. You can also start it manually:
```bash
sudo systemctl start multi-ipsix
```

### Uninstallation

To remove the application and all its components, run the uninstallation script:
```bash
sudo bash uninstall.sh
```
This will stop the service, remove the binary, and delete the `systemd` service file. It will ask for confirmation before deleting the configuration directory `/etc/multi-ipsix`.

## How It Works

The tool operates in a straightforward sequence:

1.  **Load Configuration**: It reads `/etc/multi-ipsix/atpajah.yaml` to get the list of interfaces and their desired IP configurations.
2.  **Ensure Addresses**: For each interface, it iterates through the list of IPv6 addresses and ensures they are all present using the `ip addr add` command. If an address already exists, it continues without error.
3.  **Deprecate Non-Priority IPs**: It then re-iterates through the list. If an IP address does **not** match the designated `priority_ip` for that interface, it uses the `ip addr change` command to set its `preferred_lft` (preferred lifetime) to `0`. 

This `preferred_lft 0` setting is the key. The Linux kernel will not select a deprecated address as the source for new outgoing connections unless an application explicitly binds to it.

## Future & Contribution

This tool is a foundational step towards more dynamic and intelligent IPv6 management. We welcome contributions! Feel free to open an issue or submit a pull request.

**Disclaimer**: Network configuration is sensitive. Always test in a safe environment before deploying to production systems.