package limiter

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type (
	// IndividualMemoryStore is the built-in store implementation for RateLimiter
	IndividualMemoryStore struct {
		receivers   map[string]*Receiver
		mutex       sync.Mutex
		rate        rate.Limit    // Rate of request allowed to pass as req/s
		burst       int           // Burst addtionally allows a number of requests to pass when rate limit is reached
		expiresIn   time.Duration // The duration after that a rate limiter is clean up
		lastCleanup time.Time
	}

	//Receiver signifies a unique user's limiter details
	Receiver struct {
		*rate.Limiter
		lastSeen time.Time
	}
)

func NewIndividualMemoryStore(rate rate.Limit) (store *IndividualMemoryStore) {
	return NewIndividualMemoryStoreWithConfig(IndividualMemoryStoreConfig{Rate: rate})
}

func NewIndividualMemoryStoreWithConfig(config IndividualMemoryStoreConfig) (store *IndividualMemoryStore) {
	store = &IndividualMemoryStore{}

	store.rate = config.Rate
	store.burst = config.Burst
	store.expiresIn = config.ExpiresIn
	if config.ExpiresIn == 0 {
		store.expiresIn = DefaultMemoryStoreConfig.ExpiresIn
	}
	if config.Burst == 0 {
		store.burst = int(config.Rate)
	}
	store.receivers = make(map[string]*Receiver)
	store.lastCleanup = time.Now()
	return
}

type IndividualMemoryStoreConfig struct {
	Rate      rate.Limit    // Rate of requests allowed to pass as req/s
	Burst     int           // Burst additionally allows a number of requests to pass when rate limit is reached
	ExpiresIn time.Duration // ExpiresIn is the duration after that a rate limiter is cleaned up
}

// DefaultMemoryStoreConfig provides default configuration values for RateLimiterMemoryStore
var DefaultMemoryStoreConfig = IndividualMemoryStoreConfig{
	ExpiresIn: 3 * time.Second,
}

// Allow implements RateLimiterStore.Allow
func (store *IndividualMemoryStore) Allow(identifier string) (bool, error) {
	store.mutex.Lock()
	limiter, exists := store.receivers[identifier]
	if !exists {
		limiter = new(Receiver)
		limiter.Limiter = rate.NewLimiter(store.rate, store.burst)
		store.receivers[identifier] = limiter
	}

	limiter.lastSeen = time.Now()
	if time.Now().Sub(store.lastCleanup) > store.expiresIn {
		store.cleanupStaleReceivers()
	}
	store.mutex.Unlock()

	allow := limiter.Allow()
	if allow {
		return allow, nil
	} else {
		return allow, errors.New("Individual rate limiter denied")
	}
}

// cleanupStaleReceivers helps manage the size of the receivers map by removing
// stale records of users who haven't visited again after the configured expiry
// time has elapsed
func (store *IndividualMemoryStore) cleanupStaleReceivers() {
	for id, receiver := range store.receivers {
		if time.Now().Sub(receiver.lastSeen) > store.expiresIn {
			delete(store.receivers, id)
		}
	}

	store.lastCleanup = time.Now()
}
