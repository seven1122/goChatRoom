// chatRoom client
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func MessageSend(conn net.Conn) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}

		_, err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Println("client connect failure: " + err.Error())
			break
		}

	}

}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	checkError(err)
	defer conn.Close()

	// 给服务端发送消息
	go MessageSend(conn)

	// 接收服务端的消息
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil{
			fmt.Println("你已经退出聊天室")
			os.Exit(0)
		}
		if length != 0 {
			fmt.Println("Receive server message content: " + string(buf[:length]))
		}
	}

	fmt.Println("client bye bye ....")

}
