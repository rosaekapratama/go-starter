package user

import "github.com/rosaekapratama/go-starter/redis"

type User struct {
	Uuid      string      `mapstructure:"uuid,omitempty" structs:"uuid,omitempty"`
	Username  string      `mapstructure:"username,omitempty" structs:"username,omitempty"`
	Email     string      `mapstructure:"email,omitempty" structs:"email,omitempty"`
	Fullname  string      `mapstructure:"fullname,omitempty" structs:"fullname,omitempty"`
	Realm     string      `mapstructure:"realm,omitempty" structs:"realm,omitempty"`
	IsActive  *redis.Bool `mapstructure:"isActive,omitempty" structs:"isActive,omitempty"`
	RoleGrade *redis.Int  `mapstructure:"roleGrade,omitempty" structs:"roleGrade,omitempty"`

	// For user provider realm
	TerminalId string `mapstructure:"terminalId,omitempty" structs:"terminalId,omitempty"`
	ProviderId string `mapstructure:"providerId,omitempty" structs:"providerId,omitempty"`
}

type LoadUserRequest struct {
	Id          string `avro:"id"`
	Realm       string `avro:"realm"`
	CreatedDate string `avro:"createdDate"`
}
