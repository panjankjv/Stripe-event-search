stripe-event-search
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22]

[1]: https://godoc.org/github.com/evalphobia/stripe-event-search?status.svg
[2]: https://godoc.org/github.com/evalphobia/stripe-event-search
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/stripe-event-search.svg
[6]: https://github.com/evalphobia/stripe-event-search/releases/latest
[7]: https://travis-ci.org/evalphobia/stripe-event-search.svg?branch=master
[8]: https://travis-ci.org/evalphobia/stripe-event-search
[9]: https://coveralls.io/repos/evalphobia/stripe-event-search/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/stripe-event-search?branch=master
[11]: https://codecov.io/github/evalphobia/stripe-event-search/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/stripe-event-search?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/stripe-event-search
[14]: https://goreportcard.com/report/github.com/evalphobia/stripe-event-search
[15]: https://img.shields.io/github/downloads/evalphobia/stripe-event-search/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/stripe-event-search/releases
[17]: https://img.shields.io/github/stars/evalphobia/stripe-event-search.svg
[18]: https://github.com/evalphobia/stripe-event-search/stargazers
[19]: https://codeclimate.com/github/evalphobia/stripe-event-search/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/stripe-event-search
[21]: https://bettercodehub.com/edge/badge/evalphobia/stripe-event-search?branch=master
[22]: https://bettercodehub.com/

`stripe-event-search` is a tool and golang library to search stripe events of the customers.

# Supported events

- Existing customer
- Existing payment methods
- PaymentIntent


# Installation

Install stripe-event-search by command below,

```bash
$ go install github.com/evalphobia/stripe-event-search/cmd@latest
```

# Usage

## root command

```bash
$ stripe-event-search
Commands:

  help     show help
  single   Exec searching stripe events for single customer
  multi    Exec searching stripe events for multiple customers
```

## single command

`single` command is to search single customer's payment data from Stripe.

```bash
# so far, restricted keys(rk_xxx) does not support search api.
$ export STRIPE_API_KEY=sk_test_...

$ stripe-event-search help single
Exec searching stripe events for single customer

Options:

  -h, --help                  display help information
      --apikey                api key for stripe API (e.g. --apikey='sk_test_xxx')
  -c, --customer              customer id for search query (e.g. --customer='cus_123')
  -p, --payment_type[=card]   target payment method type (e.g. --payment_type='card')
  -s, --show_metadata         metadata keys to show on output (space separated)  (e.g. --show_metadata='user_id user_name')
  -H, --hide                  ignore labels to hide from output (space separated)  (e.g. --hide='ID CardBrand')
  -k, --metakey               metadata key for search query (e.g. --metakey='user_id')
  -v, --metaval               metadata value for search query (e.g. --metaval='101')
  -A, --after                 filter payment events after this date/datetime (UTC) (e.g. --after='2022-01-31 10:00:00')
      --debug                 set if you need verbose logs --debug
```

For example, to get data of customer: `cus_123`

```bash
$ stripe-event-search single -c cus_123

Customer	CreatedTime	EventType	ID	CardBrand	CardLast4	CardFingerprint	Description	AmountCaptured	AmountRefunded	FailureCode	RiskScore
cus_123	2022-01-01T09:00:50+09:00	customer	cus_123				18424658
cus_123	2022-01-01T09:00:51+09:00	payment_intent	pi_001				1point	100	false	0		0
cus_123	2022-01-01T09:02:54+09:00	payment_intent	pi_002				10point	1000	false	0		0
cus_123	2022-01-01T09:04:37+09:00	payment_intent	pi_003	visa	0000	abcdEFG	5point	500	true	500		0
cus_123	2022-01-01T09:06:20+09:00	payment_method	pm_001	visa	0000	abcdEFG
```

If you want to search the customer by metadata, you can use `--metakey` (`-k`) and `--metaval` (`-v`) options.
Let's say the customer cus_123 has metadata `user_id` and the value is `101`, then,

```bash
$ stripe-event-search single -k user_id -v 101

search_meta_value	Customer	CreatedTime	EventType	ID	CardBrand	CardLast4	CardFingerprint	Description	AmountCaptured	AmountRefunded	FailureCode	RiskScore
101	cus_123	2022-01-01T09:00:50+09:00	customer	cus_123				18424658
101	cus_123	2022-01-01T09:00:51+09:00	payment_intent	pi_001				1point	100	false	0		0
101	cus_123	2022-01-01T09:02:54+09:00	payment_intent	pi_002				10point	1000	false	0		0
101	cus_123	2022-01-01T09:04:37+09:00	payment_intent	pi_003	visa	0000	abcdEFG	5point	500	true	500		0
101	cus_123	2022-01-01T09:06:20+09:00	payment_method	pm_001	visa	0000	abcdEFG
```

Also you can filter the time range and add/remove columns from output.

```bash
$stripe-event-search single \
  -k user_id \
  -v 101 \
  --after '2022-01-01 00:01:00' \
  -s 'user_id  country' \
  -H 'ID CardBrand CardLast4 Amount AmountRefunded RiskScore'

search_meta_value	Customer	CreatedTime	EventType	CardFingerprint	Description	Captured	FailureCode	user_id	country
101	cus_123	2022-01-01T09:00:50+09:00	customer		101			101
101	cus_123	2022-01-01T09:02:54+09:00	payment_intent		10point	false	101	jp
101	cus_123	2022-01-01T09:04:37+09:00	payment_intent	abcdEFG	5point	true		101	jp
101	cus_123	2022-01-01T09:06:20+09:00	payment_method	abcdEFG
```


## multi command

`multi` command is to search multiple customers' payment data from Stripe.

```bash
$ stripe-event-search help multi
Exec searching stripe events for multiple customers

Options:

  -h, --help                  display help information
      --apikey                api key for stripe API (e.g. --apikey='sk_test_xxx')
  -i, --input                *input csv/tsv file path (e.g. --input='./input.csv')
  -o, --output               *output tsv file path (e.g. --output='./output.tsv')
  -p, --payment_type[=card]   target payment method type (e.g. --payment_type='card')
  -s, --show_metadata         metadata keys to show on output (space separated)  (e.g. --show_metadata='user_id user_name')
  -H, --hide                  ignore labels to hide from output (space separated)  (e.g. --hide='ID CardBrand')
  -k, --metakey               metadata key for search query (e.g. --metakey='user_id')
  -A, --after                 filter payment events after this date/datetime (UTC) (e.g. --after='2022-01-31 10:00:00')
  -I, --interval              time interval after a API call to handle rate limit (ms=msec s=sec, m=min) (e.g. --interval=1.5s)
      --debug                 set if you need verbose logs
```

To use this command, specify customers on CSV/TSV file.
For example, create `input.tsv` file to get data of customer: `cus_123` and `cus_999` like below,

```bash
$ cat ./input.tsv
customer
cus_123
cus_999
```


Then use the file via `--input` option,

```bash
# so far, restricted keys(rk_xxx) does not support search api.
$ export STRIPE_API_KEY=sk_test_...

$ stripe-event-search multi -i ./input.tsv -o output.tsv
2022/03/22 13:35:56 [INFO] exec #: [1]
2022/03/22 13:35:57 [INFO] exec #: [2]
2022/03/22 13:35:58 [INFO] Finished


$ cat output.tsv
Customer	CreatedTime	EventType	ID	CardBrand	CardLast4	CardFingerprint	Description	Amount	Captured	AmountRefunded	FailureCode	RiskScore
cus_123	2022-01-01T09:00:50+09:00	customer	cus_123				18424658
cus_123	2022-01-01T09:00:51+09:00	payment_intent	pi_001				1point	100	false	0		0
cus_123	2022-01-01T09:02:54+09:00	payment_intent	pi_002				10point	1000	false	0		0
cus_123	2022-01-01T09:04:37+09:00	payment_intent	pi_003	visa	0000	abcdEFG	5point	500	true	500		0
cus_123	2022-01-01T09:06:20+09:00	payment_method	pm_001	visa	0000	abcdEFG
cus_999	2022-01-17T20:38:16+09:00	customer	cus_999				cus_999
```


If you want to search the customer by metadata, you can use `--metakey` (`-k`) option.
Let's say the customer cus_123 has metadata `user_id` and the value is `101`
and filter the time range and add/remove columns from output.

```bash
$ cat ./input.tsv
user_id
101
102


$ stripe-event-search multi \
  -i ./input.tsv \
  -o output.tsv \
  -k user_id  \
  --after '2022-01-01 00:01:00' \
  -s 'user_id  country' \
  -H 'ID CardBrand CardLast4 Amount AmountRefunded RiskScore'

2022/03/22 13:49:35 [INFO] exec #: [1]
2022/03/22 13:49:36 [INFO] exec #: [2]
2022/03/22 13:49:37 [INFO] Finished


$ cat ./output.tsv
search_meta_value	Customer	CreatedTime	EventType	CardFingerprint	Description	Captured	FailureCode	user_id	country
101	cus_123	2022-01-01T09:00:50+09:00	customer		101			101
101	cus_123	2022-01-01T09:02:54+09:00	payment_intent		10point	false		101	jp
101	cus_123	2022-01-01T09:04:37+09:00	payment_intent	abcdEFG	5point	true		101	jp
101	cus_123	2022-01-01T09:06:20+09:00	payment_method	abcdEFG
102	cus_999	2022-01-17T20:38:16+09:00	customer		102			102
```