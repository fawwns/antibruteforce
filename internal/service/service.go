package service

import (
	"net"
	"sync"
	"time"

	"github.com/fawwns/antibruteforce/internal/bucket"
	"github.com/fawwns/antibruteforce/internal/config"
	"github.com/fawwns/antibruteforce/internal/lists"
	"github.com/fawwns/antibruteforce/internal/logger"
)

type AntiBruteForceService struct {
	cfg *config.Config

	loginBuckets    *bucket.Store
	passwordBuckets *bucket.Store
	ipBuckets       *bucket.Store

	lists *lists.Lists
	mu    sync.RWMutex
}

func New(cfg *config.Config) *AntiBruteForceService {
	return &AntiBruteForceService{
		cfg:             cfg,
		loginBuckets:    bucket.NewStore(),
		passwordBuckets: bucket.NewStore(),
		ipBuckets:       bucket.NewStore(),

		lists: lists.New(),
	}
}

func (s *AntiBruteForceService) Check(login, password, ip string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Проверка whitelist/blacklist.
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}
	if s.lists.InWhitelist(clientIP) {
		return true
	} else if s.lists.InBlacklist(clientIP) {
		return false
	}

	// Проверка лимитов.
	if !s.ipBuckets.Get(login, s.cfg.LimitLogin).Allow() {
		logger.Warn.Printf("login limit exceeded: %s", ip)
		return false
	}
	if !s.passwordBuckets.Get(password, s.cfg.LimitPassword).Allow() {
		logger.Warn.Printf("password limit exceeded: %s", ip)
		return false
	}
	if !s.loginBuckets.Get(ip, s.cfg.LimitIP).Allow() {
		logger.Warn.Printf("IP limit exceeded: %s", ip)
		return false
	}

	return true

}

func (s *AntiBruteForceService) inList(ip net.IP, list []net.IPNet) bool {
	for _, netw := range list {
		if netw.Contains(ip) {
			return true
		}
	}
	return false
}

// Очистка неактивных бакетов.
func (s *AntiBruteForceService) Cleanup() {
	ttl := time.Minute * 2
	s.loginBuckets.Cleanup(ttl)
	s.passwordBuckets.Cleanup(ttl)
	s.ipBuckets.Cleanup(ttl)
}

// Методы управления whitelist/blacklist.
func (s *AntiBruteForceService) AddToWhitelist(cidr string) error {
	return s.lists.AddToWhitelist(cidr)
}

func (s *AntiBruteForceService) AddToBlacklist(cidr string) error {
	return s.lists.AddToBlacklist(cidr)
}

func (s *AntiBruteForceService) RemoveFromWhitelist(cidr string) error {
	return s.lists.RemoveFromWhitelist(cidr)
}

func (s *AntiBruteForceService) RemoveFromBlacklist(cidr string) error {
	return s.lists.RemoveFromBlacklist(cidr)
}
