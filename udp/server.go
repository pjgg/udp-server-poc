package udp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	Ctx context.Context
	// TODO: serverError maybe should be a channel
	ServerError   error
	maxBufferSize int
	address       string
	pid           int
}

type Bootstrap interface {
	Start()
	Stop()
	handlerServerSignals()
}

var once sync.Once
var Instance *Server

func NewServer(ctx context.Context, address string) *Server {
	once.Do(func() {
		Instance = new(Server)
		Instance.maxBufferSize = maxBufferSize
		Instance.pid = os.Getpid()

		if ctx == nil {
			ctx = context.Background()
		}

		Instance.Ctx = ctx
		Instance.address = address
	})

	return Instance
}

func (server Server) Stop() {

	if currentProcess, err := os.FindProcess(server.pid); err != nil {
		// this will be handler by handlerServerSignals
		currentProcess.Signal(os.Interrupt)
	} else {
		server.ServerError = err
	}
}

func (server Server) Start() {

	pc, err := net.ListenPacket(udp, server.address)
	if err != nil {
		return
	}

	defer pc.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, server.maxBufferSize)

	doneByExternalSignal := make(chan os.Signal, 1)
	signal.Notify(doneByExternalSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("Server packet-received: bytes=%d from=%s\n", n, addr.String())

			deadline := time.Now().Add(time.Duration(30) * time.Second)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			response := []byte("This is a random response")
			n, err = pc.WriteTo(response, addr)
			if err != nil {
				// this should be a general server error channel
				doneChan <- err
				return
			}

			fmt.Printf("Server packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	server.handlerServerSignals(doneChan, doneByExternalSignal, server.Ctx.Done())
	return
}

func (server Server) handlerServerSignals(errors <-chan error, signals <-chan os.Signal, cancelled <-chan struct{}) {
	select {
	case <-signals:
		fmt.Println("Server Stopped")
		_, cancel := context.WithTimeout(server.Ctx, 5*time.Second)
		defer func() {
			fmt.Println("... gracefully!")
			cancel()
		}()

	case <-cancelled:
		fmt.Println("request cancelled")
		_ = server.Ctx.Err()

	case err := <-errors:
		if err != nil {
			fmt.Println("error " + err.Error())
		}
	}
}
