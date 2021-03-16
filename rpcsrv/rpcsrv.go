package rpcsrv

import (
	"fmt"
	"log"
	"net"
	"net/http"
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

	fmt.Println("crictl reply result was:", string((*reply)[:]))
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

	fmt.Println("GetContainers crictl reply result was:", string((*reply)[:]))
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
	rpc.HandleHTTP()

	l, e := net.Listen("unix", sock)
	if e != nil {
		log.Fatal("Error starting listener:", e)
	}

	//fmt.Println("Starting container info server with address:", sock)
	http.Serve(l, nil)
}
