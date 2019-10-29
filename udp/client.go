package udp

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Ctx     context.Context
	address string
}

var onceClient sync.Once
var InstanceClient *Client

type Query interface {
	Request(path string, body string) (error, string)
}

func NewClient(ctx context.Context, address string) *Client {
	onceClient.Do(func() {
		InstanceClient = new(Client)

		if ctx == nil {
			ctx = context.Background()
		}

		InstanceClient.Ctx = ctx
		InstanceClient.address = address
	})

	return InstanceClient
}

func (client Client) Request(path string, body string) (err error, result string) {

	reader := strings.NewReader(body)
	raddr, err := net.ResolveUDPAddr(udp, client.address)
	if err != nil {
		return
	}

	conn, err := net.DialUDP(udp, nil, raddr)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)
	resultChan := make(chan string, 1)

	go func() {

		n, err := io.Copy(conn, reader)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("Client packet-written: bytes=%d\n", n)

		buffer := make([]byte, maxBufferSize)

		timeout := time.Now().Add(time.Duration(30) * time.Second)
		err = conn.SetReadDeadline(timeout)
		if err != nil {
			doneChan <- err
			return
		}

		nRead, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("Client packet-received: bytes=%d from=%s\n", nRead, addr.String())
		resultChan <- string(buffer)
		doneChan <- nil
	}()

	client.handlerClientSignals(doneChan, client.Ctx.Done())
	result = <-resultChan
	return
}

func (client Client) handlerClientSignals(errors <-chan error, cancelled <-chan struct{}) {
	select {
	case <-cancelled:
		fmt.Println("cancelled")
		_ = client.Ctx.Err()

	case err := <-errors:
		if err != nil {
			fmt.Println("error " + err.Error())
		}

	}
}
