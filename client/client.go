package client

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/WeAreInSpace/dotio/packet"
)

type ApplicationSettings struct {
	Name    string
	Address string

	Mx *sync.Mutex
}

func New(settings *ApplicationSettings) *application {
	if settings == nil {
		mx := new(sync.Mutex)

		settings = &ApplicationSettings{
			Name:    "Dot I/O application",
			Address: ":25010",

			Mx: mx,
		}
	}

	if settings.Name == "" {
		settings.Name = "Dot I/O application"
	}

	if settings.Address == "" {
		settings.Address = ":25010"
	}

	if settings.Mx == nil {
		mx := new(sync.Mutex)
		settings.Mx = mx
	}

	serverAddr, err := net.ResolveTCPAddr("tcp", settings.Address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	ib := &packet.Inbound{
		Conn: conn,
	}

	og := &packet.Outgoing{
		Conn: conn,
	}

	return &application{
		mx: settings.Mx,

		ib: ib,
		og: og,
	}
}

type application struct {
	mx *sync.Mutex

	ib *packet.Inbound
	og *packet.Outgoing
}

func (a *application) Post(path string, callback func(ib *packet.Inbound, og *packet.Outgoing) error) int32 {
	method := "post"

	firstPkBuf := a.og.Write()
	firstPkBuf.WriteString(method)
	firstPkBuf.WriteString(path)
	firstPkBuf.Sent(packet.WriteInt32(0))

	statusRes, _, err := a.ib.Read()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}

	err = callback(a.ib, a.og)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}

	return statusRes
}

func (a *application) Put(path string, callback func(ib *packet.Inbound, og *packet.Outgoing) error) int32 {
	method := "put"

	firstPkBuf := a.og.Write()
	firstPkBuf.WriteString(method)
	firstPkBuf.WriteString(path)
	firstPkBuf.Sent(packet.WriteInt32(0))

	statusRes, _, err := a.ib.Read()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}

	err = callback(a.ib, a.og)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}

	return statusRes
}
