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

package rpcsrv

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/rhdedgar/container-info/channels"
	"github.com/rhdedgar/container-info/models"
)

// InfoSrv is the base type that needs to be exported for RPC to work.
type InfoSrv struct {
}

// GetContainerInfo is the RPC-exported method that returns docker or crictl info about a container.
func (g InfoSrv) GetContainerInfo(containerID *string, reply *[]byte) error {
	//fmt.Println("Getting info for container: ", *containerID)

	channels.SetStringChan(models.ChrootChan, *containerID)
	*reply = <-models.ChrootOut

	//fmt.Println("crictl reply result was:", string((*reply)[:]))
	resultSample := string((*reply)[:])
	if len(resultSample) > 32 {
		fmt.Println(resultSample[:32])
	}
	return nil
}

// GetRuncInfo is the RPC-exported method that returns runc info about a container.
func (g InfoSrv) GetRuncInfo(containerID *string, reply *[]byte) error {
	//fmt.Println("Getting runc info for container: ", *containerID)

	channels.SetStringChan(models.RuncChan, *containerID)
	*reply = <-models.RuncOut

	fmt.Println("runc reply result was:", string((*reply)[:]))
	return nil
}

// GetContainers is the RPC-exported method that returns a list of running containers from crictl.
func (g InfoSrv) GetContainers(minContainerAge string, reply *[]byte) error {
	//fmt.Println("Getting running container info)

	channels.SetStringChan(models.ContainersChan, minContainerAge)
	*reply = <-models.ContainersOut

	//fmt.Println("GetContainers crictl reply result was:", string((*reply)[:]))
	return nil
}

// RPCSrv listens for container UIDs and returns info about that container.
func RPCSrv(sock string) {
	// Start by cleaning up any leftover sockets of the same name.
	if err := os.RemoveAll(sock); err != nil {
		log.Fatal(err)
	}

	InfoSrv := new(InfoSrv)

	rpc.Register(InfoSrv)
	//rpc.HandleHTTP()

	l, e := net.Listen("unix", sock)
	if e != nil {
		log.Fatal("Error starting listener:", e)
	}

	//fmt.Println("Starting container info server with address:", sock)
	//http.Serve(l, nil)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
