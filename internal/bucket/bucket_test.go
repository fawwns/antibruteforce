package bucket_test

import (
	"testing"
	"time"

	"github.com/fawwns/antibruteforce/internal/bucket"
)

func TestBucketIncrement(t *testing.T) {

	b := bucket.Bucket{
		Limit:      2,
		LastRefill: time.Now(),
	}

	for i := 0; i < 2; i++ {
		if !b.Allow() {
			t.Fatalf("eexpected Allow() = true on iteration %d", i)
		}
	}

	if b.Allow() {
		t.Fatalf("expected Allow() = false when limit exceeded")
	}

}

func TestBucketRefillAfterMinute(t *testing.T) {
	b := bucket.Bucket{
		Count:      1,
		Limit:      1,
		LastRefill: time.Now().Add(-2 * time.Minute),
	}

	if !b.Allow() {
		t.Fatal("expected Allow() after refill period")
	}

	if b.Count != 0 {
		t.Fatalf("expected Count reset to 0, got %d", b.Count)
	}

}
