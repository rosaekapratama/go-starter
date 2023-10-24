package location

import (
	"context"
	"github.com/rosaekapratama/go-starter/constant/timezone"
	"github.com/rosaekapratama/go-starter/log"
	"time"
)

var (
	AsiaJakarta *time.Location
)

func init() {
	ctx := context.Background()

	var err error
	AsiaJakarta, err = time.LoadLocation(timezone.AsiaJakarta)
	if err != nil {
		log.Fatal(ctx, err, "Failed to load location Asia/Jakarta")
		return
	}
}
