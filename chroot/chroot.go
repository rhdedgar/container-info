package chroot

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/rhdedgar/container-info/models"
)

var (
	// Path is the path to the container runtime interface utility
	Path = "/usr/bin/crictl"
)

// SysCmd waits for a container ID via channel input, and gathers information
// about the container when it receives an ID.
func SysCmd(cmdChan, runcChan, containersChan <-chan string) {
	_, err := ChrootPath(os.Getenv("CHROOT_PATH"))
	if err != nil {
		fmt.Println("Error getting chroot on host due to: ", err)
	}

	for {
		select {
		case containerID := <-cmdChan:
			//fmt.Println("running this: ", Path+" inspect "+containerID)
			cmd := exec.Command(Path, "inspect", containerID)

			var out bytes.Buffer
			cmd.Stdout = &out

			if cErr := cmd.Run(); err != nil {
				fmt.Println("Error running inspect command: ", cErr)
			}

			//sStr := out.String()
			//fmt.Println("Command output was", sStr)
			models.ChrootOut <- out.Bytes()

		case scanContainer := <-runcChan:
			//fmt.Println("running runc state command")
			runCmd := exec.Command("/usr/bin/runc", "state", scanContainer)

			var runOut bytes.Buffer
			runCmd.Stdout = &runOut

			if runcErr := runCmd.Run(); err != nil {
				fmt.Println("Error running state command: ", runcErr)
			}

			//runcStr := runOut.String()
			//fmt.Println("runc state command output was", runcStr)
			models.RuncOut <- runOut.Bytes()

			// Not going to use the minContainerAge value from containersChan in this ver.
			// This may be replaced with a configurable datetime in the future.
			// e.g. get containers older than minContainerAge.
		case <-containersChan:
			conCmd := exec.Command(Path, "ps", "-o", "json")

			var conOut bytes.Buffer
			conCmd.Stdout = &conOut

			if conErr := conCmd.Run(); err != nil {
				fmt.Println("Error listing running containers: ", conErr)
			}
			models.ContainersOut <- conOut.Bytes()
		}
	}
}

// ChrootPath provides chroot access to the mounted host filesystem.
func ChrootPath(chrPath string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		fmt.Println("Error getting root FD:", err)
		return nil, err
	}

	if err := syscall.Chroot(chrPath); err != nil {
		root.Close()
		fmt.Println("Error with chroot syscall:", err)
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			fmt.Println("Error with root Chdir", err)
			return err
		}
		return syscall.Chroot(".")
	}, nil
}

func init() {
	if _, err := os.Stat("/host/usr/bin/crictl"); os.IsNotExist(err) {
		fmt.Println("Cannot find /host/usr/bin/crictl, using /host/usr/bin/docker")
		Path = "/usr/bin/docker"
	}
}
