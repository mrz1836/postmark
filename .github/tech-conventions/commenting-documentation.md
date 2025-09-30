# Commenting & Documentation Standards

> Great engineers write great comments. You're not here to state the obvious—you're here to document decisions, highlight edge cases, and make sure the next dev (or AI) doesn't repeat your mistakes.

<br><br>

## 🧠 Guiding Principles

* **Comment the "why", not the "what"**

  > The code already tells us *what* it's doing. Your job is to explain *why* it's doing it that way—especially if it's non-obvious, nuanced, or a workaround.

* **Explain side effects, caveats, and constraints**

  > If the function touches global state, writes to disk, mutates shared memory, or makes assumptions—write it down.

* **Don't comment on broken code—fix or delete it**

  > Dead or disabled code with TODOs are bad signals. If it's not worth fixing now, delete it and add an issue instead.

* **Your comments are part of the product**

  > Treat them like UX copy. Make them clear, concise, and professional. You're writing for peers, not compilers.

<br><br>

## 🔤 Function Comments (Exported)

Every exported function **must** include a Go-style comment that:

* Starts with the function name
* States its purpose clearly
* Documents:
  * **Steps**: Include if the function performs a non-obvious sequence of operations.
  * **Parameters**: Always describe all parameters when present.
  * **Return values**: Document return types and their meaning if not trivially understood.
  * **Side effects**: Note any I/O, database writes, external calls, or mutations that aren't local to the function.
  * **Notes**: Include any assumptions, constraints, or important context that the caller should know.

Here is a template for function comments that is recommended to use:

```go
// FunctionName does [what this function does] in [brief context].
//
// This function performs the following steps: [if applicable, describe the main steps in a bullet list]
// - [First major action or check performed]
// - [Second action or branching logic explained, if relevant]
//    - [Details about possible outcomes or internal branching]
// - [Additional steps with sub-bullets as needed]
// - [Final steps and cleanup, if applicable]
//
// Parameters:
// - ctx: [Purpose of context in this function]
// - paramName: [Explanation of each parameter and what it controls or affects]
//
// Returns:
// - [What is returned; error behavior if any]
//
// Side Effects:
// - [Any side effects, such as modifying global state, writing to disk, etc.]
//
// Notes:
// - [Caveats, assumptions, or recommendations—e.g., transaction usage, concurrency, etc.]
// - [Any implicit contracts with the caller or system constraints]
// - [Mention if this function is part of a larger workflow or job system]
```

<br><br>

## 📦 Package-Level Comments

* Each package **must** include a package-level comment in a file named after the package (e.g., `auth.go` for package `auth`).
* If no logical file fits, add a `doc.go` with the comment block.
* Use it to explain:
    * The package purpose
    * High-level API boundaries
    * Expected use-cases and design notes

Here is a template for package comments that is recommended to use:

```go
// Package PackageName provides [brief description of the package's purpose].
//
// This package implements [high-level functionality or API] and is designed to [describe intended use cases].
//
// Key features include: [if applicable, list key features or components]
// - [Feature 1: brief description of what it does]
// - [Feature 2: brief description of what it does]
// - [Feature 3: brief description of what it does]
//
// The package is structured to [explain any architectural decisions, e.g., modularity, separation of concerns].
// It relies on [mention any key dependencies or external systems].
//
// Usage examples:
// [Provide a simple example of how to use the package, if applicable]
//
// Important notes:
// - [Any important assumptions or constraints, e.g., concurrency model, state management]
// - [Any known limitations or edge cases]
// - [Any specific configuration or initialization steps required]
//
// This package is part of the larger [Project Name] ecosystem and interacts with [mention related packages or systems].
package PackageName
```

<br><br>

## 🧱 Inline Comments

Use inline comments **strategically**, not excessively.

* Use them to explain "weird" logic, workarounds, or business rules.
* Prefer **block comments (`//`)** on their own line over trailing comments.
* Avoid obvious noise:

🚫 `i++ // increment i`

✅ `// Skip empty rows to avoid panic on CSV parse`

<br><br>

## ⚙️ Comment Style

* Use **complete sentences** with punctuation.
* Keep your tone **precise, confident, and modern**—you're not writing a novel, but you're also not writing legacy COBOL.
* Avoid filler like "simple function" or "just does X".
* Don't leave TODOs unless:
    * They are immediately actionable
    * (or) they reference an issue
    * They include a timestamp or owner

<br><br>

## 🧬 AI Agent Directives

If you're an AI contributing code:

* Treat your comments like commit messages—**use active voice, be declarative**
* Use comments to **make intent explicit**, especially for generated or AI-authored sections
* Avoid hallucinating context—if you're unsure, omit or tag with `// AI: review this logic`
* Flag areas of uncertainty or external dependency (e.g., "// AI: relies on external config structure")

<br><br>

## 🔥 Comment Hygiene

* Remove outdated comments aggressively.
* Keep comments synced with refactoring.
* Use `//nolint:<linter> // message` only with clear, justified context and explanation.

<br><br>

## 📝 Markdown Documents

Markdown files (e.g., `README.md`, `AGENTS.md`, `CONTRIBUTING.md`) are first-class citizens in this repository. Edits must follow these best practices:

* **Write with intent** — Be concise, clear, and audience-aware. Each section should serve a purpose.
* **Use proper structure** — Maintain consistent heading levels, spacing, and bullet/number list formatting.
* **Full Table Borders** — Use full borders for tables to ensure readability across different renderers.
* **Table Border Spacing** — Make sure tables have appropriate spacing for clarity.
* **Preserve voice and tone** — Match the professional tone and style used across the project documentation.
* **Preview before committing** — Always verify rendered output locally or in a PR to avoid broken formatting.
* **Update references** — When renaming files or sections, update internal links and the table of contents if present.

> Markdown updates should be treated with the same care as code—clean, purposeful, and reviewed.

<br><br>

## 📖 Documentation Examples

### Good Function Documentation

```go
// ProcessBatchWithRetry processes a batch of items with automatic retry logic
// for transient failures.
//
// This function performs the following steps:
// - Validates the input batch for basic sanity checks
// - Chunks the batch into smaller groups based on maxBatchSize
// - Processes each chunk with exponential backoff retry
//   - On success: moves to next chunk
//   - On failure: retries up to maxRetries times
// - Aggregates results and errors from all chunks
//
// Parameters:
// - ctx: Controls cancellation and timeout for the entire operation
// - items: The batch of items to process (must not be nil)
// - processor: The function that processes individual items
// - opts: Configuration options for retry behavior and batch sizing
//
// Returns:
// - []Result: Successfully processed results (may be partial on error)
// - error: First critical error encountered, or nil if all items processed
//
// Side Effects:
// - May write to configured logger on retries and failures
// - Consumes rate limit tokens from the global rate limiter
//
// Notes:
// - This function is safe for concurrent use
// - Partial results are returned even if some items fail
// - The processor function must be idempotent for retry safety
func ProcessBatchWithRetry(ctx context.Context, items []Item, processor ItemProcessor, opts *RetryOptions) ([]Result, error) {
    // Implementation...
}
```

### Good Package Documentation

```go
// Package ratelimit provides thread-safe rate limiting primitives for controlling
// API request rates and preventing resource exhaustion.
//
// This package implements token bucket and sliding window algorithms and is designed
// to work seamlessly with HTTP middleware and gRPC interceptors.
//
// Key features include:
// - Token bucket rate limiter with configurable capacity and refill rate
// - Sliding window rate limiter for more accurate rate enforcement
// - Per-key rate limiting for multi-tenant scenarios
// - Middleware helpers for easy integration with web frameworks
//
// The package is structured to minimize lock contention under high concurrency
// and uses atomic operations where possible for performance.
//
// Usage examples:
//
//     // Create a rate limiter allowing 100 requests per second
//     limiter := ratelimit.NewTokenBucket(100, time.Second)
//
//     // Check if request is allowed
//     if limiter.Allow() {
//         // Process request
//     } else {
//         // Return 429 Too Many Requests
//     }
//
// Important notes:
// - All rate limiters in this package are safe for concurrent use
// - Rate limiters do not persist state; limits reset on restart
// - For distributed rate limiting, see the ratelimit/redis subpackage
//
// This package is part of the larger ecosystem and integrates
// with the metrics and logging packages for observability.
package ratelimit
```
