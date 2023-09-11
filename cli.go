package main

import (
	"building-microservices-with-go/client"
	"building-microservices-with-go/server"
	"fmt"
)

func main() {
	go server.StartServer()

	c := client.CreateClient()
	defer c.Close()

	reply := client.PerformRequest(c)
	fmt.Println(reply.Message)
}
