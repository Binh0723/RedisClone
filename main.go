package main

import (
	"fmt"
	"net"
	"strings"
)

func main(){
	fmt.Println("Hello World")
	l, err := net.Listen("tcp", ":6379")
	if err != nil{
		fmt.Println("Error listening:", err)
		return
	}
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	defer l.Close()
	
	fmt.Println("Waiting for connections...")
	connection, err := l.Accept()
	if err != nil{
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer connection.Close() // close connection once finished
	for{
		resp := NewResp(connection)
		value,err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		
		fmt.Println(value)
		writer := NewWriter(connection)
		handler,ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}
		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}	
}