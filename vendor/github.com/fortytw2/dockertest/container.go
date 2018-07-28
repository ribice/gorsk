package dockertest

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// A Container is a container inside docker
type Container struct {
	Name string
	Args []string
	Addr string
	cmd  *exec.Cmd
}

// Shutdown ends the container
func (c *Container) Shutdown() {
	c.cmd.Process.Signal(syscall.SIGINT)
	//  Wait till the process exits.
	c.cmd.Wait()
}

// RunContainer runs a given docker container and returns a port on which the
// container can be reached
func RunContainer(container string, port string, waitFunc func(addr string) error, args ...string) (*Container, error) {
	free := freePort()
	host := getHost()
	addr := fmt.Sprintf("%s:%d", host, free)
	argsFull := append([]string{"run"}, args...)
	argsFull = append(argsFull, "-p", fmt.Sprintf("%d:%s", free, port), container)
	cmd := exec.Command("docker", argsFull...)
	// run this in the background

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not run container, %s", err)
	}
	for {
		err := waitFunc(addr)
		if err == nil {
			break
		}

		time.Sleep(time.Millisecond * 150)
	}

	return &Container{
		Name: container,
		Addr: addr,
		Args: args,
		cmd:  cmd,
	}, nil
}

func getHost() string {
	out, err := exec.Command("docker-machine", "ip", os.Getenv("DOCKER_MACHINE_NAME")).Output()
	if err == nil {
		return strings.TrimSpace(string(out[:]))
	}
	return "localhost"
}

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}
