package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/evalphobia/stripe-event-search/lib"
	"github.com/evalphobia/stripe-event-search/stripesearch"
	"github.com/mkideal/cli"
)

// parameters of 'multi' command.
type multiT struct {
	cli.Helper
	APIKey            string `cli:"apikey" usage:"api key for stripe API (e.g. --apikey='sk_test_xxx')"`
	InputCSV          string `cli:"*i,input" usage:"input csv/tsv file path (e.g. --input='./input.csv')"`
	Output            string `cli:"*o,output" usage:"output tsv file path (e.g. --output='./output.tsv')"`
	PaymentType       string `cli:"p,payment_type" usage:"target payment method type (e.g. --payment_type='card')" dft:"card"`
	ShowMetadata      string `cli:"s,show_metadata" usage:"metadata keys to show on output (space separated)  (e.g. --show_metadata='user_id user_name')"`
	HideLabels        string `cli:"H,hide" usage:"ignore labels to hide from output (space separated)  (e.g. --hide='ID CardBrand')"`
	SearchMetadataKey string `cli:"k,metakey" usage:"metadata key for search query (e.g. --metakey='user_id')"`
	TimeAfter         string `cli:"A,after" usage:"filter payment events after this date/datetime (UTC) (e.g. --after='2022-01-31 10:00:00')"`
	Interval          string `cli:"I,interval" usage:"time interval after a API call to handle rate limit (ms=msec s=sec, m=min) (e.g. --interval=1.5s)"`
	Debug             bool   `cli:"debug" usage:"set if you need verbose logs"`
}

func (p *multiT) Validate(ctx *cli.Context) error {
	if p.Interval != "" {
		if _, err := time.ParseDuration(p.Interval); err != nil {
			return fmt.Errorf("invalid 'interval' format: [%w]", err)
		}
	}
	if p.TimeAfter != "" {
		if _, err := parseTime(p.TimeAfter); err != nil {
			return err
		}
	}

	if p.APIKey == "" && getAPIKey() == "" {
		return fmt.Errorf("you must set apikey via --apikey option or STRIPE_API_KEY environment variable")
	}

	return nil
}

var multiC = &cli.Command{
	Name: "multi",
	Desc: "Exec searching stripe events for multiple customers",
	Argv: func() interface{} { return new(multiT) },
	Fn:   execMulti,
}

func execMulti(ctx *cli.Context) error {
	argv := ctx.Argv().(*multiT)

	r := newMultiRunner(*argv)
	return r.Run()
}

type MultiRunner struct {
	runOption
	InputCSV string
	Output   string
	Interval time.Duration
}

func newMultiRunner(p multiT) MultiRunner {
	interval, _ := time.ParseDuration(p.Interval)
	dt, _ := parseTime(p.TimeAfter)

	return MultiRunner{
		InputCSV: p.InputCSV,
		Output:   p.Output,
		Interval: interval,
		runOption: runOption{
			PaymentType:       p.PaymentType,
			ShowMetadataKey:   parseStringSlice(p.ShowMetadata),
			HideLabels:        parseStringSlice(p.HideLabels),
			SearchMetadataKey: p.SearchMetadataKey,
			TimeAfter:         dt,
			Debug:             p.Debug,
		},
	}
}

func (r *MultiRunner) Run() error {
	opt := r.runOption
	f, err := lib.NewCSVHandler(r.InputCSV)
	if err != nil {
		return err
	}
	defer f.Close()

	// check csv file header
	needHeader := "customer"
	if opt.SearchMetadataKey != "" {
		needHeader = opt.SearchMetadataKey
	}
	if err := f.CheckHeaders(needHeader); err != nil {
		return err
	}

	w, err := lib.NewFileHandler(r.Output)
	if err != nil {
		return err
	}
	defer w.Close()

	searchKey := opt.SearchMetadataKey
	hasSearchKey := searchKey != ""
	metaKey := opt.ShowMetadataKey

	// write new header
	headers := stripesearch.GetHeaders(metaKey, opt.HideLabels, hasSearchKey)
	w.AppendRows([]string{strings.Join(headers, "\t")})

	cli := getClient(opt)
	var counter uint64
	for {
		line, err := f.Read()
		if err != nil {
			return err
		}
		if len(line) == 0 {
			break
		}
		counter++
		cli.LogInfo("exec #: [%d]\n", counter)

		results, err := fetchCustomersAllEvents(cli, runOption{
			Customer:            line["customer"],
			PaymentType:         opt.PaymentType,
			ShowMetadataKey:     opt.ShowMetadataKey,
			SearchMetadataKey:   searchKey,
			SearchMetadataValue: line[searchKey],
			TimeAfter:           opt.TimeAfter,
			Debug:               opt.Debug,
		})

		outputs := make([]string, 0, len(results)+1)
		switch {
		case err != nil:
			cli.LogError("error on #: [%d]; err=[%v]\n", counter, err)
			v := stripesearch.EventData{
				EventType: "error: " + err.Error(),
				Customer:  line["customer"],
			}
			results = append(results, v)
		case len(results) == 0:
			cli.LogError("cannot find events on #: [%d]\n", counter)
			v := stripesearch.EventData{
				EventType: "error: no event",
				Customer:  line["customer"],
			}
			results = append(results, v)
		}

		for _, v := range results {
			v.SetSearchMetaValue(line[searchKey])
			outputs = append(outputs, v.Output("\t", headers))
		}
		w.AppendRows(outputs)
		time.Sleep(r.Interval)
	}

	cli.LogInfo("Finished")
	return nil
}
