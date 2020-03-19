package main

import (
	"fmt"
	"time"

	"github.com/rhdedgar/container-info/chroot"
	"github.com/rhdedgar/container-info/config"
	"github.com/rhdedgar/container-info/models"
	"github.com/rhdedgar/container-info/rpcsrv"
)

func main() {
	fmt.Println("container info server v0.0.4.")

	go rpcsrv.RPCSrv(config.Sock)

	// Give the server time to start locally before launching chroot functions.
	time.Sleep(5 * time.Second)

	// Wait for container IDs, gather info about the container, and return it.
	chroot.SysCmd(models.ChrootChan, models.RuncChan)
}
