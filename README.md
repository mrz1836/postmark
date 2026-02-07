<div align="center">

# 📨&nbsp;&nbsp;postmark

**Unofficial Golang Library for the Postmark API.**

<br/>

<a href="https://github.com/mrz1836/postmark/releases"><img src="https://img.shields.io/github/release-pre/mrz1836/postmark?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/mrz1836/postmark?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="https://github.com/mrz1836/postmark/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mrz1836/postmark?style=flat-square&color=blue" alt="License"></a>

<br/>

<table align="center" border="0">
  <tr>
    <td align="right">
       <code>CI / CD</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/postmark/actions"><img src="https://img.shields.io/github/actions/workflow/status/mrz1836/postmark/fortress.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build"></a>
       <a href="https://github.com/mrz1836/postmark/actions"><img src="https://img.shields.io/github/last-commit/mrz1836/postmark?style=flat-square&logo=git&logoColor=white&label=last%20update" alt="Last Commit"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Quality</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://goreportcard.com/report/github.com/mrz1836/postmark"><img src="https://goreportcard.com/badge/github.com/mrz1836/postmark?style=flat-square" alt="Go Report"></a>
       <a href="https://codecov.io/gh/mrz1836/postmark"><img src="https://codecov.io/gh/mrz1836/postmark/branch/master/graph/badge.svg?style=flat-square" alt="Coverage"></a>
    </td>
  </tr>

  <tr>
    <td align="right">
       <code>Security</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://scorecard.dev/viewer/?uri=github.com/mrz1836/postmark"><img src="https://api.scorecard.dev/projects/github.com/mrz1836/postmark/badge?style=flat-square" alt="Scorecard"></a>
       <a href=".github/SECURITY.md"><img src="https://img.shields.io/badge/policy-active-success?style=flat-square&logo=security&logoColor=white" alt="Security"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Community</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/postmark/graphs/contributors"><img src="https://img.shields.io/github/contributors/mrz1836/postmark?style=flat-square&color=orange" alt="Contributors"></a>
       <a href="https://mrz1818.com/"><img src="https://img.shields.io/badge/donate-bitcoin-ff9900?style=flat-square&logo=bitcoin" alt="Bitcoin"></a>
    </td>
  </tr>
</table>

</div>

<br/>
<br/>

<div align="center">

### <code>Project Navigation</code>

</div>

<table align="center">
  <tr>
    <td align="center" width="33%">
       🚀&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="33%">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
    <td align="center" width="33%">
       📚&nbsp;<a href="#-documentation"><code>Documentation</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
      🛠️&nbsp;<a href="#-code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
    <td align="center">
      ⚡&nbsp;<a href="#-benchmarks"><code>Benchmarks</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
      🤖&nbsp;<a href="#-ai-usage--assistant-guidelines"><code>AI&nbsp;Usage</code></a>
    </td>
    <td align="center">
       ⚖️&nbsp;<a href="#-license"><code>License</code></a>
    </td>
    <td align="center">
       👥&nbsp;<a href="#-maintainers"><code>Maintainers</code></a>
    </td>
  </tr>
</table>
<br/>

## 📦 Installation

**postmark** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get github.com/mrz1836/postmark
```

<br/>

## 💡 Usage
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
<br/>

## 📚 Documentation

View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/postmark?tab=doc)

> **Heads up!** `postmark` is intentionally light on dependencies. The only
external package it uses is the excellent `testify` suite—and that's just for
our tests. You can drop this library into your projects without dragging along
extra baggage.

<br/>

<details>
<summary><strong><code>Supported API Coverage</code></strong></summary>
<br/>

* [x] **[Email API](https://postmarkapp.com/developer/api/email-api) - ([email.go](email.go))**
	* [x] [`POST /email`](https://postmarkapp.com/developer/api/email-api#send-a-single-email) - Send a single email
	* [x] [`POST /email/batch`](https://postmarkapp.com/developer/api/email-api#send-batch-emails) - Send batch emails

* [x] **[Templates API](https://postmarkapp.com/developer/api/templates-api) - ([templates.go](templates.go))**
	* [x] [`POST /email/withTemplate`](https://postmarkapp.com/developer/api/templates-api#email-with-template) - Send email with template
	* [x] [`POST /email/batchWithTemplates`](https://postmarkapp.com/developer/api/templates-api#send-batch-with-templates) - Send batch with templates
	* [x] [`PUT /templates/push`](https://postmarkapp.com/developer/api/templates-api#push-templates) - Push templates to another server
	* [x] [`GET /templates/{templateIdOrAlias}`](https://postmarkapp.com/developer/api/templates-api#get-template) - Get a template
	* [x] [`POST /templates`](https://postmarkapp.com/developer/api/templates-api#create-template) - Create a template
	* [x] [`PUT /templates/{templateIdOrAlias}`](https://postmarkapp.com/developer/api/templates-api#edit-template) - Edit a template
	* [x] [`GET /templates`](https://postmarkapp.com/developer/api/templates-api#list-templates) - List templates
	* [x] [`DELETE /templates/{templateIdOrAlias}`](https://postmarkapp.com/developer/api/templates-api#delete-template) - Delete a template
	* [x] [`POST /templates/validate`](https://postmarkapp.com/developer/api/templates-api#validate-template) - Validate a template

* [x] **[Bounce API](https://postmarkapp.com/developer/api/bounce-api) - ([bounce.go](bounce.go))**
	* [x] [`GET /deliverystats`](https://postmarkapp.com/developer/api/bounce-api#get-delivery-stats) - Get delivery stats
	* [x] [`GET /bounces`](https://postmarkapp.com/developer/api/bounce-api#get-bounces) - Get bounces
	* [x] [`GET /bounces/{bounceid}`](https://postmarkapp.com/developer/api/bounce-api#get-bounce) - Get a single bounce
	* [x] [`GET /bounces/{bounceid}/dump`](https://postmarkapp.com/developer/api/bounce-api#get-bounce-dump) - Get bounce dump
	* [x] [`PUT /bounces/{bounceid}/activate`](https://postmarkapp.com/developer/api/bounce-api#activate-bounce) - Activate a bounce
	* [x] [`GET /bounces/tags`](https://postmarkapp.com/developer/api/bounce-api#get-bounced-tags) - Get bounced tags

* [x] **[Messages API](https://postmarkapp.com/developer/api/messages-api) - ([messages_outbound.go](messages_outbound.go), [messages_inbound.go](messages_inbound.go))**
	* [x] [`GET /messages/outbound`](https://postmarkapp.com/developer/api/messages-api#outbound-message-search) - Search outbound messages
	* [x] [`GET /messages/outbound/{messageid}/details`](https://postmarkapp.com/developer/api/messages-api#outbound-message-details) - Get outbound message details
	* [x] [`GET /messages/outbound/{messageid}/dump`](https://postmarkapp.com/developer/api/messages-api#outbound-message-dump) - Get outbound message dump
	* [x] [`GET /messages/outbound/opens`](https://postmarkapp.com/developer/api/messages-api#message-opens) - Get message opens
	* [x] [`GET /messages/outbound/opens/{messageid}`](https://postmarkapp.com/developer/api/messages-api#opens-for-single-message) - Get opens for single message
	* [x] [`GET /messages/outbound/clicks`](https://postmarkapp.com/developer/api/messages-api#message-clicks) - Get message clicks
	* [x] [`GET /messages/outbound/clicks/{messageid}`](https://postmarkapp.com/developer/api/messages-api#clicks-for-single-message) - Get clicks for single message
	* [x] [`GET /messages/inbound`](https://postmarkapp.com/developer/api/messages-api#inbound-message-search) - Search inbound messages
	* [x] [`GET /messages/inbound/{messageid}/details`](https://postmarkapp.com/developer/api/messages-api#inbound-message-details) - Get inbound message details
	* [x] [`PUT /messages/inbound/{messageid}/bypass`](https://postmarkapp.com/developer/api/messages-api#bypass-inbound-message-rules) - Bypass inbound message rules
	* [x] [`PUT /messages/inbound/{messageid}/retry`](https://postmarkapp.com/developer/api/messages-api#retry-inbound-message-processing) - Retry inbound message processing

* [x] **[Message Streams API](https://postmarkapp.com/developer/api/message-streams-api) - ([message_streams.go](message_streams.go))**
	* [x] [`GET /message-streams`](https://postmarkapp.com/developer/api/message-streams-api#list-message-streams) - List message streams
	* [x] [`GET /message-streams/{stream_ID}`](https://postmarkapp.com/developer/api/message-streams-api#get-message-stream) - Get a message stream
	* [x] [`PATCH /message-streams/{stream_ID}`](https://postmarkapp.com/developer/api/message-streams-api#edit-message-stream) - Edit a message stream
	* [x] [`POST /message-streams`](https://postmarkapp.com/developer/api/message-streams-api#create-message-stream) - Create a message stream
	* [x] [`POST /message-streams/{stream_ID}/archive`](https://postmarkapp.com/developer/api/message-streams-api#archive-message-stream) - Archive a message stream
	* [x] [`POST /message-streams/{stream_ID}/unarchive`](https://postmarkapp.com/developer/api/message-streams-api#unarchive-message-stream) - Unarchive a message stream

* [x] **[Domains API](https://postmarkapp.com/developer/api/domains-api) - ([domains.go](domains.go))**
	* [x] [`GET /domains`](https://postmarkapp.com/developer/api/domains-api#list-domains) - List domains
	* [x] [`GET /domains/{domainid}`](https://postmarkapp.com/developer/api/domains-api#get-domain-details) - Get domain details
	* [x] [`POST /domains`](https://postmarkapp.com/developer/api/domains-api#create-domain) - Create a domain
	* [x] [`PUT /domains/{domainid}`](https://postmarkapp.com/developer/api/domains-api#edit-domain) - Edit a domain
	* [x] [`DELETE /domains/{domainid}`](https://postmarkapp.com/developer/api/domains-api#delete-domain) - Delete a domain
	* [x] [`PUT /domains/{domainid}/verifyDkim`](https://postmarkapp.com/developer/api/domains-api#verify-dkim) - Verify DKIM status
	* [x] [`PUT /domains/{domainid}/verifyReturnPath`](https://postmarkapp.com/developer/api/domains-api#verify-return-path) - Verify return-path status
	* [x] [`POST /domains/{domainid}/rotatedkim`](https://postmarkapp.com/developer/api/domains-api#rotate-dkim) - Rotate DKIM keys

* [x] **[Sender Signatures API](https://postmarkapp.com/developer/api/signatures-api) - ([sender_signatures.go](sender_signatures.go))**
	* [x] [`GET /senders`](https://postmarkapp.com/developer/api/signatures-api#list-sender-signatures) - List sender signatures
	* [x] [`GET /senders/{signatureid}`](https://postmarkapp.com/developer/api/signatures-api#get-sender-signature-details) - Get sender signature details
	* [x] [`POST /senders`](https://postmarkapp.com/developer/api/signatures-api#create-signature) - Create a signature
	* [x] [`PUT /senders/{signatureid}`](https://postmarkapp.com/developer/api/signatures-api#edit-signature) - Edit a signature
	* [x] [`DELETE /senders/{signatureid}`](https://postmarkapp.com/developer/api/signatures-api#delete-signature) - Delete a signature
	* [x] [`POST /senders/{signatureid}/resend`](https://postmarkapp.com/developer/api/signatures-api#resend-confirmation) - Resend confirmation

* [x] **[Stats API](https://postmarkapp.com/developer/api/stats-api) - ([stats.go](stats.go))**
	* [x] [`GET /stats/outbound`](https://postmarkapp.com/developer/api/stats-api#get-outbound-overview) - Get outbound overview
	* [x] [`GET /stats/outbound/sends`](https://postmarkapp.com/developer/api/stats-api#get-sent-counts) - Get sent counts
	* [x] [`GET /stats/outbound/bounces`](https://postmarkapp.com/developer/api/stats-api#get-bounce-counts) - Get bounce counts
	* [x] [`GET /stats/outbound/spam`](https://postmarkapp.com/developer/api/stats-api#get-spam-complaints) - Get spam complaints
	* [x] [`GET /stats/outbound/tracked`](https://postmarkapp.com/developer/api/stats-api#get-tracked-email-counts) - Get tracked email counts
	* [x] [`GET /stats/outbound/opens`](https://postmarkapp.com/developer/api/stats-api#get-email-open-counts) - Get email open counts
	* [x] [`GET /stats/outbound/opens/platforms`](https://postmarkapp.com/developer/api/stats-api#get-email-platform-usage) - Get email platform usage
	* [x] [`GET /stats/outbound/opens/emailclients`](https://postmarkapp.com/developer/api/stats-api#get-email-client-usage) - Get email client usage
	* [x] [`GET /stats/outbound/clicks`](https://postmarkapp.com/developer/api/stats-api#get-click-counts) - Get click counts
	* [x] [`GET /stats/outbound/clicks/browserfamilies`](https://postmarkapp.com/developer/api/stats-api#get-browser-usage) - Get browser usage
	* [x] [`GET /stats/outbound/clicks/platforms`](https://postmarkapp.com/developer/api/stats-api#get-browser-platform-usage) - Get browser platform usage
	* [x] [`GET /stats/outbound/clicks/location`](https://postmarkapp.com/developer/api/stats-api#get-click-location) - Get click location

* [x] **[Webhooks API](https://postmarkapp.com/developer/api/webhooks-api) - ([webhooks.go](webhooks.go))**
	* [x] [`GET /webhooks`](https://postmarkapp.com/developer/api/webhooks-api#list-webhooks) - List webhooks
	* [x] [`GET /webhooks/{Id}`](https://postmarkapp.com/developer/api/webhooks-api#get-webhook) - Get a webhook
	* [x] [`POST /webhooks`](https://postmarkapp.com/developer/api/webhooks-api#create-webhook) - Create a webhook
	* [x] [`PUT /webhooks/{Id}`](https://postmarkapp.com/developer/api/webhooks-api#edit-webhook) - Edit a webhook
	* [x] [`DELETE /webhooks/{Id}`](https://postmarkapp.com/developer/api/webhooks-api#delete-webhook) - Delete a webhook

* [x] **[Suppressions API](https://postmarkapp.com/developer/api/suppressions-api) - ([suppressions.go](suppressions.go))**
	* [x] [`GET /message-streams/{stream_id}/suppressions/dump`](https://postmarkapp.com/developer/api/suppressions-api#suppression-dump) - Suppression dump
	* [x] [`POST /message-streams/{stream_id}/suppressions`](https://postmarkapp.com/developer/api/suppressions-api#create-suppression) - Create suppressions
	* [x] [`POST /message-streams/{stream_id}/suppressions/delete`](https://postmarkapp.com/developer/api/suppressions-api#delete-suppression) - Delete suppressions

* [x] **[Servers API](https://postmarkapp.com/developer/api/servers-api) - ([server.go](server.go), [servers.go](servers.go))**
	* [x] [`GET /server`](https://postmarkapp.com/developer/api/servers-api#get-server) - Get current server
	* [x] [`PUT /server`](https://postmarkapp.com/developer/api/servers-api#edit-server) - Edit current server
	* [x] [`GET /servers/{serverid}`](https://postmarkapp.com/developer/api/servers-api#get-server) - Get a server
	* [x] [`POST /servers`](https://postmarkapp.com/developer/api/servers-api#create-server) - Create a server
	* [x] [`PUT /servers/{serverid}`](https://postmarkapp.com/developer/api/servers-api#edit-server) - Edit a server
	* [x] [`GET /servers`](https://postmarkapp.com/developer/api/servers-api#list-servers) - List servers
	* [x] [`DELETE /servers/{serverid}`](https://postmarkapp.com/developer/api/servers-api#delete-server) - Delete a server

* [x] **[Inbound Rules Triggers API](https://postmarkapp.com/developer/api/inbound-rules-triggers-api) - ([inbound_rules_triggers.go](inbound_rules_triggers.go))**
	* [x] [`GET /triggers/inboundrules`](https://postmarkapp.com/developer/api/inbound-rules-triggers-api#list-inbound-rule-triggers) - List inbound rule triggers
	* [x] [`POST /triggers/inboundrules`](https://postmarkapp.com/developer/api/inbound-rules-triggers-api#create-inbound-rule-trigger) - Create an inbound rule trigger
	* [x] [`DELETE /triggers/inboundrules/{triggerid}`](https://postmarkapp.com/developer/api/inbound-rules-triggers-api#delete-trigger) - Delete a single trigger

* [x] **[Data Removal API](https://postmarkapp.com/developer/api/data-removals-api) - ([data_removals.go](data_removals.go))**
	* [x] [`POST /data-removals`](https://postmarkapp.com/developer/api/data-removals-api#create-data-removal-request) - Create a data removal request
	* [x] [`GET /data-removals/{id}`](https://postmarkapp.com/developer/api/data-removals-api#check-data-removal-status) - Check a data removal request status

</details>

<details>
<summary><strong><code>Custom HTTPClient Support</code></strong></summary>
<br/>

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
</details>

<details>
<summary><strong><code>Development Setup (Getting Started)</code></strong></summary>
<br/>

Install [MAGE-X](https://github.com/mrz1836/mage-x) build tool for development:

```bash
# Install MAGE-X for development and building
go install github.com/mrz1836/mage-x/cmd/magex@latest
magex update:install
```
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.

Then create and push a new Git tag using:

```bash
magex version:bump bump=patch push=true branch=master
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Build Commands</code></strong></summary>
<br/>

View all build commands

```bash script
magex help
```

</details>

<details>
<summary><strong>GitHub Workflows</strong></summary>
<br/>

All workflows are driven by modular configuration in [`.github/env/`](.github/env/README.md) — no YAML editing required.

**[View all workflows and the control center →](.github/docs/workflows.md)**

</details>

<details>
<summary><strong><code>Updating Dependencies</code></strong></summary>
<br/>

To update all dependencies (Go modules, linters, and related tools), run:

```bash
magex deps:update
```

This command ensures all dependencies are brought up to date in a single step, including Go modules and any managed tools. It is the recommended way to keep your development environment and CI in sync with the latest versions.

</details>

<br/>

## Examples & Tests
## 🧪 Examples & Tests

All unit tests and fuzz tests run via [GitHub Actions](https://github.com/mrz1836/postmark/actions) and use [Go version 1.18.x](https://go.dev/doc/go1.18). View the [configuration file](.github/workflows/fortress.yml).

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
```

<br/>

## ⚡ Benchmarks

Run the Go benchmarks:

```bash script
magex bench
```

### 📊 Performance Results

All benchmarks measure **real API client performance** including HTTP request setup, JSON marshaling/unmarshalling, and response processing against mock servers. Results collected on Apple M1 Max (10 cores).

#### 🎯 Performance Overview

| Metric                | Value           | Description              |
|-----------------------|-----------------|--------------------------|
| **Fastest Operation** | 36.7 µs         | Get Bounced Tags         |
| **Average Latency**   | 41.2 µs         | Across all 47 operations |
| **Throughput**        | ~24,000 ops/sec | Per operation average    |
| **Memory Efficiency** | 7.7 KB/op       | Average memory usage     |
| **Allocations**       | 97 allocs/op    | Average per operation    |

<details>
<summary><strong>Bounce API Performance</strong></summary>
<br/>

| Operation          | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|--------------------|--------------|----------------------|--------|--------|
| Get Delivery Stats | 38.0         | 26,300               | 6.8 KB | 86     |
| Get Bounces        | 41.6         | 24,000               | 7.8 KB | 110    |
| Get Bounce         | 40.4         | 24,800               | 7.1 KB | 89     |
| Get Bounce Dump    | 37.4         | 26,700               | 6.6 KB | 84     |
| Activate Bounce    | 39.9         | 25,100               | 7.2 KB | 92     |
| Get Bounced Tags   | 36.7         | 27,200               | 6.6 KB | 85     |

</details>

<details>
<summary><strong>Data Removal API Performance</strong></summary>
<br/>

| Operation               | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|-------------------------|--------------|----------------------|--------|--------|
| Create Data Removal     | 41.5         | 24,100               | 7.6 KB | 100    |
| Get Data Removal Status | 38.6         | 25,900               | 6.8 KB | 86     |

</details>

<details>
<summary><strong>Domains API Performance</strong></summary>
<br/>

| Operation          | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|--------------------|--------------|----------------------|--------|--------|
| Get Domains        | 40.2         | 24,900               | 7.2 KB | 101    |
| Get Domain         | 41.9         | 23,900               | 7.3 KB | 89     |
| Create Domain      | 41.3         | 24,200               | 7.8 KB | 100    |
| Edit Domain        | 41.7         | 24,000               | 8.3 KB | 107    |
| Delete Domain      | 38.2         | 26,200               | 7.1 KB | 89     |
| Verify DKIM Status | 39.6         | 25,200               | 7.4 KB | 91     |
| Verify Return Path | 39.2         | 25,500               | 7.4 KB | 90     |
| Rotate DKIM        | 40.2         | 24,900               | 7.6 KB | 93     |

</details>

<details>
<summary><strong>Inbound Rules Triggers API Performance</strong></summary>
<br/>

| Operation                   | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|-----------------------------|--------------|----------------------|--------|--------|
| Get Inbound Rule Triggers   | 39.8         | 25,100               | 7.1 KB | 101    |
| Create Inbound Rule Trigger | 41.0         | 24,400               | 7.6 KB | 99     |
| Delete Inbound Rule Trigger | 40.4         | 24,700               | 6.7 KB | 84     |

</details>

<details>
<summary><strong>Message Streams API Performance</strong></summary>
<br/>

| Operation                | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|--------------------------|--------------|----------------------|--------|--------|
| List Message Streams     | 44.4         | 22,500               | 7.4 KB | 93     |
| Get Message Stream       | 42.6         | 23,500               | 7.0 KB | 89     |
| Edit Message Stream      | 46.8         | 21,400               | 8.1 KB | 106    |
| Create Message Stream    | 44.5         | 22,500               | 8.1 KB | 104    |
| Archive Message Stream   | 40.4         | 24,800               | 6.8 KB | 86     |
| Unarchive Message Stream | 42.6         | 23,500               | 7.1 KB | 90     |

</details>

<details>
<summary><strong>Messages API Performance</strong></summary>
<br/>

| Operation                    | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|------------------------------|--------------|----------------------|--------|--------|
| Get Outbound Messages Clicks | 47.5         | 21,100               | 8.5 KB | 118    |
| Get Outbound Message Clicks  | 43.2         | 23,100               | 7.9 KB | 109    |

</details>

<details>
<summary><strong>Sender Signatures API Performance</strong></summary>
<br/>

| Operation                     | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|-------------------------------|--------------|----------------------|--------|--------|
| Get Sender Signatures         | 40.6         | 24,600               | 7.3 KB | 104    |
| Get Sender Signature          | 40.6         | 24,600               | 7.5 KB | 92     |
| Create Sender Signature       | 42.2         | 23,700               | 8.1 KB | 101    |
| Edit Sender Signature         | 47.1         | 21,200               | 8.6 KB | 108    |
| Delete Sender Signature       | 38.8         | 25,800               | 7.1 KB | 89     |
| Resend Signature Confirmation | 39.0         | 25,700               | 7.2 KB | 90     |

</details>

<details>
<summary><strong>Stats API Performance</strong></summary>
<br/>

| Operation                 | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|---------------------------|--------------|----------------------|--------|--------|
| Get Click Counts          | 40.3         | 24,800               | 7.4 KB | 103    |
| Get Browser Family Counts | 42.2         | 23,700               | 7.6 KB | 103    |
| Get Click Location Counts | 42.8         | 23,400               | 7.4 KB | 103    |
| Get Click Platform Counts | 42.0         | 23,800               | 7.5 KB | 103    |
| Get Email Client Counts   | 41.3         | 24,200               | 7.6 KB | 103    |

</details>

<details>
<summary><strong>Templates API Performance</strong></summary>
<br/>

| Operation                  | Latency (µs) | Throughput (ops/sec) | Memory | Allocs |
|----------------------------|--------------|----------------------|--------|--------|
| Get Template               | 39.9         | 25,100               | 7.5 KB | 92     |
| Get Templates              | 41.3         | 24,200               | 7.5 KB | 103    |
| Get Templates Filtered     | 40.0         | 25,000               | 7.4 KB | 103    |
| Create Template            | 44.7         | 22,400               | 7.9 KB | 99     |
| Edit Template              | 42.5         | 23,500               | 8.4 KB | 106    |
| Delete Template            | 39.3         | 25,500               | 7.1 KB | 89     |
| Validate Template          | 44.9         | 22,300               | 8.5 KB | 110    |
| Send Templated Email       | 44.2         | 22,600               | 8.8 KB | 110    |
| Send Templated Email Batch | 46.1         | 21,700               | 9.0 KB | 117    |
| Push Templates             | 42.9         | 23,300               | 7.9 KB | 105    |

</details>

> **Note:** All benchmarks use mock HTTP servers for consistent, reproducible measurements. Real-world performance will vary based on network latency and Postmark API response times.

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 🤖 AI Usage & Assistant Guidelines
Read the [AI Usage & Assistant Guidelines](.github/tech-conventions/ai-compliance.md) for details on how AI is used in this project and how to interact with the AI assistants.

<br/>

## 👥 Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## 🤝 Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap:
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=postmark&utm_term=postmark&utm_content=postmark) to ensure this journey continues indefinitely! :rocket:


[![Stars](https://img.shields.io/github/stars/mrz1836/postmark?label=Please%20like%20us&style=social)](https://github.com/mrz1836/postmark/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/github/license/mrz1836/postmark.svg?style=flat)](LICENSE)
