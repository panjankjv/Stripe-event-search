package stripesearch

import (
	"time"

	stripe "github.com/stripe/stripe-go/v72"
)

func (c StripeClient) GetCustomer(opt CustomerOption) (*EventData, error) {
	params := stripe.CustomerParams{}

	cust, err := c.OfficialClient.Customers.Get(opt.Customer, &params)
	if err != nil {
		return nil, err
	}
	return &EventData{
		EventType:   eventTypeCustomer,
		Customer:    cust.ID,
		CreatedTime: time.Unix(cust.Created, 0),
		Metadata:    cust.Metadata,
		CustomerEvent: &CustomerEvent{
			ID:          cust.ID,
			Description: cust.Description,
			Name:        cust.Name,
			Email:       cust.Email,
			Phone:       cust.Phone,
		},
	}, nil
}

type CustomerOption struct {
	Customer string
}
