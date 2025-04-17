# Postmark
> Fork from [Keighl's Postmark](https://github.com/keighl/postmark) (A Golang package for the using Postmark API)

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/postmark)](https://golang.org/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mrz1836/postmark/run-tests.yml?branch=master&logo=github&v=3)](https://github.com/mrz1836/postmark/actions)
[![codecov](https://codecov.io/gh/mrz1836/postmark/graph/badge.svg?token=j1KI99qHyW)](https://codecov.io/gh/mrz1836/postmark)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrz1836/postmark)](https://goreportcard.com/report/github.com/mrz1836/postmark)
[![Release](https://img.shields.io/github/release-pre/mrz1836/postmark.svg?style=flat&v=1)](https://github.com/mrz1836/postmark/releases)
[![GoDoc](https://godoc.org/github.com/mrz1836/postmark?status.svg)](https://pkg.go.dev/github.com/mrz1836/postmark)

<br/>

### Installation
```shell script
go get -u github.com/mrz1836/postmark
```

<br/>

### Basic Usage
Grab your [`Server Token`](https://account.postmarkapp.com/servers/XXXX/credentials), and your [`Account Token`](https://account.postmarkapp.com/account/edit).

```go
package main

import (
	"context"

	"github.com/mrz1836/postmark"
)

func main() {
	client := postmark.NewClient("[SERVER-TOKEN]", "[ACCOUNT-TOKEN]")

	email := postmark.Email{
		From:       "no-reply@example.com",
		To:         "tito@example.com",
		Subject:    "Reset your password",
		HTMLBody:   "...",
		TextBody:   "...",
		Tag:        "pw-reset",
		TrackOpens: true,
	}

	_, err := client.SendEmail(context.Background(), email)
	if err != nil {
		panic(err)
	}
}
```
Swap out HTTPClient for use on Google App Engine:

```go
package main

import (
    "github.com/mrz1836/postmark"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
)

// ....

client := postmark.NewClient("[SERVER-TOKEN]", "[ACCOUNT-TOKEN]")

ctx := appengine.NewContext(req)
client.HTTPClient = urlfetch.Client(ctx)

// ...
```

<br/>

### API Coverage

* [x] Emails
    * [x] `POST /email`
    * [x] `POST /email/batch`
    * [x] `POST /email/withTemplate`
    * [x] `POST /email/batchWithTemplates`
* [x] Bounces
    * [x] `GET /deliverystats`
    * [x] `GET /bounces`
    * [x] `GET /bounces/:id`
    * [x] `GET /bounces/:id/dump`
    * [x] `PUT /bounces/:id/activate`
    * [x] `GET /bounces/tags`
* [ ] Templates
    * [x] `GET /templates`
    * [x] `POST /templates`
    * [x] `GET /templates/:id`
    * [x] `PUT /templates/:id`
    * [x] `DELETE /templates/:id`
    * [x] `POST /templates/validate`
* [x] Suppressions
    * [x] `GET /suppressions/dump`
    * [x] `POST /suppressions`
    * [x] `POST /suppressions/delete`
* [x] Servers
    * [x] `GET /servers/:id`
    * [x] `PUT /servers/:id`
* [x] Outbound Messages
    * [x] `GET /messages/outbound`
    * [x] `GET /messages/outbound/:id/details`
    * [x] `GET /messages/outbound/:id/dump`
    * [x] `GET /messages/outbound/opens`
    * [x] `GET /messages/outbound/opens/:id`
* [x] Inbound Messages
    * [x] `GET /messages/inbound`
    * [x] `GET /messages/inbound/:id/details`
    * [x] `PUT /messages/inbound/:id/bypass`
    * [x] `PUT /messages/inbound/:id/retry`
* [x] Message Streams
    * [x] `GET /message-streams`
    * [x] `POST /message-streams`
    * [x] `GET /message-streams/{stream_ID}`
    * [x] `PATCH /message-streams/{stream_ID}`
    * [x] `POST /message-streams/{stream_ID}/archive`
    * [x] `POST /message-streams/{stream_ID}/unarchive`
* [ ] Sender signatures
    * [x] `GET /senders`
    * [ ] Get a sender signature’s details
    * [ ] Create a signature
    * [ ] Edit a signature
    * [ ] Delete a signature
    * [ ] Resend a confirmation
    * [ ] Verify an SPF record
    * [ ] Request a new DKIM
* [ ] Stats
    * [x] `GET /stats/outbound`
    * [x] `GET /stats/outbound/sends`
    * [x] `GET /stats/outbound/bounces`
    * [x] `GET /stats/outbound/spam`
    * [x] `GET /stats/outbound/tracked`
    * [x] `GET /stats/outbound/opens`
    * [x] `GET /stats/outbound/platform`
    * [ ] Get email client usage
    * [ ] Get email read times
* [ ] Triggers
    * [ ] Tags triggers
        * [ ] Create a trigger for a tag
        * [ ] Get a single trigger
        * [ ] Edit a single trigger
        * [ ] Delete a single trigger
        * [ ] Search triggers
    * [ ] Inbound rules triggers
        * [ ] Create a trigger for inbound rule
        * [ ] Delete a single trigger
        * [ ] List triggers
* [x] Webhooks
    * [x] List webhooks
    * [x] Get webhooks
    * [x] Create webhooks
    * [x] Edit webhooks
    * [x] Delete webhooks

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to GitHub and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                   Runs multiple commands
clean                 Remove previous builds and any test cache data
clean-mods            Remove all the Go mod cache
coverage              Shows the test coverage
diff                  Show the git diff
generate              Runs the go generate command in the base of the repo
godocs                Sync the latest tag with GoDocs
help                  Show this help message
install               Install the application
install-go            Install the application (Using Native Go)
install-releaser      Install the GoReleaser application
lint                  Run the golangci-lint application (install if not found)
release               Full production release (creates release in GitHub)
release               Runs common.release then runs godocs
release-snap          Test the full release (build binaries)
release-test          Full production test release (everything except deploy)
replace-version       Replaces the version in HTML/JS (pre-deploy)
run-examples          Runs all the examples
tag                   Generate a new tag and push (tag version=0.0.0)
tag-remove            Remove a tag if found (tag-remove version=0.0.0)
tag-update            Update an existing tag to current commit (tag-update version=0.0.0)
test                  Runs lint and ALL tests
test-ci               Runs all tests via CI (exports coverage)
test-ci-no-race       Runs all tests via CI (no race) (exports coverage)
test-ci-short         Runs unit tests via CI (exports coverage)
test-no-lint          Runs just tests
test-short            Runs vet, lint and tests (excludes integration tests)
test-unit             Runs tests and outputs coverage
uninstall             Uninstall the application (and remove files)
update-linter         Update the golangci-lint package (macOS only)
vet                   Run the Go vet application
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [GitHub Actions](https://github.com/mrz1836/postmark/actions) and
uses [Go version 1.18.x](https://golang.org/doc/go1.18). View the [configuration file](.github/workflows/run-tests.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

Run the examples:
```shell script
make run-examples
```

## License

[![License](https://img.shields.io/github/license/mrz1836/go-ses.svg?style=flat&v=2)](LICENSE)
