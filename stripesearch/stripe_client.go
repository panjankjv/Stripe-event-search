package stripesearch

import "github.com/stripe/stripe-go/v72/client"

const (
	defaultBaseURL = "https://api.stripe.com"
)

type StripeClientConfig struct {
	APIKey string
	Logger Logger
	Debug  bool
}

type StripeClient struct {
	RESTClient     RESTClient
	OfficialClient *client.API
	Logger         Logger
}

func NewStripeClient(conf StripeClientConfig) StripeClient {
	logger := conf.Logger
	if logger == nil {
		logger = &DummyLogger{}
	}
	cli := &client.API{}
	cli.Init(conf.APIKey, nil)
	return StripeClient{
		Logger:         logger,
		OfficialClient: cli,
		RESTClient: RESTClient{
			BasicAuthUser: conf.APIKey,
			Option: Option{
				Debug:   conf.Debug,
				BaseURL: defaultBaseURL,
			},
		},
	}
}

func (c StripeClient) LogInfo(format string, v ...interface{}) {
	c.Logger.Infof(format, v...)
}

func (c StripeClient) LogError(format string, v ...interface{}) {
	c.Logger.Errorf(format, v...)
}
