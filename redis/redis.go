package redis

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/healthcheck"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
	"github.com/rosaekapratama/go-starter/utils"
	"strconv"
	"strings"
	"time"
)

const (
	pong         = "PONG"
	modeSingle   = "single"
	modeSentinel = "sentinel"
)

var (
	Client redis.UniversalClient
	Locker ILocker
)

// Init Initiate redis client with given configuration
func Init(ctx context.Context, config config.Config) {
	cfg := config.GetObject().Redis
	if cfg == nil || cfg.Disabled {
		log.Warn(ctx, "Redis client is disabled")
		return
	}

	if nil == cfg {
		log.Fatal(ctx, response.InitFailed, "Missing redis configuration")
		return
	}

	// Init redis client based on its mode
	mode := strings.ToLower(cfg.Mode)
	if mode == modeSingle {
		if strings.TrimSpace(cfg.Addr) == str.Empty {
			log.Fatal(ctx, response.ConfigNotFound, "Missing redis address")
			return
		}
		err := initSingleMode(ctx, cfg)
		if err != nil {
			log.Fatal(ctx, err, "Redis single mode init failed")
			return
		}
		log.Info(ctx, "Redis client is initiated in single mode")
	} else if mode == modeSentinel {
		if len(cfg.SentinelAddrs) == integer.Zero {
			log.Fatal(ctx, response.ConfigNotFound, "Missing redis sentinel addresses")
			return
		}
		err := initSentinelMode(ctx, cfg)
		if err != nil {
			log.Fatal(ctx, err, "Redis sentinel mode init failed")
			return
		}
		log.Info(ctx, "Redis client is initiated in sentinel mode")
	} else {
		log.Fatalf(ctx, response.InvalidConfig, "Unsupported mode '%s', valid mode are single or sentinel", mode)
		return
	}

	if Client == nil {
		log.Fatal(ctx, response.InitFailed, "Redis is not initiated")
		return
	}

	// Create a new distributed lock client.
	Locker = redislock.New(Client)

	// Add health checker
	healthcheck.AddChecker("redis", func(ctx context.Context) error {
		res, err := Client.Ping(ctx).Result()
		if err != nil {
			return err
		}
		if res != pong {
			return fmt.Errorf("ping %s", res)
		}
		return nil
	})
}

func initSingleMode(ctx context.Context, cfg *config.RedisConfig) error {
	singleConfig := cfg
	Client = redis.NewClient(&redis.Options{
		Network:         singleConfig.Network,
		Addr:            singleConfig.Addr,
		Username:        singleConfig.Username,
		Password:        singleConfig.Password,
		DB:              singleConfig.DB,
		MaxRetries:      singleConfig.MaxRetries,
		MinRetryBackoff: singleConfig.MinRetryBackoff.Duration,
		MaxRetryBackoff: singleConfig.MaxRetryBackoff.Duration,
		DialTimeout:     singleConfig.DialTimeout.Duration,
		ReadTimeout:     singleConfig.ReadTimeout.Duration,
		WriteTimeout:    singleConfig.WriteTimeout.Duration,
		PoolFIFO:        singleConfig.PoolFIFO,
		PoolSize:        singleConfig.PoolSize,
		PoolTimeout:     singleConfig.PoolTimeout.Duration,
		MinIdleConns:    singleConfig.MinIdleConns,
		MaxIdleConns:    singleConfig.MaxIdleConns,
		ConnMaxIdleTime: singleConfig.ConnMaxIdleTime.Duration,
		ConnMaxLifetime: singleConfig.ConnMaxLifetime.Duration,
	})
	ping, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Error(ctx, err, "Ping failed")
		return err
	}
	log.Trace(ctx, "Redis ping status:", ping)
	return nil
}

func initSentinelMode(ctx context.Context, cfg *config.RedisConfig) error {
	sentinelConfig := cfg
	Client = redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:              sentinelConfig.MasterName,
		SentinelAddrs:           sentinelConfig.SentinelAddrs,
		SentinelUsername:        sentinelConfig.SentinelUsername,
		SentinelPassword:        sentinelConfig.SentinelPassword,
		RouteByLatency:          sentinelConfig.RouteByLatency,
		RouteRandomly:           sentinelConfig.RouteRandomly,
		ReplicaOnly:             sentinelConfig.ReplicaOnly,
		UseDisconnectedReplicas: sentinelConfig.UseDisconnectedReplicas,
		Username:                sentinelConfig.Username,
		Password:                sentinelConfig.Password,
		DB:                      sentinelConfig.DB,
		MaxRetries:              sentinelConfig.MaxRetries,
		MinRetryBackoff:         sentinelConfig.MinRetryBackoff.Duration,
		MaxRetryBackoff:         sentinelConfig.MaxRetryBackoff.Duration,
		DialTimeout:             sentinelConfig.DialTimeout.Duration,
		ReadTimeout:             sentinelConfig.ReadTimeout.Duration,
		WriteTimeout:            sentinelConfig.WriteTimeout.Duration,
		PoolFIFO:                sentinelConfig.PoolFIFO,
		PoolSize:                sentinelConfig.PoolSize,
		PoolTimeout:             sentinelConfig.PoolTimeout.Duration,
		MinIdleConns:            sentinelConfig.MinIdleConns,
		MaxIdleConns:            sentinelConfig.MaxIdleConns,
		ConnMaxIdleTime:         sentinelConfig.ConnMaxIdleTime.Duration,
		ConnMaxLifetime:         sentinelConfig.ConnMaxLifetime.Duration,
	})
	ping, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Error(ctx, err, "Ping failed")
		return err
	}
	log.Trace(ctx, "Redis ping status:", ping)
	return nil
}

func (i Int) MarshalBinary() ([]byte, error) {
	return []byte(strconv.Itoa(int(i))), nil
}

func (p *Int) UnmarshalBinary(data []byte) error {
	v, err := strconv.Atoi(string(data))
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Int(v)
	return nil
}

func (i Int) MarshalText() (text []byte, err error) {
	return i.MarshalBinary()
}

func (p *Int) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (u Int8) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(u), integer.Ten)), nil
}

func (p *Int8) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseInt(string(data), integer.Ten, integer.Eight)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Int8(v)
	return nil
}

func (i Int8) MarshalText() (text []byte, err error) {
	return i.MarshalBinary()
}

func (p *Int8) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (i Int16) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), integer.Ten)), nil
}

func (p *Int16) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseInt(string(data), integer.Ten, integer.I16)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Int16(v)
	return nil
}

func (i Int16) MarshalText() (text []byte, err error) {
	return i.MarshalBinary()
}

func (p *Int16) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (i Int16s) MarshalBinary() ([]byte, error) {
	return []byte(utils.SliceOfInt16ToString(i, sym.Comma)), nil
}

func (p *Int16s) UnmarshalBinary(data []byte) error {
	uints, err := utils.StringToSliceOfInt16(string(data), sym.Comma)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = uints
	return nil
}

func (i Int16s) MarshalText() (text []byte, err error) {
	return i.MarshalBinary()
}

func (p *Int16s) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (i Int32) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), integer.Ten)), nil
}

func (p *Int32) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseInt(string(data), integer.Ten, integer.I32)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Int32(v)
	return nil
}

func (i Int32) MarshalText() (text []byte, err error) {
	return i.MarshalBinary()
}

func (p *Int32) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (i Int64) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), integer.Ten)), nil
}

func (p *Int64) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseInt(string(data), integer.Ten, integer.I64)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Int64(v)
	return nil
}

func (i Int64) MarshalText() (text []byte, err error) {
	return i.MarshalBinary()
}

func (p *Int64) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (u Uint8) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(u), integer.Ten)), nil
}

func (p *Uint8) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseInt(string(data), integer.Ten, integer.Eight)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Uint8(v)
	return nil
}

func (u Uint8) MarshalText() (text []byte, err error) {
	return u.MarshalBinary()
}

func (p *Uint8) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}
func (u Uint16) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(u), integer.Ten)), nil
}

func (p *Uint16) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseInt(string(data), integer.Ten, integer.I16)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Uint16(v)
	return nil
}

func (u Uint16) MarshalText() (text []byte, err error) {
	return u.MarshalBinary()
}

func (p *Uint16) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}
func (u Uint32) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(u), integer.Ten)), nil
}

func (p *Uint32) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseInt(string(data), integer.Ten, integer.I32)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Uint32(v)
	return nil
}

func (u Uint32) MarshalText() (text []byte, err error) {
	return u.MarshalBinary()
}

func (p *Uint32) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (u Uint64) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(u), integer.Ten)), nil
}

func (p *Uint64) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseUint(string(data), integer.Ten, integer.I64)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Uint64(v)
	return nil
}

func (u Uint64) MarshalText() (text []byte, err error) {
	return u.MarshalBinary()
}

func (p *Uint64) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (u Uint64s) MarshalBinary() ([]byte, error) {
	return []byte(utils.SliceOfUint64ToString(u, sym.Comma)), nil
}

func (p *Uint64s) UnmarshalBinary(data []byte) error {
	uints, err := utils.StringToSliceOfUint64(string(data), sym.Comma)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = uints
	return nil
}

func (u Uint64s) MarshalText() (text []byte, err error) {
	return u.MarshalBinary()
}

func (p *Uint64s) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (f Float64) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f), fmtF, integer.Two, integer.I64)), nil
}

func (p *Float64) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseFloat(string(data), integer.I64)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Float64(v)
	return nil
}

func (f Float64) MarshalText() (text []byte, err error) {
	return f.MarshalBinary()
}

func (p *Float64) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (f Float64s) MarshalBinary() ([]byte, error) {
	return []byte(utils.SliceOfFloat64ToString(f, sym.Comma)), nil
}

func (p *Float64s) UnmarshalBinary(data []byte) error {
	uints, err := utils.StringToSliceOfFloat64(string(data), sym.Comma)
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = uints
	return nil
}

func (f Float64s) MarshalText() (text []byte, err error) {
	return f.MarshalBinary()
}

func (p *Float64s) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (b Bool) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatBool(bool(b))), nil
}

func (p *Bool) UnmarshalBinary(data []byte) error {
	v, err := strconv.ParseBool(string(data))
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Bool(v)
	return nil
}

func (b Bool) MarshalText() (text []byte, err error) {
	return b.MarshalBinary()
}

func (p *Bool) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

func (t Time) MarshalBinary() ([]byte, error) {
	return []byte(time.Time(t).Format(time.RFC3339)), nil
}

func (p *Time) UnmarshalBinary(data []byte) error {
	v, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		log.Error(context.Background(), err)
		return err
	}
	*p = Time(v)
	return nil
}

func (t Time) MarshalText() (text []byte, err error) {
	return t.MarshalBinary()
}

func (p *Time) UnmarshalText(text []byte) error {
	return p.UnmarshalBinary(text)
}

// BoolP returns a pointer of a boolean variable
func BoolP(value Bool) *Bool {
	return &value
}

// PBool returns a boolean value from a pointer
func PBool(value *Bool) Bool {
	if value == nil {
		return false
	}
	return *value
}

// IntP returns a pointer of an integer variable
func IntP(value Int) *Int {
	return &value
}

// PInt returns an integer value from a pointer
func PInt(value *Int) Int {
	if value == nil {
		return 0
	}
	return *value
}

// Int8P returns a pointer of an int8 variable
func Int8P(value Int8) *Int8 {
	return &value
}

// PInt8 returns an int8 value from a pointer
func PInt8(value *Int8) Int8 {
	if value == nil {
		return 0
	}
	return *value
}

// Int16P returns a pointer of an int8 variable
func Int16P(value Int16) *Int16 {
	return &value
}

// PInt16 returns an int8 value from a pointer
func PInt16(value *Int16) Int16 {
	if value == nil {
		return 0
	}
	return *value
}

// Int32P returns a pointer of an Int32 variable
func Int32P(value Int32) *Int32 {
	return &value
}

// PInt32 returns an Int32 value from a pointer
func PInt32(value *Int32) Int32 {
	if value == nil {
		return 0
	}
	return *value
}

// Int64P returns a pointer of an Int64 variable
func Int64P(value Int64) *Int64 {
	return &value
}

// PInt64 returns an Int64 value from a pointer
func PInt64(value *Int64) Int64 {
	if value == nil {
		return 0
	}
	return *value
}

// Uint8P returns a pointer of an uint8 variable
func Uint8P(value Uint8) *Uint8 {
	return &value
}

// PUint8 returns an uint8 value from a pointer
func PUint8(value *Uint8) Uint8 {
	if value == nil {
		return 0
	}
	return *value
}

// Uint16P returns a pointer of an uint8 variable
func Uint16P(value Uint16) *Uint16 {
	return &value
}

// PUint16 returns an uint8 value from a pointer
func PUint16(value *Uint16) Uint16 {
	if value == nil {
		return 0
	}
	return *value
}

// Uint32P returns a pointer of an uint8 variable
func Uint32P(value Uint32) *Uint32 {
	return &value
}

// PUint32 returns an uint8 value from a pointer
func PUint32(value *Uint32) Uint32 {
	if value == nil {
		return 0
	}
	return *value
}

// Uint64P returns a pointer of an uint64 variable
func Uint64P(value Uint64) *Uint64 {
	return &value
}

// PUint64 returns an uint64 value from a pointer
func PUint64(value *Uint64) Uint64 {
	if value == nil {
		return 0
	}
	return *value
}

// Float64P returns a pointer of a float64 variable
func Float64P(value Float64) *Float64 {
	return &value
}

// PFloat64 returns an flaot64 value from a pointer
func PFloat64(value *Float64) Float64 {
	if value == nil {
		return 0
	}
	return *value
}
