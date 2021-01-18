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
	fmt.Println("container info server v0.0.6.")

	go rpcsrv.RPCSrv(config.Sock)

	// Give the server time to start locally before launching chroot functions.
	time.Sleep(5 * time.Second)

	// Wait for RPC calls, gather info about containers, and return it to the caller.
	chroot.SysCmd(models.ChrootChan, models.RuncChan, models.ContainersChan)
}
