package main

import (
	"flag"
	"github.com/ftrvxmtrx/fd"
	"log"
	"net"
	"os"
)

var (
	filename string
	socket   string
	create bool
)

func init() {
	flag.StringVar(&filename, "f", "", "filename")
	flag.StringVar(&socket, "s", "/tmp/sendfd.sock", "socket")
	flag.BoolVar(&create, "c", false, "create file")
}

func main() {
	flag.Parse()

	if !flag.Parsed() || filename == "" || socket == "" {
		flag.Usage()
		os.Exit(1)
	}

	var (
		f *os.File
		err error
	)
	if create {
		f, err = os.Create(filename)
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = os.Chmod(socket, 0o777)
	if err != nil {
		log.Fatal(err)
	}

	var a net.Conn
	a, err = l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	listenConn := a.(*net.UnixConn)
	if err = fd.Put(listenConn, f); err != nil {
		log.Fatal(err)
	}
}
