package stripesearch

import (
	"time"

	stripe "github.com/stripe/stripe-go/v72"
)

func (c StripeClient) GetPaymentMethodList(opt PaymentMethodOption) ([]EventData, error) {
	params := stripe.PaymentMethodListParams{
		Customer: stripe.String(opt.Customer),
		Type:     stripe.String(opt.Type),
	}

	iter := c.OfficialClient.PaymentMethods.List(&params)
	if err := iter.Err(); err != nil {
		return nil, err
	}

	var list []EventData
	for iter.Next() {
		pm := iter.PaymentMethod()
		card := newPaymentMethodCard(pm.Card)
		ev := EventData{
			EventType:   eventTypePaymentMethod,
			Metadata:    pm.Metadata,
			Customer:    pm.Customer.ID,
			CreatedTime: time.Unix(pm.Created, 0),
			PaymentMethodEvent: &PaymentMethodEvent{
				ID:              pm.ID,
				Type:            string(pm.Type),
				CardBrand:       card.Brand,
				CardLast4:       card.Last4,
				CardFingerprint: card.Fingerprint,
			},
		}
		list = append(list, ev)
	}
	return list, nil
}

type PaymentMethodOption struct {
	Customer string
	Type     string // card
}
