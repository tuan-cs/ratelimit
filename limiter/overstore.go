package limiter

import (
	"errors"
	"time"

	"golang.org/x/time/rate"
)

// type TotalStore interface {
// 	Allow() (bool, error)
// }

type TotalMemoryStore struct {
	rate  rate.Limit // Rate of request allowed to pass as req/s
	burst int        // Burst addtionally allows a number of requests to pass when rate limit is reached

	limiter *rate.Limiter
}

func NewTotalMemoryStore(rate rate.Limit) (store *TotalMemoryStore) {
	return NewTotalMemoryStoreWithConfig(TotalMemoryStoreConfig{Rate: rate})
}

func NewTotalMemoryStoreWithConfig(config TotalMemoryStoreConfig) (store *TotalMemoryStore) {
	store = &TotalMemoryStore{}
	store.rate = config.Rate
	store.burst = config.Burst
	store.limiter = rate.NewLimiter(store.rate, store.burst)

	return
}

type TotalMemoryStoreConfig struct {
	Rate      rate.Limit
	Burst     int
	ExpiresIn time.Duration
}

var DefaultTotalMemoryStoreConfig = TotalMemoryStoreConfig{
	ExpiresIn: 3 * time.Minute,
}

func (store *TotalMemoryStore) Allow(identifier string) (allow bool, err error) {
	allow = store.limiter.Allow()
	if allow == true {
		return allow, nil
	} else {
		err = errors.New("Total rate limiter denied")
		return allow, err
	}
}
