// +build linux

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	containerRoot = "/home/rafael/container"
	containerHostname = "container"
)

type ContainerConfig struct {
	command string
	args    []string
	stdin   *os.File
	stdout  *os.File
	stderr  *os.File
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: %s <command> [args...]", os.Args[0])
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	switch os.Args[1] {
	case "run":
		return handleRun(ctx)
	case "child":
		return handleChild(ctx)
	default:
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func handleRun(ctx context.Context) error {
	config := &ContainerConfig{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	cmd := exec.CommandContext(ctx, "/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = config.stdin
	cmd.Stdout = config.stdout
	cmd.Stderr = config.stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	return cmd.Run()
}

func handleChild(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(3)
	go setupHostname(&wg, errChan)
	go setupRootFS(&wg, errChan)
	go setupWorkDir(&wg, errChan)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("container setup failed: %v", err)
		}
	}

	log.Printf("Running %v as PID %d", os.Args[2:], os.Getpid())
	
	return executeCommand(ctx)
}

func setupHostname(wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	if err := syscall.Sethostname([]byte(containerHostname)); err != nil {
		errChan <- fmt.Errorf("failed to set hostname: %v", err)
	}
}

func setupRootFS(wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	if err := syscall.Chroot(containerRoot); err != nil {
		errChan <- fmt.Errorf("failed to change root: %v", err)
	}
}

func setupWorkDir(wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	if err := os.Chdir("/"); err != nil {
		errChan <- fmt.Errorf("failed to change directory: %v", err)
	}
}

func executeCommand(ctx context.Context) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("no command specified")
	}

	cmd := exec.CommandContext(ctx, os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
