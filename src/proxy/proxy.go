package proxy

import (
	"net"
	"io"
	"log"
	"math/rand"
	"time"
	"fmt"
	"os"
)

var (
	sessions = map[string]*Session{}
)

type Session struct {
	ID string
	from string
	to string
	bytesUp int64
	bytesDown int64
}

func cp(dst io.ReadWriteCloser, src io.ReadWriteCloser,sessionID string,direction string, errChan chan error){
	buf := make([]byte, 32*1024)
	var err  error= nil
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written := int64(nw)
				session, ok := sessions[sessionID]
				if ok {
					switch direction {
					case "up":
						session.bytesUp += written
					case "down":
						session.bytesDown += written
					}
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	if (err != nil){
		errChan <- err
	}
}

func connect(upStream io.ReadWriteCloser, downStream io.ReadWriteCloser, sessionID string){
	errChan := make(chan error, 1)
	go cp(downStream, upStream, sessionID, "down", errChan)
	go cp(upStream, downStream, sessionID, "up", errChan)

	if err := <- errChan;err != nil {
		log.Printf("[%s] %s\n", sessionID, err)
	}

	upStream.Close()
	downStream.Close()
}

func selectBackend(backends[] string)(upStream string){
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(backends))
	return backends[r]
}

func genSessionId(from string, to string)(sessionId string){
	return fmt.Sprintf("%s|%s|%s", time.Now().Format("060102_150405.999"), from , to)
}

func sessionLog(){
	accessLoger := log.New(os.Stdout, "", log.LstdFlags)
	for{
		<-time.After(time.Second*10)
		for _, session := range sessions {
			accessLoger.Printf("%s %s %s %d %d",session.ID, session.from, session.to, session.bytesUp, session.bytesDown)
		}
	}
}

func ProxyServer(backends []string, listen string){
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		log.Printf("error: listen %s, %s\n", listen, err)
		return
	}
	log.Println("wormhole server started")
	log.Printf("listen to addr [%s]\n", listen)
	log.Printf("ready to proxy to backends %v", backends)
	go sessionLog()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error: sock accept %s\n", err)
			continue
		}
		backend := selectBackend(backends)
		go func(conn net.Conn){
			upStream, er := net.Dial("tcp", backend)
			if(er != nil){
				log.Printf("error: dail to upstream %s, %s\n", backend, er)
				conn.Close()
				return
			}
			sessionID := genSessionId(conn.RemoteAddr().String(), backend)
			sessions[sessionID] = &Session{ID: sessionID, from: conn.RemoteAddr().String(), to: backend, bytesUp: 0, bytesDown: 0}
			log.Printf("[%s] proxy from [%s] to [%s] created", sessionID , conn.RemoteAddr(), backend)
			connect(upStream, conn, sessionID)
			delete(sessions, sessionID)
			log.Printf("[%s] proxy from [%s] to [%s] closed", sessionID , conn.RemoteAddr(), backend)
		}(conn)
	}
}

