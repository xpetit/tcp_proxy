package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	. "github.com/xpetit/x/v2"
)

func main() {
	var remote, local string
	switch len(os.Args) {
	case 3:
		remote = os.Args[1]
		local = os.Args[2]
	case 2:
		remote = os.Args[1]
		_, port := C3(net.SplitHostPort(remote))
		local = "127.0.0.1:" + port
	default:
		fmt.Fprintln(os.Stderr, "invalid number of arguments")
		fmt.Println("Usage:  tcp_proxy REMOTE_ADDRESS [LOCAL_ADDRESS]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("tcp_proxy postgres-host.lan:5432                        listens to 127.0.0.1:5432")
		fmt.Println("tcp_proxy postgres-host.lan:5432 :8888                  listens to localhost:8888")
		fmt.Println("tcp_proxy postgres-host.lan:5432 [::1]:8888             listens to [::1]:8888")
		os.Exit(1)
	}

	// try dialing first to fail early
	C(C2(net.DialTimeout("tcp", remote, 15*time.Second)).Close())

	l := C2(net.Listen("tcp", local))
	defer Closing(l)
	fmt.Println("Forwarding from", l.Addr(), "to", remote)
	for {
		src := C2(l.Accept())
		go func() {
			defer Closing(src)

			dst := C2(net.DialTimeout("tcp", remote, 15*time.Second))
			go func() {
				C2(io.Copy(dst, src))
				C(dst.Close())
			}()
			_, err := io.Copy(src, dst)
			Assert(errors.Is(err, net.ErrClosed), "dst connection is closed")
		}()
	}
}
