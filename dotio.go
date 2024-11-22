package dotio

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
	Wg *sync.WaitGroup
}

func New(settings *ApplicationSettings) *application {
	if settings == nil {
		mx := new(sync.Mutex)
		wg := new(sync.WaitGroup)

		settings = &ApplicationSettings{
			Name:    "Dot I/O application",
			Address: ":25010",

			Mx: mx,
			Wg: wg,
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

	if settings.Wg == nil {
		wg := new(sync.WaitGroup)
		settings.Wg = wg
	}

	fmt.Printf("Starting: %s\n", settings.Name)
	fmt.Printf("Listen on: %s\n", settings.Address)

	serverAddr, err := net.ResolveTCPAddr("tcp", settings.Address)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	device := make(map[net.Addr]*device)
	route := make(map[[2]string]*route)

	return &application{
		listener: listener,

		mx: settings.Mx,
		wg: settings.Wg,

		devices: device,
		routes:  route,
	}
}

type application struct {
	listener *net.TCPListener

	mx *sync.Mutex
	wg *sync.WaitGroup

	devices map[net.Addr]*device
	routes  map[[2]string]*route
}

type device struct {
	addr net.Addr
	conn *net.TCPConn

	ib *packet.Inbound
	og *packet.Outgoing
}

func (a *application) Listen() {
	defer a.wg.Done()

	go a.connectionHanler()
	a.wg.Add(1)

	a.wg.Wait()
}

func (a *application) connectionHanler() {
	defer a.listener.Close()

	for {
		conn, err := a.listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		go a.deviceHanler(conn)
	}
}

func (a *application) deviceHanler(conn *net.TCPConn) {
	if _, exits := a.devices[conn.RemoteAddr()]; !exits {
		addr := conn.RemoteAddr()
		ib := packet.Inbound{
			Conn: conn,
		}
		og := packet.Outgoing{
			Conn: conn,
		}

		newDevice := &device{
			addr: addr,
			conn: conn,
			ib:   &ib,
			og:   &og,
		}

		a.devices[addr] = newDevice
		log.Printf("Device: %v\n", a.devices)

		a.load(addr)

		defer a.deleteDevice(addr)
	}
}

func (a *application) load(addr net.Addr) {
	device := a.devices[addr]

	for {
		id, firstPkBuf, err := device.ib.Read()
		if err != nil {
			log.Printf("ERROR: %s", err)
			break
		}

		if id == 0 {
			method := firstPkBuf.ReadString()
			path := firstPkBuf.ReadString()

			log.Printf("Requst from: %s to: %s method: %s\n", addr.String(), path, method)

			route, exits := a.routes[[2]string{method, path}]
			if exits {
				err := route.callback(device.ib, device.og)
				if err != nil {
					log.Printf("ERROR: %s", err)
					break
				}

				res := device.og.Write()
				res.Sent(packet.WriteInt32(0))
			} else {
				log.Printf("ERROR: Function '%s' at '%s' does not exits\n", method, path)
				res := device.og.Write()
				res.Sent(packet.WriteInt32(1))
			}
			break
		}
	}
}

func (a *application) deleteDevice(addr net.Addr) {
	delete(a.devices, addr)
	log.Printf("Device: %v\n", a.devices)
}

type route struct {
	method   string
	callback func(ib *packet.Inbound, og *packet.Outgoing) error
}

func (a *application) Post(path string, callback func(ib *packet.Inbound, og *packet.Outgoing) error) {
	method := "post"
	fmt.Printf("Register: method '%s' at '%s'\n", method, path)

	if _, exits := a.routes[[2]string{method, path}]; exits {
		log.Printf("ERROR: Function '%s' at '%s' already exits.", method, path)
		return
	}

	newRoute := &route{
		method:   method,
		callback: callback,
	}

	a.routes[[2]string{method, path}] = newRoute
}

func (a *application) Put(path string, callback func(ib *packet.Inbound, og *packet.Outgoing) error) {
	method := "put"
	fmt.Printf("Register: method '%s' at '%s'\n", method, path)

	if _, exits := a.routes[[2]string{method, path}]; exits {
		log.Printf("ERROR: Function '%s' at '%s' already exits.", method, path)
		return
	}

	newRoute := &route{
		method:   method,
		callback: callback,
	}

	a.routes[[2]string{method, path}] = newRoute
}
