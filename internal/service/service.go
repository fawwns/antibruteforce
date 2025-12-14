package service

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/fawwns/antibruteforce/internal/bucket"
	"github.com/fawwns/antibruteforce/internal/config"
	"github.com/fawwns/antibruteforce/internal/lists"
	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/fawwns/antibruteforce/internal/postgres"
)

type AntiBruteForceService struct {
	cfg *config.Config

	loginBuckets    *bucket.Store
	passwordBuckets *bucket.Store
	ipBuckets       *bucket.Store

	lists *lists.Lists
	mu    sync.RWMutex
}

func New(cfg *config.Config, storage *postgres.Postgres) *AntiBruteForceService {
	lists := lists.New(storage)
	lists.LoadFromDB(context.Background())

	return &AntiBruteForceService{
		cfg:             cfg,
		loginBuckets:    bucket.NewStore(),
		passwordBuckets: bucket.NewStore(),
		ipBuckets:       bucket.NewStore(),

		lists: lists,
	}
}

func (s *AntiBruteForceService) Check(login, password, ip string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	if s.lists.InWhitelist(clientIP) {
		return true
	}
	if s.lists.InBlacklist(clientIP) {
		return false
	}

	// Login limit.
	if !s.loginBuckets.Get(login, s.cfg.LimitLogin).Allow() {
		logger.Warn.Printf("login limit exceeded: %s", login)
		return false
	}

	// Password limit.
	if !s.passwordBuckets.Get(password, s.cfg.LimitPassword).Allow() {
		logger.Warn.Printf("password limit exceeded")
		return false
	}

	// IP limit.
	if !s.ipBuckets.Get(ip, s.cfg.LimitIP).Allow() {
		logger.Warn.Printf("IP limit exceeded: %s", ip)
		return false
	}

	return true
}

func (s *AntiBruteForceService) InList(ip net.IP, list []net.IPNet) bool {
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

func (s *AntiBruteForceService) ResetBuckets() {
	s.loginBuckets.ResetAll()
	s.passwordBuckets.ResetAll()
	s.ipBuckets.ResetAll()
}

// Методы управления whitelist/blacklist.
func (s *AntiBruteForceService) AddToWhitelist(cidr string) error {
	ctx := context.Background()
	if err := s.lists.AddToWhitelistDB(ctx, cidr); err != nil {
		return err
	}
	return s.lists.AddToWhitelist(cidr)
}

func (s *AntiBruteForceService) RemoveFromWhitelist(cidr string) error {
	ctx := context.Background()
	if err := s.lists.RemoveFromWhitelistDB(ctx, cidr); err != nil {
		return err
	}
	return s.lists.RemoveFromWhitelist(cidr)
}

func (s *AntiBruteForceService) AddToBlacklist(cidr string) error {
	ctx := context.Background()
	if err := s.lists.AddToBlacklistDB(ctx, cidr); err != nil {
		return err
	}
	return s.lists.AddToBlacklist(cidr)
}

func (s *AntiBruteForceService) RemoveFromBlacklist(cidr string) error {
	ctx := context.Background()
	if err := s.lists.RemoveFromBlacklistDB(ctx, cidr); err != nil {
		return err
	}
	return s.lists.RemoveFromBlacklist(cidr)
}
