package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	Port     int
	Name     string
	Conn     net.Conn
	Flag     int
}

var quit chan bool

func NewClient(serverIp string, port int) *Client {
	client := &Client{ServerIp: serverIp, Port: port, Flag: 999}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.Port))

	if err != nil {
		fmt.Println("链接失败" + err.Error())
		return nil
	}
	//defer conn.Close()

	client.Conn = conn
	return client
}

var serviceIp string

var servicePort int

func (c *Client) menu() bool {
	var f int
	fmt.Println("1.群聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出程序")

	_, _ = fmt.Scanln(&f)

	if f < 0 || f > 3 {
		fmt.Println("》》》》》》请选择正确的模式")
		return false
	}

	c.Flag = f
	return true
}

func (c *Client) Run() {

	for c.Flag != 0 {
		for c.menu() != true {
		}
		switch c.Flag {
		case 1:
			c.PubicChat()
			break
		case 2:
			c.PrivateChat()
			break
		case 3:
			c.UpdateName()
			break
		}
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println("请输入需要更新的用户名!")
	var newUserName string
	_, _ = fmt.Scan(&newUserName)

	fakeMsg := "rename|" + newUserName + "\n"
	_, err := c.Conn.Write([]byte(fakeMsg))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true

}
func init() {
	flag.StringVar(&serviceIp, "ip", "127.0.0.1", "default server ip：127.0.0.1")
	flag.IntVar(&servicePort, "port", 8888, "default server port：8888")
	quit = make(chan bool)
}
func main() {
	flag.Parse()
	client := NewClient(serviceIp, servicePort)
	if client == nil {
		fmt.Println("链接失败！")
		return
	}
	go client.DialResponse()
	go client.ReceiverSignal()
	fmt.Println("链接成功！")
	go client.Run()
	<-quit
	fmt.Println("process exit")
}
func (c *Client) DialResponse() {
	_, _ = io.Copy(os.Stdout, c.Conn)
}

func (c *Client) PubicChat() {
	var chatMsg string
	fmt.Println("请输入内容！")
	_, _ = fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.Conn.Write([]byte(sendMsg))
			if err != nil {

				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入内容！")
		_, _ = fmt.Scanln(&chatMsg)
	}

}

func (c *Client) PrivateChat() {
	var toUserName string
	var sendMsg string
	c.SelectUser()
	fmt.Println("请输入要私聊的用户名!exit 退出")
	_, _ = fmt.Scanln(&toUserName)
	for toUserName != "exit" {
		fmt.Println("请输入消息内容！exit退出")
		_, _ = fmt.Scanln(&sendMsg)
		for sendMsg != "exit" {
			if len(sendMsg) != 0 {
				sendMsg = "to|" + toUserName + "|" + sendMsg + "\n"
				fmt.Println("send msg :" + sendMsg)
				_, err := c.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err:", err)
					break
				}
			}
			sendMsg = ""
			fmt.Println("请输入内容！")
			_, _ = fmt.Scanln(&sendMsg)
		}
		c.SelectUser()
		fmt.Println("请输入要私聊的用户名!exit 退出")
		_, _ = fmt.Scanln(&toUserName)
	}

}

func (c *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := c.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err", err)
		return
	}
}

func (c *Client) ReceiverSignal() {
	buffer := make([]byte, 1024)
	for {
		n, err := c.Conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
				quit <- true
			}
			break
		}

		fmt.Printf("Received %d bytes: %s\n", n, string(buffer[:n]))
	}

}
