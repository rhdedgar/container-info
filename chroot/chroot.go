/*
Copyright 2020 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chroot

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/rhdedgar/container-info/models"
)

var (
	// Path is the path to the container runtime interface utility
	Path = "/usr/bin/crictl"
)

func inspectContainer(containerID string) ([]byte, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	fmt.Println("running this: ", Path+" inspect "+containerID)
	cmd := exec.Command(Path, "inspect", containerID)

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	for i := 0; i <= 5; i++ {
		if i >= 5 {
			return out.Bytes(), fmt.Errorf("inspectContainer: Error inspecting container after multiple retries: %v\n", stderr.String())
		}

		time.Sleep(time.Duration(i) * time.Second)

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running inspect command: ", err)
			fmt.Println("Stderr output was: ", stderr.String())
			fmt.Println("inspectContainer: Error inspecting container, waiting a few seconds in case it just isn't available yet.")
			continue
		}

		sStr := out.String()
		sStderr := stderr.String()

		//fmt.Println("Command output was", sStr)
		fmt.Println("Command stderr output was", sStderr)

		if sStderr != "" {
			fmt.Println("stderr was returned from the command: ", sStderr)
			fmt.Println("inspectContainer: Error inspecting container, waiting a few seconds in case it just isn't available yet.")
			continue
		}

		if sStr != "" {
			fmt.Println("Found container information; returning.")
			return out.Bytes(), nil
		}
		fmt.Println("inspectContainer output is empty.")
	}

	return out.Bytes(), fmt.Errorf("inspectContainer could not return data for this container.\n")
}

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
			result, err := inspectContainer(containerID)
			if err != nil {
				fmt.Println("Error returned from inspectContainer: ", err)
			}
			models.ChrootOut <- result

		case scanContainer := <-runcChan:
			fmt.Println("running runc state command")
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
