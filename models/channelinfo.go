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

package models

var (
	// ChrootChan is the input channel for containerIDs
	ChrootChan = make(chan string)
	// ChrootOut is the output channel for crictl commands run by chroot.SysCmd
	ChrootOut = make(chan []byte)
	// RuncChan is the input channel for containerIDs used in runc inspection commands
	RuncChan = make(chan string)
	// RuncOut is the output channel for runc commands run by chroot.SysCmd
	RuncOut = make(chan []byte)
	// ContainersChan is the input channel for containerIDs
	ContainersChan = make(chan string)
	// ContainersOut is the output channel for crictl commands run by chroot.SysCmd
	ContainersOut = make(chan []byte)
)
