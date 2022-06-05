package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/evalphobia/stripe-event-search/stripesearch"
)

type runOption struct {
	APIKey              string
	Customer            string
	PaymentType         string
	ShowMetadataKey     []string
	HideLabels          []string
	SearchMetadataKey   string
	SearchMetadataValue string
	Debug               bool
	TimeAfter           time.Time
}

func getClient(opt runOption) stripesearch.StripeClient {
	apiKey := getAPIKey()
	if apiKey == "" {
		apiKey = opt.APIKey
	}
	return stripesearch.NewStripeClient(stripesearch.StripeClientConfig{
		APIKey: apiKey,
		Debug:  opt.Debug,
		Logger: &stripesearch.StdLogger{},
	})
}

func searchCustomers(cli stripesearch.StripeClient, metaKey, metaValue string) ([]stripesearch.EventData, error) {
	results, err := cli.SearchCustomer(stripesearch.SearchOption{
		Metadata: map[string]string{
			metaKey: metaValue,
		},
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(stripesearch.EventSortAsc(results))
	return results, nil
}

func searchPaymentIntentEvents(cli stripesearch.StripeClient, customer, paymentType string) ([]stripesearch.EventData, error) {
	var results []stripesearch.EventData
	events, err := cli.GetPaymentMethodList(stripesearch.PaymentMethodOption{
		Customer: customer,
		Type:     paymentType,
	})
	if err != nil {
		return nil, err
	}
	results = append(results, events...)

	events, err = cli.SearchPaymentIntent(stripesearch.SearchOption{
		Customer: customer,
	})
	if err != nil {
		return nil, err
	}
	results = append(results, events...)

	sort.Sort(stripesearch.EventSortAsc(results))
	return results, nil
}

// get all of the target events for the customer.
func fetchSingleCustomerAllEvents(cli stripesearch.StripeClient, customer, paymentType string, includeCustomer bool) ([]stripesearch.EventData, error) {
	var results []stripesearch.EventData
	if includeCustomer {
		event, err := cli.GetCustomer(stripesearch.CustomerOption{
			Customer: customer,
		})
		if err != nil {
			return nil, err
		}
		results = append(results, *event)
	}

	res, err := searchPaymentIntentEvents(cli, customer, paymentType)
	if err != nil {
		return nil, err
	}
	results = append(results, res...)
	return results, nil
}

// get all of the target events for the customer(s).
func fetchCustomersAllEvents(cli stripesearch.StripeClient, opt runOption) ([]stripesearch.EventData, error) {
	if opt.Customer != "" {
		results, err := fetchSingleCustomerAllEvents(cli, opt.Customer, opt.PaymentType, true)
		if err != nil {
			return nil, err
		}
		return filterEventsByTimeAfter(results, opt.TimeAfter), nil
	}

	customers, err := cli.SearchCustomer(stripesearch.SearchOption{
		Metadata: map[string]string{
			opt.SearchMetadataKey: opt.SearchMetadataValue,
		},
	})
	if err != nil {
		return nil, err
	}

	var results []stripesearch.EventData
	for _, c := range customers {
		res, err := fetchSingleCustomerAllEvents(cli, c.Customer, opt.PaymentType, false)
		if err != nil {
			return nil, err
		}

		results = append(results, c)
		results = append(results, res...)
	}
	return filterEventsByTimeAfter(results, opt.TimeAfter), nil
}

func filterEventsByTimeAfter(list []stripesearch.EventData, dt time.Time) []stripesearch.EventData {
	if dt.IsZero() {
		return list
	}

	results := make([]stripesearch.EventData, 0, len(list))
	for _, v := range list {
		switch {
		case v.IsCustomerEvent(),
			v.CreatedTime.After(dt):
			results = append(results, v)
		}
	}
	return results
}

func parseTime(s string) (time.Time, error) {
	if dt, err := time.Parse(time.RFC3339, s); err == nil {
		return dt, nil
	}
	if dt, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return dt, nil
	}
	if dt, err := time.Parse("2006-01-02", s); err == nil {
		return dt, nil
	}

	return time.Time{}, fmt.Errorf("cannot parse to time: [%s]", s)
}

func parseStringSlice(s string) []string {
	tmpList := strings.Split(s, " ")
	results := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		results = append(results, v)
	}
	return results
}
