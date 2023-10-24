package redis

import (
	"context"
	"github.com/bsm/redislock"
	"time"
)

const fmtF = 'f'

type Int int
type Int8 int8
type Int16 int16
type Int16s []int16
type Int32 int32
type Int64 int64
type Uint8 uint8
type Uint16 uint16
type Uint32 uint32
type Uint64 uint64
type Uint64s []uint64
type Float64 float64
type Float64s []float64
type Bool bool
type Time time.Time

type ILocker interface {
	Obtain(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (*redislock.Lock, error)
}

type ILock interface {
	// Key returns the redis key used by the lock.
	Key() string

	// Token returns the token value set by the lock.
	Token() string

	// Metadata returns the metadata of the lock.
	Metadata() string

	// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
	TTL(ctx context.Context) (time.Duration, error)

	// Refresh extends the lock with a new TTL.
	// May return ErrNotObtained if refresh is unsuccessful.
	Refresh(ctx context.Context, ttl time.Duration, opt *redislock.Options) error

	// Release manually releases the lock.
	// May return ErrLockNotHeld.
	Release(ctx context.Context) error
}
