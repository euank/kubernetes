/*
Copyright 2016 The Kubernetes Authors.

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

package rktshim

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/rktlet/rktlet"
	"google.golang.org/grpc"

	runtimeApi "k8s.io/kubernetes/pkg/kubelet/api/v1alpha1/runtime"
	"k8s.io/kubernetes/pkg/util/interrupt"
)

const (
	// defaultEndpoint is the default address of rktshim grpc server socket.
	defaultAddress = "/var/run/rktshim.sock"
	// unixProtocol is the network protocol of unix socket.
	unixProtocol = "unix"
)

// RktletServer is the grpc server of rktshim
type RktletServer struct {
	// addr is the address to serve on.
	addr string
	// service is the rktlet service which implements runtime and image services.
	service rktlet.ContainerAndImageService

	// server is the grpc server.
	server *grpc.Server
}

// NewRktletServer creates the rktlet grpc server.
func NewRktletServer(addr string, s rktlet.ContainerAndImageService) *RktletServer {
	if addr == "" {
		addr = defaultAddress
	}
	return &RktletServer{
		addr:    addr,
		service: s,
	}
}

// Start starts the rktlet grpc server.
func (s *RktletServer) Start() error {
	glog.V(2).Infof("Starting rktlet grpc server")
	// Unlink to cleanup the previous socket file.
	err := syscall.Unlink(s.addr)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to unlink socket file %q: %v", s.addr, err)
	}
	l, err := net.Listen(unixProtocol, s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %q: %v", s.addr, err)
	}
	// Create the grpc server and register runtime and image services.
	s.server = grpc.NewServer()
	runtimeApi.RegisterRuntimeServiceServer(s.server, s.service)
	runtimeApi.RegisterImageServiceServer(s.server, s.service)
	go func() {
		// Use interrupt handler to make sure the server to be stopped properly.
		h := interrupt.New(nil, s.Stop)
		err := h.Run(func() error { return s.server.Serve(l) })
		if err != nil {
			glog.Errorf("Failed to serve connections: %v", err)
		}
	}()
	return nil
}

// Stop stops the rktlet grpc server.
func (s *RktletServer) Stop() {
	glog.V(2).Infof("Stop rktlet server")
	s.server.Stop()
}

// Endpoint returns the rktlet grpc server unix endpoint
func (s *RktletServer) Endpoint() string {
	return s.addr
}
