# mini-docker

mini-docker is a minimal containerization project in Go. It demonstrates basic Docker principles using Linux namespaces and chroot for process isolation.

## Requirements

- **OS:** Linux (or a Linux environment like WSL2, Podman, or a VM)
- **Go:** 1.13+
- **Privileges:** Root (or privileged container mode)

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/mini-docker.git
cd mini-docker
```

# How to Run the Project

Follow these steps to run the project:

## 1. Clone the Repository and Enter the Project Directory

```bash
git clone https://github.com/yourusername/mini-docker.git
cd mini-docker
```
## 2. Prepare a Minimal Filesystem (rootfs)
The project uses chroot for isolation, so you need to create a rootfs directory containing a minimal Linux environment. For example, you can use debootstrap:

```bash
sudo apt update -o Acquire::Check-Valid-Until=false
sudo apt install -y debootstrap
mkdir -p /path/to/rootfs
sudo debootstrap --arch=arm64 focal /path/to/rootfs http://ports.ubuntu.com/ubuntu-ports
```
Then, move (or create a symbolic link to) the generated filesystem into the project directory if the project expects rootfs at the root.

## 3. Run the Project
To run the project directly from the source code, execute:

```bash
go run main.go run /bin/sh
```
This command will create an isolated environment and run the /bin/sh command inside it.

## 4. (Optional) Run in a Container with Podman
If you want to test the project in a Linux environment using Podman, start a privileged container with the project mounted:

```bash
podman run --privileged -v "$(pwd):/app" -it ubuntu:20.04 /bin/bash
```

Inside the container, navigate to the /app directory and run:

```bash
go run main.go run /bin/sh
```
