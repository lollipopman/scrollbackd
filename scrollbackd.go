package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func handleConnection(conn net.Conn, scrollbackFile *os.File) {
	defer conn.Close()
	_, err := scrollbackFile.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	scrollbackScanner := bufio.NewScanner(scrollbackFile)
	for scrollbackScanner.Scan() {
		io.WriteString(conn, scrollbackScanner.Text()+"\r\n")
	}

	connScanner := bufio.NewScanner(conn)
	for connScanner.Scan() {
		log.Print(connScanner.Text())
		io.WriteString(scrollbackFile, connScanner.Text()+"\n")
	}
}

func main() {
	log.SetPrefix("scrollbackd: ")
	log.SetFlags(0)
	var port = flag.Int("p", 8000, "Port to listen on")
	var ip = flag.String("l", "localhost", "IP to listen on, 0.0.0.0 for all interfaces")
	var scrollbackFilename = flag.String("f", "", "File to store scrollback data")
	flag.Parse()

	if *scrollbackFilename == "" {
		flag.PrintDefaults()
		log.Fatal("You must supply a file to write the scrollback data")
	}

	scrollbackFile, err := os.OpenFile(*scrollbackFilename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", *ip+":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handleConnection(conn, scrollbackFile)
	}
}
