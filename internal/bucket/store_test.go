package bucket_test

import (
	"testing"
	"time"

	"github.com/fawwns/antibruteforce/internal/bucket"
)

func TestStoreGetCreatesBucket(t *testing.T) {
	b := bucket.NewStore()

	first := b.Get("123", 2)

	if first == nil {
		t.Fatal("expected bucket, got nil")
	}

	second := b.Get("123", 2)
	if first != second {
		t.Fatal("expected same bucket instance on repeated Get()")
	}

}

func TestStoreCleanup(t *testing.T) {
	s := bucket.NewStore()

	b := s.Get("old", 5)
	b.LastRefill = time.Now().Add(-3 * time.Minute)

	s.Cleanup(1 * time.Minute)

	if _, exists := s.Buckets["old"]; exists {
		t.Fatal("expected old bucket to be cleaned up")
	}
}
