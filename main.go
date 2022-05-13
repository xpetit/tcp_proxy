package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

func check(a ...any) {
	for _, v := range a {
		if err, ok := v.(error); ok && err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Println("Usage:  tp REMOTE_ADDRESS [LOCAL_ADDRESS]")
			fmt.Println()
			fmt.Println("Examples:")
			fmt.Println("tp postgres-host.lan:5432                        listens to 127.0.0.1:5432")
			fmt.Println("tp postgres-host.lan:5432 :8888                  listens to localhost:8888")
			fmt.Println("tp postgres-host.lan:5432 [::1]:8888             listens to [::1]:8888")
			os.Exit(1)
		}
	}
}

func main() {
	var remote, local string
	switch len(os.Args) {
	case 3:
		remote = os.Args[1]
		local = os.Args[2]
	case 2:
		remote = os.Args[1]
		_, port, err := net.SplitHostPort(remote)
		check(err)
		local = "127.0.0.1:" + port
	default:
		check(errors.New("invalid number of arguments"))
	}

	conn, err := net.DialTimeout("tcp", remote, 15*time.Second)
	check(err)
	check(conn.Close())

	l, err := net.Listen("tcp", local)
	check(err)
	defer func() { check(l.Close()) }()
	fmt.Println("Forwarding from", l.Addr(), "to", remote)
	for {
		src, err := l.Accept()
		check(err)
		go func() {
			defer func() { check(src.Close()) }()
			dst, err := net.DialTimeout("tcp", remote, 15*time.Second)
			check(err)
			defer func() { check(dst.Close()) }()
			errC := make(chan error)
			cp := func(dst, src net.Conn) {
				_, err := io.Copy(dst, src)
				if errors.Is(err, syscall.ECONNRESET) {
					err = nil
				}
				errC <- err
			}
			go cp(src, dst)
			go cp(dst, src)
			check(<-errC)
			check(<-errC)
		}()
	}
}
