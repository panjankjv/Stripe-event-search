package stripesearch

import (
	"time"

	"github.com/stripe/stripe-go/v72"
)

type BaseResponse struct {
	// common
	StatusCode int `json:"status_code"`
	// error only
	Err ResponseError `json:"error"`
	// success only
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	URL     string `json:"url"`
}

func (r BaseResponse) HasError() bool {
	return r.Err.HasError()
}

func (r BaseResponse) Error() error {
	return r.Err.Error()
}

type CustomerResponse struct {
	BaseResponse
	Data []Customer `json:"data"`
}

type PaymentIntentResponse struct {
	BaseResponse
	Data []PaymentIntent `json:"data"`
}

type PaymentIntent struct {
	ID                        string            `json:"id"`
	Object                    string            `json:"object"`
	Amount                    int64             `json:"amount"`
	AmountCapturable          int64             `json:"amount_capturable"`
	AmountReceived            int64             `json:"amount_received"`
	Application               interface{}       `json:"application"`
	ApplicationFeeAmount      interface{}       `json:"application_fee_amount"`
	ApplicationPaymentMethods interface{}       `json:"application_payment_methods"`
	CanceledAt                interface{}       `json:"canceled_at"`
	CancellationReason        interface{}       `json:"cancellation_reason"`
	CaptureMethod             string            `json:"capture_method"`
	Charges                   Charges           `json:"charges"`
	ClientSecret              string            `json:"client_secret"`
	ConfitmationMethod        string            `json:"confirmation_method"`
	Created                   int64             `json:"created"`
	Currency                  string            `json:"currency"`
	Customer                  string            `json:"customer"`
	Description               string            `json:"description"`
	Invoice                   string            `json:"invoice"`
	LastPaymentError          interface{}       `json:"last_payment_error"`
	Livemode                  bool              `json:"livemode"`
	Metadata                  map[string]string `json:"metadata"`
	NextAction                interface{}       `json:"next_action"`
	OnBehalfOf                interface{}       `json:"on_behalf_of"`
	PaymentMethod             string            `json:"payment_method"`
	PaymentMethodOptions      interface{}       `json:"payment_method_options"`
	PaymentMethodTypes        []string          `json:"payment_method_types"`
	Processing                string            `json:"processing"`
	ReceiptEmail              string            `json:"receipt_email"`
	Review                    string            `json:"review"`
	SetupFutureUsage          string            `json:"setup_future_usage"`
	Shipping                  string            `json:"shipping"`
	Source                    string            `json:"source"`
	StatementDescriptor       string            `json:"statement_descriptor"`
	StatementDescriptorSuffix string            `json:"statement_descriptor_suffix"`
	Status                    string            `json:"status"`
	TransferData              string            `json:"transfer_data"`
	TransferGroup             string            `json:"transfer_group"`
}

type Charges struct {
	Object     string   `json:"object"`
	Data       []Charge `json:"data"`
	HasMore    bool     `json:"has_more"`
	TotalCount int64    `json:"total_count"`
	URL        string   `json:"url"`
}

func (c Charges) GetFirstCharge() (Charge, bool) {
	if len(c.Data) == 0 {
		return Charge{}, false
	}
	return c.Data[0], true
}

type Charge struct {
	ID                            string               `json:"id"`
	Object                        string               `json:"object"`
	Amount                        int64                `json:"amount"`
	AmountCaptured                int64                `json:"amount_captured"`
	AmountRefunded                int64                `json:"amount_refunded"`
	Application                   interface{}          `json:"application"`
	ApplicationFee                interface{}          `json:"application_fee"`
	ApplicationFeeAmount          interface{}          `json:"application_fee_amount"`
	BalanceTransaction            interface{}          `json:"balance_transaction"`
	BillingDetails                BillingDetails       `json:"billing_details"`
	CalculatedStatementDescriptor string               `json:"calculated_statement_descriptor"`
	Captured                      bool                 `json:"captured"`
	Created                       int64                `json:"created"`
	Currency                      string               `json:"currency"`
	Customer                      string               `json:"customer"`
	Description                   string               `json:"description"`
	Destination                   interface{}          `json:"destination"`
	Dispute                       interface{}          `json:"dispute"`
	Disputed                      bool                 `json:"disputed"`
	FailureCode                   string               `json:"failure_code"`
	FailureMessage                string               `json:"failure_message"`
	FraudDetails                  FraudDetails         `json:"fraud_details"`
	Invoice                       string               `json:"invoice"`
	Livemode                      bool                 `json:"livemode"`
	Metadata                      map[string]string    `json:"metadata"`
	OnBehalfOf                    interface{}          `json:"on_behalf_of"`
	Order                         interface{}          `json:"order"`
	Outcome                       Outcome              `json:"outcome"`
	Paid                          bool                 `json:"paid"`
	PaymentIntent                 string               `json:"payment_intent"`
	PaymentMethod                 string               `json:"payment_method"`
	PaymentMethodDetails          PaymentMethodDetails `json:"payment_method_details"`
	ReceiptEmail                  string               `json:"receipt_email"`
	ReceiptNumber                 string               `json:"receipt_number"`
	ReceiptURL                    string               `json:"receipt_url"`
	Refunded                      bool                 `json:"refunded"`
	Refunds                       Refunds              `json:"refunds"`
	Review                        interface{}          `json:"review"`
	Shipping                      interface{}          `json:"shipping"`
	Source                        string               `json:"source"`
	SourceTransfer                interface{}          `json:"source_transfer"`
	StatementDescriptor           string               `json:"statement_descriptor"`
	StatementDescriptorSuffix     string               `json:"statement_descriptor_suffix"`
	Status                        string               `json:"status"`
	TransferData                  interface{}          `json:"transfer_data"`
	TransferGroup                 string               `json:"transfer_group"`
}

type LastPaymentError struct {
	Charge        string        `json:"charge"`
	Code          string        `json:"code"`
	DeclineCode   string        `json:"decline_code"`
	DocURL        string        `json:"doc_url"`
	Message       string        `json:"message"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Type          string        `json:"type"`
}

type PaymentMethod struct {
	ID             string            `json:"id"`
	Object         string            `json:"object"`
	BillingDetails BillingDetails    `json:"billing_details"`
	Card           PaymentMethodCard `json:"card"`
	Created        int64             `json:"created"`
	Customer       string            `json:"customer"`
	Livemode       bool              `json:"livemode"`
	Metadata       map[string]string `json:"metadata"`
	Type           string            `json:"type"`

	// optional
	CreatedTime time.Time `json:"created_time"`
}

type BillingDetails struct {
	Address Address `json:"address"`
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
}

type Address struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

type PaymentMethodCard struct {
	Brand             string            `json:"brand"`
	Checks            CardChecks        `json:"checks"`
	Country           string            `json:"country"`
	ExpMonth          int64             `json:"exp_month"`
	ExpYear           int64             `json:"exp_year"`
	Fingerprint       string            `json:"fingerprint"`
	Funding           string            `json:"funding"`
	GeneratedFrom     string            `json:"generated_from"`
	Last4             string            `json:"last4"`
	Networks          interface{}       `json:"networks"`
	ThreeDSecureUsage ThreeDSecureUsage `json:"three_d_secure_usage"`
	Wallet            interface{}       `json:"wallet"`
}

func newPaymentMethodCard(card *stripe.PaymentMethodCard) PaymentMethodCard {
	if card == nil {
		return PaymentMethodCard{}
	}
	return PaymentMethodCard{
		Brand:       string(card.Brand),
		Fingerprint: card.Fingerprint,
		Last4:       card.Last4,
	}
}

type CardChecks struct {
	AddressLine1Check      string `json:"address_line1_check"`
	AddressPostalCodeCheck string `json:"address_postal_code_check"`
	CVCCheck               string `json:"cvc_check"`
}

type ThreeDSecureUsage struct {
	Supported bool `json:"supported"`
}

type PaymentMethodDetails struct {
	Card PaymentMethodCard `json:"card"`
}

type Outcome struct {
	NetworkStatus string `json:"network_status"`
	Reason        string `json:"reason"`
	RiskLevel     string `json:"risk_level"`
	RsikScore     int64  `json:"risk_score"`
	SellerMessage string `json:"seller_message"`
	Type          string `json:"type"`
}

type FraudDetails struct {
	UserReport string `json:"user_report"`
}

type Refunds struct {
	Object     string   `json:"object"`
	Data       []Refund `json:"data"`
	HasMore    bool     `json:"has_more"`
	TotalCount int64    `json:"total_count"`
	URL        string   `json:"url"`
}

type Refund struct {
	Amount                 int64             `json:"amount"`
	BalanceTransaction     string            `json:"balance_transaction"`
	Charge                 string            `json:"charge"`
	Created                int64             `json:"created"`
	Currency               string            `json:"currency"`
	ID                     string            `json:"id"`
	Metadata               map[string]string `json:"metadata"`
	Object                 string            `json:"object"`
	PaymentIntent          string            `json:"payment_intent"`
	Reason                 string            `json:"reason"`
	ReceiptNumber          string            `json:"receipt_number"`
	SourceTransferReversal interface{}       `json:"source_transfer_reversal"`
	Status                 string            `json:"status"`
	TransaferReversal      interface{}       `json:"transfer_reversal"`
}

type Customer struct {
	Address             Address           `json:"address"`
	Balance             int64             `json:"balance"`
	Created             int64             `json:"created"`
	Currency            string            `json:"currency"`
	DefaultSource       string            `json:"default_source"`
	Delinquent          bool              `json:"delinquent"`
	Description         string            `json:"description"`
	Discount            interface{}       `json:"discount"`
	Email               string            `json:"email"`
	ID                  string            `json:"id"`
	InvoicePrefix       string            `json:"invoice_prefix"`
	InvoiceSettings     InvoiceSettings   `json:"invoice_settings"`
	Livemode            bool              `json:"livemode"`
	Metadata            map[string]string `json:"metadata"`
	Name                string            `json:"name"`
	NextInvoiceSequence int64             `json:"next_invoice_sequence"`
	Object              string            `json:"object"`
	Phone               string            `json:"phone"`
	PreferredLocales    []string          `json:"preferred_locales"`
	Shipping            Address           `json:"shipping"`
	TaxExempt           string            `json:"tax_exempt"`
	TestClock           interface{}       `json:"test_clock"`
}

type InvoiceSettings struct {
	CustomFields         interface{} `json:"custom_fields"`
	DefaultPaymentMethod interface{} `json:"default_payment_method"`
	Footer               interface{} `json:"footer"`
}
