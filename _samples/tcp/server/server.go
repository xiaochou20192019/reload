package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"reload"
	"syscall"
)

var (
	port int
)

func main() {
	flag.IntVar(&port, "p", 8888, `要发送的内容`)
	flag.Parse()

	log.Printf("Actual pid is %d\n", syscall.Getpid())

	listener, err := reload.GetListener(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
	}

	var s = reload.NewService(listener)
	log.Printf("isChild : %v ,listener: %v\n", s.IsChild(), listener)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			s.Add(1)
			log.Println("Accept ", conn.RemoteAddr())
			go recvConnMsg(conn, s)
		}
	}()

	s.Start()

}

func recvConnMsg(conn net.Conn, s reload.Service) {
	//  var buf [4096]byte
	buf := make([]byte, 4096)

	defer conn.Close()
	defer s.Done()

	for {
		n, err := conn.Read(buf)

		if err == io.EOF {
			//连接结束
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		var recv = string(buf[0:n])
		log.Printf("Rev Data : %v", recv)

		conn.Write([]byte(recv))

	}
}