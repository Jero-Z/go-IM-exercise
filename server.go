package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{Ip: ip, Port: port, OnlineMap: make(map[string]*User), Message: make(chan string)}
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net listener err")
		return
	}
	defer listen.Close()
	// 启动监听全局消息分发器
	go s.DispatchMsg()

	for {
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println("net listener err", err)
			continue
		}
		go s.Handler(accept)

	}
}

func (s *Server) Handler(conn net.Conn) {
	conn.RemoteAddr()

	user := NewUser(conn, s)
	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {

			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMsg(msg)

			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:

		case <-time.After(10 * time.Minute):
			user.SendMsg("超时下线！")
			close(user.Ch)
			_ = conn.Close()
			return
		}
	}

}

func (s *Server) BroadCast(fromUser *User, msg string) {
	sendMsg := "[" + fromUser.Addr + fromUser.Name + "]" + msg

	s.Message <- sendMsg

}
func (s *Server) DispatchMsg() {
	for {
		msg := <-s.Message

		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.Ch <- msg
		}
		s.mapLock.Unlock()
	}
}
