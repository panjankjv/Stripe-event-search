package main

import (
	"fmt"
	"strings"

	"github.com/mkideal/cli"

	"github.com/evalphobia/stripe-event-search/stripesearch"
)

// parameters of 'single' command.
type singleT struct {
	cli.Helper
	APIKey              string `cli:"apikey" usage:"api key for stripe API (e.g. --apikey='sk_test_xxx')"`
	Customer            string `cli:"c,customer" usage:"customer id for search query (e.g. --customer='cus_123')"`
	PaymentType         string `cli:"p,payment_type" usage:"target payment method type (e.g. --payment_type='card')" dft:"card"`
	ShowMetadata        string `cli:"s,show_metadata" usage:"metadata keys to show on output (space separated)  (e.g. --show_metadata='user_id user_name')"`
	HideLabels          string `cli:"H,hide" usage:"ignore labels to hide from output (space separated)  (e.g. --hide='ID CardBrand')"`
	SearchMetadataKey   string `cli:"k,metakey" usage:"metadata key for search query (e.g. --metakey='user_id')"`
	SearchMetadataValue string `cli:"v,metaval" usage:"metadata value for search query (e.g. --metaval='101')"`
	TimeAfter           string `cli:"A,after" usage:"filter payment events after this date/datetime (UTC) (e.g. --after='2022-01-31 10:00:00')"`
	Debug               bool   `cli:"debug" usage:"set if you need verbose logs --debug"`
}

func (p *singleT) Validate(ctx *cli.Context) error {
	if p.APIKey == "" && getAPIKey() == "" {
		return fmt.Errorf("you must set apikey via --apikey option or STRIPE_API_KEY environment variable")
	}
	if p.TimeAfter != "" {
		if _, err := parseTime(p.TimeAfter); err != nil {
			return err
		}
	}

	hasVal := false
	switch {
	case p.Customer != "",
		p.SearchMetadataKey != "" && p.SearchMetadataValue != "":
		hasVal = true
	case p.SearchMetadataKey != "",
		p.SearchMetadataValue != "":
		return fmt.Errorf("you must set both of --metakey and --metaval")
	}
	if !hasVal {
		return fmt.Errorf("at least, one of the query parameter is required (e.g. --customer='cus_123')")
	}

	return nil
}

var singleC = &cli.Command{
	Name: "single",
	Desc: "Exec searching stripe events for single customer",
	Argv: func() interface{} { return new(singleT) },
	Fn:   execSingle,
}

func execSingle(ctx *cli.Context) error {
	argv := ctx.Argv().(*singleT)

	r := newSingleRunner(*argv)
	return r.Run()
}

type SingleRunner struct {
	// parameters
	runOption
}

func newSingleRunner(p singleT) SingleRunner {
	dt, _ := parseTime(p.TimeAfter)

	return SingleRunner{
		runOption: runOption{
			Customer:            p.Customer,
			PaymentType:         p.PaymentType,
			ShowMetadataKey:     parseStringSlice(p.ShowMetadata),
			HideLabels:          parseStringSlice(p.HideLabels),
			SearchMetadataKey:   p.SearchMetadataKey,
			SearchMetadataValue: p.SearchMetadataValue,
			TimeAfter:           dt,
			Debug:               p.Debug,
		},
	}
}

func (r *SingleRunner) Run() error {
	opt := r.runOption
	cli := getClient(opt)
	results, err := fetchCustomersAllEvents(cli, opt)
	if err != nil {
		return err
	}

	hasSearchKey := opt.SearchMetadataKey != ""
	metaKey := r.ShowMetadataKey
	outputs := make([]string, 0, len(results)+1)
	headers := stripesearch.GetHeaders(metaKey, opt.HideLabels, hasSearchKey)
	outputs = append(outputs, strings.Join(headers, "\t"))
	for _, v := range results {
		v.SetSearchMetaValue(r.SearchMetadataValue)
		outputs = append(outputs, v.Output("\t", headers))
	}

	fmt.Printf("%s\n", strings.Join(outputs, "\n"))
	return nil
}
