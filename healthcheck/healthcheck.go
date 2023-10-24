package healthcheck

import (
	"context"
	"github.com/etherlabsio/healthcheck/v2"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"net/http"
	"time"
)

const URLPathRegex = sym.Circumflex + "/v./health" + sym.Dollars

var options []healthcheck.Option

func AddChecker(name string, f func(ctx context.Context) error) {
	options = append(
		options,
		healthcheck.WithChecker(
			name,
			healthcheck.CheckerFunc(f),
		))
}

func HandlerV1() http.Handler {
	// WithTimeout allows you to set a max overall timeout.
	options = append(options, healthcheck.WithTimeout(5*time.Second))
	return healthcheck.Handler(options...)
}
