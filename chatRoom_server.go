// chatRoom server
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var onlineConns = make(map[string]net.Conn)
var messageQueue = make(chan string, 1000)
var quitQueue = make(chan bool)
// 初始日志输出
var logFile *os.File
var Logger *log.Logger
const (
	LOG_DIR = "./test.log"
)

//处理错误的函数
func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func ProcessInfo(conn net.Conn) {
	buff := make([]byte, 1024)
	//客户端断开后从map清除
	defer func(conn net.Conn) {
		addr :=  fmt.Sprintf("%s", conn.RemoteAddr())
		delete(onlineConns,addr)
		conn.Close()
		for i := range onlineConns{
			fmt.Println("Now online client:" + i)
		}
	}(conn)
	for {
		numOfMessage, err := conn.Read(buff)
		if err != nil {
			break
		}
		if numOfMessage != 0 {
			message := string(buff[:numOfMessage])
			// 拼接发送消息
			messageQueue <- message
		}
	}
}

func consumeMessage() {
	for {
		select {
		case message := <- messageQueue:
			//对消息进行处理
			doProcessMessage(message)
		case <- quitQueue:
			break
		}
	}

}

func doProcessMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		//127.0.0.1:54527#contents：给用户转发聊天内容
		receiverAddr := strings.Trim(contents[0], "")
		sendMessage := strings.Join(contents[1:], "#")
		//如果接收用户在线，转发消息
		if conn, ok := onlineConns[receiverAddr]; ok {
			_, err := conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Printf("%s has offline\n", receiverAddr)
			}

		}
	}else {
		//127.0.0.1:54527*list：给用户返回当前在线用户列表
		contents = strings.Split(message, "*")
		if strings.ToUpper(contents[1]) == "LIST"{
			var ips string = "在线用户列表："
			for i := range onlineConns{
				ips = ips + "|" + i
			}
			if conn, ok := onlineConns[contents[0]]; ok {
				_, err := conn.Write([]byte(ips))
				if err != nil {
					fmt.Printf("%s has offline\n", contents[0])
				}
			}
		}
	}
}

func main() {
	//日志配置
	logFile, err := os.OpenFile(LOG_DIR, os.O_RDWR|os.O_CREATE, 0)
	if err != nil{
		fmt.Println("LogFile open failure")
		os.Exit(-1)
	}
	defer logFile.Close()
	Logger = log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	//开始启动服务
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	checkError(err)
	defer listen.Close()
	Logger.Println("Server has running....")
	fmt.Println("Server has running ......")
	// 转发消息
	go consumeMessage()

	// 循环进行监听请求
	for {
		conn, err := listen.Accept()
		checkError(err)
		//将连接用户存入在线列表
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		onlineConns[addr] = conn
		for i := range onlineConns {
			fmt.Printf("在线用户：%s\n", i)
		}
		//接收消息
		go ProcessInfo(conn)

	}
}
