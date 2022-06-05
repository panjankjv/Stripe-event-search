package stripesearch

import (
	"fmt"
	"strings"
	"time"
)

func (c StripeClient) SearchPaymentIntent(opt SearchOption) ([]EventData, error) {
	resp, err := c.RawSearchPaymentIntent(opt)
	if err != nil {
		return nil, err
	}
	if resp.HasError() {
		return nil, resp.Error()
	}
	if len(resp.Data) == 0 {
		return nil, nil
	}

	list := make([]EventData, len(resp.Data))
	for i, v := range resp.Data {
		d := PaymentIntentEvent{
			ID:          v.ID,
			Amount:      v.Amount,
			Description: v.Description,
			Currency:    v.Currency,
		}
		ch, _ := v.Charges.GetFirstCharge()
		d.AmountRefunded = ch.AmountRefunded
		d.BillingEmail = ch.BillingDetails.Email
		d.BillingName = ch.BillingDetails.Name
		d.BillingPhone = ch.BillingDetails.Phone
		d.Captured = ch.Captured
		d.FailureCode = ch.FailureCode
		d.RiskScore = ch.Outcome.RsikScore

		card := ch.PaymentMethodDetails.Card
		d.CardBrand = card.Brand
		d.CardFingerprint = card.Fingerprint
		d.CardLast4 = card.Last4
		d.CVCCheck = card.Checks.CVCCheck

		ev := EventData{
			EventType:          eventTypePaymentIntent,
			Customer:           v.Customer,
			Metadata:           v.Metadata,
			CreatedTime:        time.Unix(v.Created, 0),
			PaymentIntentEvent: &d,
		}
		list[i] = ev
	}
	return list, nil
}

func (c StripeClient) RawSearchPaymentIntent(opt SearchOption) (PaymentIntentResponse, error) {
	params := make(map[string]string)
	params["query"] = opt.ToQueryParameter()

	resp := PaymentIntentResponse{}
	code, err := c.RESTClient.CallGET("/v1/search/payment_intents", params, &resp)
	resp.StatusCode = code
	return resp, err
}

func (c StripeClient) SearchCustomer(opt SearchOption) ([]EventData, error) {
	resp, err := c.RawSearchCustomer(opt)
	if err != nil {
		return nil, err
	}
	if resp.HasError() {
		return nil, resp.Error()
	}
	if len(resp.Data) == 0 {
		return nil, nil
	}

	list := make([]EventData, len(resp.Data))
	for i, v := range resp.Data {
		ev := EventData{
			EventType:   eventTypeCustomer,
			Customer:    v.ID,
			Metadata:    v.Metadata,
			CreatedTime: time.Unix(v.Created, 0),
			CustomerEvent: &CustomerEvent{
				ID:          v.ID,
				Description: v.Description,
				Name:        v.Name,
				Email:       v.Email,
				Phone:       v.Phone,
			},
		}
		list[i] = ev
	}
	return list, nil
}

func (c StripeClient) RawSearchCustomer(opt SearchOption) (CustomerResponse, error) {
	params := make(map[string]string)
	params["query"] = opt.ToQueryParameter()

	resp := CustomerResponse{}
	code, err := c.RESTClient.CallGET("/v1/search/customers", params, &resp)
	resp.StatusCode = code
	return resp, err
}

type SearchOption struct {
	// common
	Metadata map[string]string

	// for payment intent
	Customer string

	// for customer
	Email string
}

func (o SearchOption) ToQueryParameter() string {
	list := []string{}
	if o.Customer != "" {
		list = append(list, fmt.Sprintf("customer:'%s'", o.Customer))
	}
	if o.Email != "" {
		list = append(list, fmt.Sprintf("email:'%s'", o.Email))
	}
	for key, val := range o.Metadata {
		list = append(list, fmt.Sprintf("metadata['%s']:'%s'", key, val))
	}
	return strings.Join(list, " AND ")
}
