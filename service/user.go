package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	Ch     chan string
	Conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{Name: userAddr, Addr: userAddr, Ch: make(chan string), Conn: conn, server: server}
	go user.ListenMessage()

	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Ch
		_, _ = u.Conn.Write([]byte(msg + "\n"))
	}
}
func (u *User) Online() {

	u.server.mapLock.Lock()
	_, ok := u.server.OnlineMap[u.Name]
	if !ok {
		u.server.OnlineMap[u.Name] = u
	}
	u.server.mapLock.Unlock()
	u.server.BroadCast(u, "online")
}
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	u.server.BroadCast(u, "offline")
}
func (u *User) DoMsg(msg string) {

	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			if user == u {
				continue
			}
			sendMsg := "[" + user.Addr + "]" + user.Name + "在线"
			u.SendMsg(sendMsg)
		}
		u.server.mapLock.Unlock()
		return
	}
	// 改名
	if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("用户名已存在!")
			return
		}

		u.server.mapLock.Lock()
		delete(u.server.OnlineMap, u.Name)
		u.server.OnlineMap[newName] = u
		u.server.mapLock.Unlock()
		u.Name = newName
		u.SendMsg("已经更新为：" + newName)
		return
	}
	if len(msg) > 4 && msg[:3] == "to|" {
		toName := strings.Split(msg, "|")[1]

		if toName == "" {
			u.SendMsg("消息格式不正确，请使用\"to|name|msg格式\"")
			return
		}
		if _, ok := u.server.OnlineMap[toName]; !ok {
			u.SendMsg("该用户名称不存在")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("发送消息为空\n")
			return
		}
		toUser := u.server.OnlineMap[toName]
		toUser.SendMsg("from:" + u.Name + "|send msg:" + content)
		return

	}
	u.server.BroadCast(u, msg)

}

func (u *User) SendMsg(msg string) {
	_, _ = u.Conn.Write([]byte(msg + "\n"))
}
