package fuzzyproxy

import (
	"net"
	"strconv"
	"time"
)

// Resolve will check the cache and try find the original IP
func (P *Proxy) Resolve(addr net.Addr) (*IP, error) {
	_, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	P.mutex.RLock()
	entry, ok := P.cache[id]
	P.mutex.RUnlock()

	if ok {
		if entry.Expire.Unix() > time.Now().Unix() {
			return entry, nil
		}

		// removes from cache
		delete(P.cache, id)
	}

	entry = new(IP)
	entry.Expire = time.Now().Add(5 * time.Minute)
	entry.Proxy = P
	entry.ID = id

	if err := entry.fetch(); err != nil {
		return nil, err
	}

	return entry, nil
}