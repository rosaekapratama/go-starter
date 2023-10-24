package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/redis"
	"github.com/rosaekapratama/go-starter/response"
	"strconv"
)

const (
	spanGet                           = "cache.user.Get"
	spanGetFullname                   = "cache.user.GetFullname"
	spanGetFullnameWithDefaultValue   = "cache.user.GetFullnameWithDefaultValue"
	spanGetUsername                   = "cache.user.GetUsername"
	spanGetUsernameWithDefaultValue   = "cache.user.GetUsernameWithDefaultValue"
	spanGetEmail                      = "cache.user.GetEmail"
	spanGetEmailWithDefaultValue      = "cache.user.GetEmailWithDefaultValue"
	spanGetRealm                      = "cache.user.GetRealm"
	spanGetRealmWithDefaultValue      = "cache.user.GetRealmWithDefaultValue"
	spanGetTerminalId                 = "cache.user.GetTerminalId"
	spanGetTerminalIdWithDefaultValue = "cache.user.GetTerminalIdWithDefaultValue"
	spanGetProviderId                 = "cache.user.GetProviderId"
	spanGetProviderIdWithDefaultValue = "cache.user.GetProviderIdWithDefaultValue"
	spanGetIsActive                   = "cache.user.GetIsActive"
	spanGetIsActiveWithDefaultValue   = "cache.user.GetIsActiveWithDefaultValue"
	spanGetRoleGrade                  = "cache.user.GetRoleGrade"
	spanGetRoleGradeWithDefaultValue  = "cache.user.GetRoleGradeWithDefaultValue"
	spanSet                           = "cache.user.Set"
	spanSetWithTx                     = "cache.user.SetWithTx"
	spanDel                           = "cache.user.Del"
	spanDelWithTx                     = "cache.user.DelWithTx"
	spanDelAll                        = "cache.user.DelAll"
	spanLen                           = "cache.user.Len"

	prefixUser      = "user:%s"
	fullnameField   = "fullname"
	usernameField   = "username"
	emailField      = "email"
	realmField      = "realm"
	terminalIdField = "terminalId"
	providerIdField = "providerId"
	isActiveField   = "isActive"
	roleGradeField  = "roleGrade"
)

// Get will return nil if no data found in redis
func Get(ctx context.Context, uuid string) (*User, bool) {
	ctx, span := otel.Trace(ctx, spanGet)
	defer span.End()

	m, err := redis.Client.HGetAll(ctx, fmt.Sprintf(prefixUser, uuid)).Result()
	if errors.Is(err, goredis.Nil) || len(m) == integer.Zero {
		log.Tracef(ctx, "User cache not found, uuid=%s", uuid)
		return nil, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user from cache, uuid=%s", uuid)
		return nil, false
	}

	u := User{}
	mds, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.TextUnmarshallerHookFunc(),
		Result:     &u,
	})
	if err != nil {
		log.Errorf(ctx, err, "Failed to create decoder for user cache, uuid=%s", uuid)
		return nil, false
	}

	err = mds.Decode(m)
	if err != nil {
		log.Errorf(ctx, err, "Failed to decode user map to struct from redis, uuid=%s", uuid)
		return nil, false
	}
	return &u, true
}

// GetFullname will return string empty if no data or fullname found in redis
func GetFullname(ctx context.Context, uuid string) (string, bool) {
	ctx, span := otel.Trace(ctx, spanGetFullname)
	defer span.End()

	name, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), fullnameField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User fullname cache not found, uuid=%s", uuid)
		return str.Empty, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user fullname from cache, uuid=%s", uuid)
		return str.Empty, false
	}
	return name, true
}

// GetFullnameWithDefaultValue will return default value if no data or fullname found in redis
func GetFullnameWithDefaultValue(ctx context.Context, uuid string, defaultValue string) string {
	ctx, span := otel.Trace(ctx, spanGetFullnameWithDefaultValue)
	defer span.End()

	if name, ok := GetFullname(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

// GetUsername will return string empty if no data or username found in redis
func GetUsername(ctx context.Context, uuid string) (string, bool) {
	ctx, span := otel.Trace(ctx, spanGetUsername)
	defer span.End()

	name, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), usernameField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User username cache not found, uuid=%s", uuid)
		return str.Empty, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user username from cache, uuid=%s", uuid)
		return str.Empty, false
	}
	return name, true
}

// GetUsernameWithDefaultValue will return default value if no data or username found in redis
func GetUsernameWithDefaultValue(ctx context.Context, uuid string, defaultValue string) string {
	ctx, span := otel.Trace(ctx, spanGetUsernameWithDefaultValue)
	defer span.End()

	if name, ok := GetUsername(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

// GetEmail will return string empty if no data or email found in redis
func GetEmail(ctx context.Context, uuid string) (string, bool) {
	ctx, span := otel.Trace(ctx, spanGetEmail)
	defer span.End()

	name, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), emailField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User email cache not found, uuid=%s", uuid)
		return str.Empty, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user email from cache, uuid=%s", uuid)
		return str.Empty, false
	}
	return name, true
}

// GetEmailWithDefaultValue will return default value if no data or email found in redis
func GetEmailWithDefaultValue(ctx context.Context, uuid string, defaultValue string) string {
	ctx, span := otel.Trace(ctx, spanGetEmailWithDefaultValue)
	defer span.End()

	if name, ok := GetEmail(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

// GetRealm will return string empty if no data or realm found in redis
func GetRealm(ctx context.Context, uuid string) (string, bool) {
	ctx, span := otel.Trace(ctx, spanGetRealm)
	defer span.End()

	name, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), realmField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User realm cache not found, uuid=%s", uuid)
		return str.Empty, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user realm from cache, uuid=%s", uuid)
		return str.Empty, false
	}
	return name, true
}

// GetRealmWithDefaultValue will return default value if no data or realm found in redis
func GetRealmWithDefaultValue(ctx context.Context, uuid string, defaultValue string) string {
	ctx, span := otel.Trace(ctx, spanGetRealmWithDefaultValue)
	defer span.End()

	if name, ok := GetRealm(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

// GetTerminalId will return string empty if no data or realm found in redis
func GetTerminalId(ctx context.Context, uuid string) (string, bool) {
	ctx, span := otel.Trace(ctx, spanGetTerminalId)
	defer span.End()

	name, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), terminalIdField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User terminal ID cache not found, uuid=%s", uuid)
		return str.Empty, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user terminal ID from cache, uuid=%s", uuid)
		return str.Empty, false
	}
	return name, true
}

// GetTerminalIdWithDefaultValue will return default value if no data or realm found in redis
func GetTerminalIdWithDefaultValue(ctx context.Context, uuid string, defaultValue string) string {
	ctx, span := otel.Trace(ctx, spanGetTerminalIdWithDefaultValue)
	defer span.End()

	if name, ok := GetTerminalId(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

// GetProviderId will return string empty if no data or realm found in redis
func GetProviderId(ctx context.Context, uuid string) (string, bool) {
	ctx, span := otel.Trace(ctx, spanGetProviderId)
	defer span.End()

	name, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), providerIdField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User provider ID cache not found, uuid=%s", uuid)
		return str.Empty, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user provider ID from cache, uuid=%s", uuid)
		return str.Empty, false
	}
	return name, true
}

// GetProviderIdWithDefaultValue will return default value if no data or realm found in redis
func GetProviderIdWithDefaultValue(ctx context.Context, uuid string, defaultValue string) string {
	ctx, span := otel.Trace(ctx, spanGetProviderIdWithDefaultValue)
	defer span.End()

	if name, ok := GetProviderId(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

// GetIsActive will return string empty if no data or realm found in redis
func GetIsActive(ctx context.Context, uuid string) (string, bool) {
	ctx, span := otel.Trace(ctx, spanGetIsActive)
	defer span.End()

	name, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), isActiveField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User is active cache not found, uuid=%s", uuid)
		return str.Empty, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user is active from cache, uuid=%s", uuid)
		return str.Empty, false
	}
	return name, true
}

// GetIsActiveWithDefaultValue will return default value if no data or realm found in redis
func GetIsActiveWithDefaultValue(ctx context.Context, uuid string, defaultValue string) string {
	ctx, span := otel.Trace(ctx, spanGetIsActiveWithDefaultValue)
	defer span.End()

	if name, ok := GetIsActive(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

// GetRoleGrade will return integer zero if no data in redis
func GetRoleGrade(ctx context.Context, uuid string) (int, bool) {
	ctx, span := otel.Trace(ctx, spanGetRoleGrade)
	defer span.End()

	roleGrade, err := redis.Client.HGet(ctx, fmt.Sprintf(prefixUser, uuid), roleGradeField).Result()
	if errors.Is(err, goredis.Nil) {
		log.Tracef(ctx, "User role grade cache not found, uuid=%s", uuid)
		return integer.Zero, false
	} else if err != nil {
		log.Errorf(ctx, err, "Failed to get user role grade from cache, uuid=%s", uuid)
		return integer.Zero, false
	}

	v, err := strconv.Atoi(roleGrade)
	if err != nil {
		log.Errorf(ctx, err, "Failed to parse user role grade to int, uuid=%s", uuid)
		return integer.Zero, false
	}

	return v, true
}

// GetRoleGradeWithDefaultValue will return default value if no data or realm found in redis
func GetRoleGradeWithDefaultValue(ctx context.Context, uuid string, defaultValue int) int {
	ctx, span := otel.Trace(ctx, spanGetRoleGradeWithDefaultValue)
	defer span.End()

	if name, ok := GetRoleGrade(ctx, uuid); ok {
		return name
	}
	return defaultValue
}

func Set(ctx context.Context, u *User) error {
	ctx, span := otel.Trace(ctx, spanSet)
	defer span.End()

	if u.Uuid == str.Empty {
		log.Errorf(ctx, response.InvalidArgument, "Set user cache failed, UUID must not empty, %v", u)
		return response.InvalidArgument
	}

	temp := structs.Map(u)
	err := redis.Client.HSet(ctx, fmt.Sprintf(prefixUser, u.Uuid), temp).Err()
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

func SetWithTx(ctx context.Context, pipe goredis.Pipeliner, u *User) error {
	ctx, span := otel.Trace(ctx, spanSetWithTx)
	defer span.End()

	if u.Uuid == str.Empty {
		log.Errorf(ctx, response.InvalidArgument, "Set user cache failed, UUID must not empty, %v", u)
		return response.InvalidArgument
	}

	temp := structs.Map(u)
	err := pipe.HSet(ctx, fmt.Sprintf(prefixUser, u.Uuid), temp).Err()
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

func Del(ctx context.Context, uuid string) error {
	ctx, span := otel.Trace(ctx, spanDel)
	defer span.End()

	err := redis.Client.Del(ctx, fmt.Sprintf(prefixUser, uuid)).Err()
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

func DelWithTx(ctx context.Context, pipe goredis.Pipeliner, uuid string) error {
	ctx, span := otel.Trace(ctx, spanDelWithTx)
	defer span.End()

	err := pipe.Del(ctx, fmt.Sprintf(prefixUser, uuid)).Err()
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// DelAll will delete all caches related to this package
func DelAll(ctx context.Context) error {
	ctx, span := otel.Trace(ctx, spanDelAll)
	defer span.End()

	pattern := fmt.Sprintf(prefixUser, sym.Asterisk)
	keys, err := redis.Client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Errorf(ctx, err, "Failed to get user cache keys, prefix=%s", pattern)
		return err
	}

	if len(keys) > integer.Zero {
		err = redis.Client.Del(ctx, keys...).Err()
		if err != nil && errors.Is(err, goredis.Nil) {
			log.Tracef(ctx, "User caches are not found, prefix=%s", pattern)
			return nil
		} else if err != nil {
			log.Errorf(ctx, err, "Failed to delete user caches, prefix=%s", pattern)
			return err
		}
	}

	return nil
}

func Len(ctx context.Context) (int, error) {
	ctx, span := otel.Trace(ctx, spanLen)
	defer span.End()

	pattern := fmt.Sprintf(prefixUser, sym.Asterisk)
	keys, err := redis.Client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Errorf(ctx, err, "Failed to get user cache keys, prefix=%s", pattern)
		return integer.Zero, err
	}

	return len(keys), nil
}
