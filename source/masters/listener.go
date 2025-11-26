package masters

import (
	"Nosviak4/source/functions/fuzzyProxy"
	"fmt"
	"net"
)

// Listen will bind to the address information presented
func (m *Masters) Listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", m.binder.Address, m.binder.Port))
	if err != nil {
		return err
	}

	go m.spawnMasterTicker() 
	proxy := fuzzyproxy.New(fmt.Sprintf("%s:%d", m.binder.FuzzyProxy.Address, m.binder.FuzzyProxy.Port), m.binder.FuzzyProxy.Secret)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go m.acceptConn(conn, proxy)
	}
}
