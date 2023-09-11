package server

import (
	"building-microservices-with-go/contract"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

const port = 1234

func main() {
	log.Printf("Server starting on port %v\n", port)
	StartServer()
}

func StartServer() {
	helloWorld := &HelloWorldHandler{}
	rpc.Register(helloWorld)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to listen on given port: %s", err))
	}

	log.Printf("Server starting on Port %v\n", port)

	http.Serve(l, nil)
}

type HelloWorldHandler struct{}

func (h *HelloWorldHandler) HelloWorld(args *contract.HelloWorldRequest, reply *contract.HelloWorldResponse) error {
	reply.Message = "Hello " + args.Name
	return nil
}
