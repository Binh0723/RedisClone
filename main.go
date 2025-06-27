package main

import (
	"fmt"
	"net"
	"io"
	"os"
)

func main(){
	fmt.Println("Hello World")
	l, err := net.Listen("tcp", ":6379")
	if err != nil{
		fmt.Println("Error listening:", err)
		return
	}
	defer l.Close()
	
	fmt.Println("Waiting for connections...")
	connection, err := l.Accept()
	if err != nil{
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer connection.Close() // close connection once finished
	for{
		buf := make([]byte,1024)

		_,err := connection.Read(buf)
		if err != nil{
			if err == io.EOF{
				break
			}
			fmt.Println("error reading from client: ", err.Error())
        	os.Exit(1)
		}

		connection.Write([]byte("+OK\r\n"))
	}	
}