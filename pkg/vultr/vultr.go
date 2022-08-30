package vultr

import (
	"context"
	"time"

	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
)

// Config is the struct for govultr configuration
type Config struct {
	APIKey     string
	RateLimit  int
	RetryLimit int
}

// Init configures and returns an initialized govultr client
func (c *Config) Init() (*govultr.Client, error) {
	userAgent := "govultr"
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.APIKey,
	})

	client := oauth2.NewClient(context.Background(), tokenSrc)

	vultrClient := govultr.NewClient(client)
	vultrClient.SetUserAgent(userAgent)

	if c.RateLimit != 0 {
		vultrClient.SetRateLimit(time.Duration(c.RateLimit) * time.Millisecond)
	}

	if c.RetryLimit != 0 {
		vultrClient.SetRetryLimit(c.RetryLimit)
	}

	return vultrClient, nil
}
