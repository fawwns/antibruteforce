package lists

import (
	"net"
	"sync"

	"github.com/fawwns/antibruteforce/internal/logger"
)

type Lists struct {
	whitelist []net.IPNet
	blacklist []net.IPNet
	mu        sync.RWMutex
}

func New() *Lists {
	return &Lists{
		whitelist: make([]net.IPNet, 0),
		blacklist: make([]net.IPNet, 0),
	}
}

// AddToWhitelist добавляет CIDR-сеть в белый список.
func (l *Lists) AddToWhitelist(cidr string) error {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.whitelist = append(l.whitelist, *network)
	logger.Info.Println("Network " + network.String() + " added to the whitelist")
	return nil
}

// AddToBlacklist добавляет CIDR-сеть в белый список.
func (l *Lists) AddToBlacklist(cidr string) error {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.blacklist = append(l.blacklist, *network)
	logger.Info.Println("Network " + network.String() + " added to the blacklist")
	return nil
}

// RemoveFromWhitelist — удаление сети из whitelist.
func (l *Lists) RemoveFromWhitelist(cidr string) error {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	newList := make([]net.IPNet, 0)
	for _, n := range l.whitelist {
		if n.String() != network.String() {
			newList = append(newList, n)
		}
	}

	l.whitelist = newList
	logger.Info.Println("Network " + network.String() + " removed from whitelist")
	return nil
}

// RemoveFromBlacklist — удаление сети из blacklist.
func (l *Lists) RemoveFromBlacklist(cidr string) error {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	newList := make([]net.IPNet, 0)
	for _, n := range l.blacklist {
		if n.String() != network.String() {
			newList = append(newList, n)
		}
	}

	l.whitelist = newList
	logger.Info.Println("Network " + network.String() + " removed from blacklist")
	return nil
}

// InWhitelist проверяет, входит ли IP в белый список.
func (l *Lists) InWhitelist(ip net.IP) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, val := range l.whitelist {
		if val.Contains(ip) {
			return true
		}
	}
	return false
}

// InBlacklist проверяет, входит ли IP в чёрный список.
func (l *Lists) InBlacklist(ip net.IP) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, val := range l.blacklist {
		if val.Contains(ip) {
			return true
		}
	}
	return false
}
