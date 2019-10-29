package main

import (
	"context"
	"fmt"
	"golang-udp-server/udp"
	"os"
	"strconv"
)

func main() {
	fmt.Println("process pid " + strconv.Itoa(os.Getpid()))
	ctx := context.Background()
	server := udp.NewServer(ctx, "0.0.0.0:8080")
	go server.Start()

	client := udp.NewClient(ctx, "127.0.0.1:8080")
	_, result := client.Request("/test", "hello world!")

	fmt.Println(result)

	/*	if err := udp.Server(context.Background(), "0.0.0.0:8080"); err != nil {
		fmt.Errorf("UDP server error %s", err)
	}*/

}
