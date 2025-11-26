package masters

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	fuzzyproxy "Nosviak4/source/functions/fuzzyProxy"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/signal"
	"encoding/binary"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

// accept will handle the connection and prepare to pass onto the terminal package
func (m *Masters) acceptConn(conn net.Conn, proxy *fuzzyproxy.Proxy) {
	termLog, IP := m.logger.WithTerminal(), conn.RemoteAddr()

	// if proxy is enabled, we attempt to resolve it
	if m.binder.FuzzyProxy.Enabled && strings.Contains(IP.String(), source.OPTIONS.String("ssh", "fuzzyProxy", "address")) {
		proxier, err := proxy.Resolve(IP)
		if err == nil {
			IP = proxier
		}
	}

	serverConn, channels, requests, err := ssh.NewServerConn(conn, m.config)
	if err != nil {
		termLog.WriteLog(gologr.ERROR, "[SSH-CONN] fatal-conn error occurred with %s: %v", conn.RemoteAddr().String(), err)
		return
	}

	go ssh.DiscardRequests(requests)
	for channel := range channels {
		if channel.ChannelType() != "session" {
			termLog.WriteLog(gologr.ERROR, "[SSH-CONN] rejected channel type with conn from %s: %s", conn.RemoteAddr().String(), channel.ChannelType())
			channel.Reject(ssh.UnknownChannelType, "UnknownChannelType")
			return
		}

		channel, requests, err := channel.Accept()
		if err != nil {
			termLog.WriteLog(gologr.ERROR, "[SSH-CONN] unable to accept channel from %s: %v", conn.RemoteAddr().String(), err)
			return
		}

		term := terminal.NewTerminal(serverConn, channel, termLog, m.binder, IP)
		term.Signal = signal.NewSignaller(channel)

		go func() {
			for request := range requests {
				switch request.Type {

				default:
					request.Reply(true, make([]byte, 0))

				case "pty-req", "shell", "exec":
					err := request.Reply(true, make([]byte, 0))
					if err != nil {
						return
					}

				case "window-change":
					term.X, term.Y = binary.BigEndian.Uint32(request.Payload), binary.BigEndian.Uint32(request.Payload[4:])
					if !request.WantReply {
						continue
					}

					request.Reply(true, make([]byte, 0))
				}
			}
		}()

		if err := Login(term); err != nil {
			termLog.WriteLog(gologr.ERROR, "[SSH-CONN] fatal error occurred inside internal session for %s: %v", conn.RemoteAddr().String(), err)
			return
		}
	}
}
