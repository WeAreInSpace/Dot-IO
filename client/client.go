package client

import (
	"fmt"
	"log"
	"math"
	"net"
	"sync"

	"github.com/WeAreInSpace/dotio/packet"
)

type ApplicationSettings struct {
	Name    string
	Address string

	Mx *sync.Mutex
}

func New(settings *ApplicationSettings) *Application {
	if settings == nil {
		mx := new(sync.Mutex)

		settings = &ApplicationSettings{
			Name:    "Dot I/O Application",
			Address: ":25010",

			Mx: mx,
		}
	}

	if settings.Name == "" {
		settings.Name = "Dot I/O Application"
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

	return &Application{
		mx: settings.Mx,

		ib: ib,
		og: og,
	}
}

type Application struct {
	mx *sync.Mutex

	ib *packet.Inbound
	og *packet.Outgoing
}

func (a *Application) Post(path string, callback func(ib *packet.Inbound, og *packet.Outgoing) error) int32 {
	method := "post"

	reqConnPkBuf := a.og.Write()
	reqConnPkBuf.Sent(packet.WriteInt32(math.MaxInt32))

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

func (a *Application) Put(path string, callback func(ib *packet.Inbound, og *packet.Outgoing) error) int32 {
	method := "put"

	reqConnPkBuf := a.og.Write()
	reqConnPkBuf.Sent(packet.WriteInt32(math.MaxInt32))

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
