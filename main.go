package main

import (
	"fmt"

	"github.com/rhdedgar/container-info/chroot"
	"github.com/rhdedgar/container-info/config"
	"github.com/rhdedgar/container-info/models"
	"github.com/rhdedgar/container-info/rpcsrv"
)

func main() {
	fmt.Println("container info server v0.0.1.")

	// A goroutine to wait for container IDs, gather info about the container, and return it.
	go chroot.SysCmd(models.ChrootChan, models.RuncChan)

	rpcsrv.RPCSrv(config.Sock)
}
