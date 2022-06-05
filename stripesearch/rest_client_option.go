package stripesearch

import (
	"fmt"
	"time"
)

const (
	clientVersion  = "v0.0.1"
	defaultTimeout = 20 * time.Second
)

var (
	defaultUserAgent     = fmt.Sprintf("stripe-event-search/%s", clientVersion)
	defaultStripeVersion = "2020-08-27;search_api_beta=v1"
)

// Option contains optional setting of RESTClient.
type Option struct {
	BaseURL       string
	Headers       map[string]string
	UserAgent     string
	StripeVersion string
	Timeout       time.Duration
	Debug         bool
	Retry         bool
	LogFn         func(msg string, opts ...interface{})
}

func (o Option) LogInfo(msg string, opts ...interface{}) {
	if o.LogFn == nil {
		return
	}
	o.LogFn(msg, opts...)
}

func (o Option) getUserAgent() string {
	if o.UserAgent != "" {
		return o.UserAgent
	}
	return defaultUserAgent
}

func (o Option) getStripeVersion() string {
	if o.StripeVersion != "" {
		return o.StripeVersion
	}
	return defaultStripeVersion
}

func (o Option) getTimeout() time.Duration {
	if o.Timeout > 0 {
		return o.Timeout
	}
	return defaultTimeout
}

func (o Option) getHeaders() map[string]string {
	h := o.Headers
	if h == nil {
		h = make(map[string]string)
	}

	ver := o.getStripeVersion()
	if _, ok := h["Stripe-Version"]; !ok && ver != "" {
		h["Stripe-Version"] = ver
	}

	if len(h) != 0 {
		return h
	}
	return nil
}
