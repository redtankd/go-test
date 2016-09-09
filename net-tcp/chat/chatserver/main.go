package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"runtime"
	"flag"
)

func main() {
	var (
		host = flag.String("host", "localhost", "host")
		port = flag.String("port", "6000", "port")
	)
	flag.Parse()

	log.SetFlags(log.Lshortfile)
	runtime.GOMAXPROCS(2)

	ln, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	b := newBroadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("Connected from %s (%d clients)", conn.RemoteAddr(), b.Count()+1)
		b.Add(conn)
		go handleConnection(conn, b)
	}
}

func handleConnection(conn net.Conn, bc *broadcaster) {
	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				conn.Close()
				bc.Remove(conn)
				break
			}
			bc.Send(line)
		}
	}()
}
