package lists

import (
	"context"
	"net"
	"sync"

	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/fawwns/antibruteforce/internal/postgres"
)

type Lists struct {
	whitelist []net.IPNet
	blacklist []net.IPNet
	mu        sync.RWMutex

	storage *postgres.Postgres
}

func New(storage *postgres.Postgres) *Lists {
	return &Lists{
		whitelist: make([]net.IPNet, 0),
		blacklist: make([]net.IPNet, 0),
		storage:   storage,
	}
}

func (l *Lists) LoadFromDB(ctx context.Context) error {
	w, err := l.storage.GetWhitelist(ctx)
	if err != nil {
		return err
	}
	b, err := l.storage.GetBlacklist(ctx)
	if err != nil {
		return err
	}

	for _, cidr := range w {
		l.AddToWhitelist(cidr)
	}
	for _, cidr := range b {
		l.AddToBlacklist(cidr)
	}

	return nil
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

	l.blacklist = newList
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

func (l *Lists) AddToWhitelistDB(ctx context.Context, cidr string) error {
	return l.storage.AddToWhitelist(ctx, cidr)
}

func (l *Lists) RemoveFromWhitelistDB(ctx context.Context, cidr string) error {
	return l.storage.RemoveFromWhitelist(ctx, cidr)
}

func (l *Lists) AddToBlacklistDB(ctx context.Context, cidr string) error {
	return l.storage.AddToBlacklist(ctx, cidr)
}

func (l *Lists) RemoveFromBlacklistDB(ctx context.Context, cidr string) error {
	return l.storage.RemoveFromBlacklist(ctx, cidr)
}
