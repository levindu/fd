package main

import (
	"flag"
	"github.com/ftrvxmtrx/fd"
	"io"
	"log"
	"net"
	"os"
)

var (
	socket string
	create bool
)

func init() {
	flag.StringVar(&socket, "s", "/tmp/sendfd.sock", "socket")
	flag.BoolVar(&create, "c", false, "create file")
}

func main() {
	flag.Parse()

	if !flag.Parsed() || socket == "" {
		flag.Usage()
		os.Exit(1)
	}

	c, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	fdConn := c.(*net.UnixConn)

	var fs []*os.File
	fs, err = fd.Get(fdConn, 1, []string{"a file"})
	if err != nil {
		log.Fatal(err)
	}
	f := fs[0]
	defer f.Close()

	if create {
		_, err = io.Copy(f, os.Stdin)
	} else {
		_, err = io.Copy(os.Stdout, f)
	}
	if err != nil {
		log.Fatal(err)
	}
}
