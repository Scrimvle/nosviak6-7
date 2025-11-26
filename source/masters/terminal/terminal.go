package terminal

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source/masters/terminal/signal"
	"bytes"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Terminal can even be referenced to with blank fields
type Terminal struct {
	X, Y       uint32
	Config     *ServerConfig
	Screen     *bytes.Buffer
	Channel    ssh.Channel
	Conn       *ssh.ServerConn
	connLogger *gologr.Terminal
	Signal     *signal.TerminalSignal
	ConnTime   time.Time
	writeLock  sync.Mutex
	Addr       net.Addr
	
	// if the ssh client supports XTerm (well certain escapes)
	XTerm      bool
}

type ServerConfig struct {
	Port            int    `json:"port"`
	Address         string `json:"address"`
	ServerKey       string `json:"server_key"`
	MaxAuthAttempts int    `json:"max_auth_attempts"`
	LoginTimeout    int    `json:"login_timeout"`
	FuzzyProxy      struct {
		Enabled bool   `json:"enabled"`
		Address string `json:"address"`
		Port    int    `json:"port"`
		Secret  string `json:"secret"`
	} `toml:"fuzzyProxy"`
	Captcha struct {
		Enabled     bool     `json:"enabled"`
		IgnorePerms []string `json:"ignore_perms"`
		Max         int      `json:"max"`
		Min         int      `json:"min"`
	} `toml:"captcha"`

	// CustomAuth will implement the custom authentication screen. 
	CustomAuth      bool   `json:"custom_auth"`
}

// NewTerminal will create an interface for interacting with connections
func NewTerminal(conn *ssh.ServerConn, channel ssh.Channel, logger *gologr.Terminal, config *ServerConfig, addr net.Addr) *Terminal {
	return &Terminal{
		X:          80,
		Y:          24,
		Addr:       addr,
		Conn:       conn,
		XTerm:      false,
		Config:     config,
		Screen:     bytes.NewBuffer(make([]byte, 0)),
		Channel:    channel,
		connLogger: logger,
		ConnTime:   time.Now(),
	}
}

// Logger will return the logger which is used for returning values to the terminal
func (t *Terminal) Logger() *gologr.Terminal {
	return t.connLogger
}