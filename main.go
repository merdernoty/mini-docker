// +build linux
package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <command> [args...]", os.Args[0])
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func run() {

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running the command: %v", err)
	}
}

func child() {
	log.Printf("Running %v as PID %d", os.Args[2:], os.Getpid())
	if err := syscall.Sethostname([]byte("container")); err != nil {
		log.Fatalf("Error setting hostname: %v", err)
	}

	if err := syscall.Chroot("/home/rafael/container"); err != nil {
		log.Fatalf("Error changing root: %v", err)
	}

	if err := os.Chdir("/"); err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running the command: %v", err)
	}
}
