package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/xkenmon/go-tools/net-tools/util"
	"net"
	"os"
)

func main() {
	network := os.Args[1]
	addr := os.Args[2]
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP(network, tcpAddr)
	defer listener.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
  fmt.Println("start listening")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
		}
		go HandlerConn(conn)
	}
}

func HandlerConn(conn *net.TCPConn) {
	fmt.Println("connection established.")
	go util.PrintReceived(conn)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}
		data, err := util.EscapeStr(string(line))
		if err != nil {
			fmt.Println(err)
			continue
		}
		n, err := conn.Write(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n != len(data) {
			fmt.Printf("total of %d bytes of data, but only sent %d bytes\n", len(data), n)
		}
		fmt.Println("--------------------send msg start--------------------")
		fmt.Println(hex.Dump(data))
		fmt.Println("---------------------send msg end---------------------")
	}
}
