# 🤝 Contributing Guide

Thanks for taking the time to contribute! This project thrives on clear, well-tested, idiomatic Go code. Here's how you can help:

<br/>

## 📦 How to Contribute

1. Fork the repo.
2. Create a new branch.
3. Install the [pre-commit hooks](https://github.com/mrz1836/go-pre-commit).
4. Commit *one feature per commit*.
5. Write tests.
6. Open a pull request with a clear list of changes.

More info on [pull requests](http://help.github.com/pull-requests/).

<br/>

## 🧪 Testing

All tests follow standard Go patterns. We love:

* ✅ [Go Tests](https://golang.org/pkg/testing/)
* 📘 [Go Examples](https://golang.org/pkg/testing/#hdr-Examples)
* ⚡ [Go Benchmarks](https://golang.org/pkg/testing/#hdr-Benchmarks)

Tests should be:

* Easy to understand
* Focused on one behavior
* Fast

This project aims for >= **90% code coverage**. Every code path must be tested to
keep the Codecov badge green and CI passing.

<br/>

## 🧹 Coding Conventions

We follow [Effective Go](https://golang.org/doc/effective_go.html), plus:

* 📖 [godoc](https://godoc.org/golang.org/x/tools/cmd/godoc)
* 🧼 [golangci-lint](https://golangci-lint.run/)

Format your code with `gofmt`, lint with `golangci-lint`, and keep your diffs minimal.

<br/>

## 📚 More Guidance

For detailed workflows, commit standards, branch naming, PR templates, and more—read [AGENTS.md](./AGENTS.md). It’s the rulebook.

<br/>

Let’s build something great. 💪
