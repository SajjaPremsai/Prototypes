package server

import (
	"log"
	"net"
	"syscall"

	"github.com/redis/config"
	"github.com/redis/core"
)

var con_clients int = 0

func RunAsyncTCPServer() error {
	log.Printf("Starting an asynchronous TCP server on %s:%d", config.Host, config.Port)
	max_clients := 20000
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	if err := syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	ip4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[12], ip4[13], ip4[14], ip4[15]},
	}); err != nil {
		return err
	}

	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer syscall.Close(epollFD)

	var socketServerEvents syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvents); err != nil {
		return err
	}

	for {
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)

		if e != nil {
			return e
		}

		for i := 0; i < nevents; i++ {
			if int(events[i].Fd) == serverFD {

				fd, addr, err := syscall.Accept(serverFD)

				if err != nil {

					if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
						continue
					}

					log.Printf("Error accepting connection: %v", err)
					continue
				}

				con_clients++

				if err := syscall.SetNonblock(fd, true); err != nil {
					log.Printf("Error setting client socket non-blocking: %v", err)
					syscall.Close(fd)
					continue
				}

				if clientAddr, ok := addr.(*syscall.SockaddrInet4); ok {

					ip := net.IP(clientAddr.Addr[:])

					log.Printf(
						"Client connected: %s:%d (FD=%d) Concurrent clients: %d",
						ip.String(),
						clientAddr.Port,
						fd,
						con_clients,
					)
				}

				var socketClientEvents syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvents); err != nil {
					log.Printf("Error adding client socket to epoll: %v", err)
					syscall.Close(fd)
					continue
				}

			} else {
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommand(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					con_clients -= 1
					log.Printf("Client disconnected: %d (Concurrent clients: %d)", events[i].Fd, con_clients)
					continue
				}
				respond(cmd, comm)
			}
		}
	}
}
