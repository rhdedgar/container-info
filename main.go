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
	fmt.Println("container info server v0.0.16.")

	go rpcsrv.RPCSrv(config.Sock)

	// Give the server time to start locally before launching chroot functions.
	time.Sleep(5 * time.Second)

	// Wait for RPC calls, gather info about containers, and return it to the caller.
	chroot.SysCmd(models.ChrootChan, models.RuncChan, models.ContainersChan)
}
