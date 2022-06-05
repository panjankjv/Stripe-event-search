package stripesearch

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	eventTypeCustomer      = "customer"
	eventTypePaymentMethod = "payment_method"
	eventTypePaymentIntent = "payment_intent"
)

var defaultOutputKeys = []string{
	"Customer",
	"CreatedTime",
	"EventType",
	"ID",
	"CardBrand",
	"CardLast4",
	"CardFingerprint",
	"Description",
	"Amount",
	"Captured",
	"AmountRefunded",
	"FailureCode",
	"RiskScore",
}

var withSearchMetaOutputKeys = append([]string{
	"search_meta_value",
}, defaultOutputKeys...)

func GetHeaders(metaKeys, hideLabels []string, hasSearchMeta bool) []string {
	tmpList := append(defaultOutputKeys, metaKeys...)
	if hasSearchMeta {
		tmpList = append(withSearchMetaOutputKeys, metaKeys...)
	}

	hideMap := make(map[string]struct{})
	for _, v := range hideLabels {
		hideMap[v] = struct{}{}
	}

	results := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		if _, ok := hideMap[v]; ok {
			continue
		}
		results = append(results, v)
	}
	return results
}

type EventData struct {
	EventType          string              `json:"event_type"`
	Customer           string              `json:"customer"`
	Metadata           map[string]string   `json:"metadata"`
	CreatedTime        time.Time           `json:"created_time"`
	PaymentMethodEvent *PaymentMethodEvent `json:"payment_method,omitempty"`
	PaymentIntentEvent *PaymentIntentEvent `json:"payment_intent,omitempty"`
	CustomerEvent      *CustomerEvent      `json:"customer_data,omitempty"`

	searchMetaValue string
}

func (d EventData) IsCustomerEvent() bool {
	return d.EventType == eventTypeCustomer
}

func (d EventData) IsPaymentMethodEvent() bool {
	return d.EventType == eventTypePaymentMethod
}

func (d EventData) IsPaymentIntentEvent() bool {
	return d.EventType == eventTypePaymentIntent
}

func (d *EventData) SetSearchMetaValue(s string) {
	d.searchMetaValue = s
}

func (d EventData) Output(delimiter string, labels []string) string {
	mapVal := make(map[string]string)
	switch {
	case d.IsCustomerEvent():
		mapVal = toMapData(*d.CustomerEvent)
	case d.IsPaymentMethodEvent():
		mapVal = toMapData(*d.PaymentMethodEvent)
	case d.IsPaymentIntentEvent():
		mapVal = toMapData(*d.PaymentIntentEvent)
	}
	for _, key := range labels {
		if v, ok := d.Metadata[key]; ok {
			mapVal[key] = v
		}
	}
	mapVal["Customer"] = d.Customer
	mapVal["CreatedTime"] = d.CreatedTime.Format(time.RFC3339)
	mapVal["EventType"] = d.EventType

	outputs := make([]string, 0, len(mapVal))
	if d.searchMetaValue != "" {
		outputs = append(outputs, d.searchMetaValue)
		labels = labels[1:]
	}

	for _, v := range labels {
		outputs = append(outputs, mapVal[v])
	}
	return strings.Join(outputs, delimiter)
}

type EventSortAsc []EventData

func (e EventSortAsc) Len() int {
	return len(e)
}

func (e EventSortAsc) Less(i, j int) bool {
	return e[i].CreatedTime.Before(e[j].CreatedTime)
}

func (e EventSortAsc) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type CustomerEvent struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type PaymentMethodEvent struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	CardBrand       string `json:"card_brand"`
	CardFingerprint string `json:"card_fingerprint"`
	CardLast4       string `json:"card_last4"`
}

type PaymentIntentEvent struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Amount         int64  `json:"amount"`
	AmountRefunded int64  `json:"amount_refunded"`
	BillingEmail   string `json:"billing_email"`
	BillingName    string `json:"billing_name"`
	BillingPhone   string `json:"billing_phone"`
	Description    string `json:"description"`
	Captured       bool   `json:"captured"`
	Currency       string `json:"currency"`
	FailureCode    string `json:"failure_code"`
	RiskScore      int64  `json:"risk_score"`

	// card data
	CardBrand       string `json:"card_brand"`
	CardFingerprint string `json:"card_fingerprint"`
	CardLast4       string `json:"card_last4"`
	CVCCheck        string `json:"cvc_check"`
}

func toMapData(v interface{}) map[string]string {
	mapVal := make(map[string]string)

	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	for i, max := 0, vt.NumField(); i < max; i++ {
		field := vt.Field(i)
		key := field.Name
		value := fmt.Sprint(vv.FieldByName(key).Interface())
		mapVal[key] = value
	}
	return mapVal
}
