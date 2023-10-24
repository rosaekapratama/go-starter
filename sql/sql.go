package sql

import (
	"database/sql"
	"github.com/rosaekapratama/go-starter/constant/str"
	"time"
)

func GetNullString(val string) sql.NullString {
	if val != str.Empty {
		return sql.NullString{
			String: val,
			Valid:  true,
		}
	} else {
		return sql.NullString{
			Valid: false,
		}
	}
}

func GetNullTime(val time.Time) sql.NullTime {
	if val.IsZero() {
		return sql.NullTime{
			Valid: false,
		}
	} else {
		return sql.NullTime{
			Time:  val,
			Valid: true,
		}
	}
}
