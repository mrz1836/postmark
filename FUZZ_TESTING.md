# Fuzz Testing for Postmark Go Library

This library includes comprehensive fuzz tests to help identify security vulnerabilities, edge cases, and robustness issues. The fuzz tests use Go's native fuzzing framework introduced in Go 1.18.

## Overview

The fuzz tests cover critical areas of the library that handle external input:

- **JSON unmarshaling and error handling** (`postmark_fuzz_test.go`)
- **Email address validation and processing** (`email_fuzz_test.go`)
- **Bounce data parsing and query parameters** (`bounce_fuzz_test.go`)
- **Template validation and processing** (`templates_fuzz_test.go`)

## Running Fuzz Tests

### Run All Fuzz Tests (Short Duration)
```bash
go test -fuzz=. -fuzztime=10s
```

### Run Specific Fuzz Test
```bash
go test -fuzz=FuzzDoRequestJSONUnmarshal -fuzztime=30s
go test -fuzz=FuzzEmailAddressValidation -fuzztime=30s
go test -fuzz=FuzzGetBouncedTagsJSON -fuzztime=30s
go test -fuzz=FuzzTemplateValidation -fuzztime=30s
```

### Run Only Fuzz Tests (Skip Regular Tests)
```bash
go test -run=^$ -fuzz=FuzzDoRequestJSONUnmarshal -fuzztime=60s
```

### Extended Fuzzing Session
```bash
go test -fuzz=FuzzEmailAddressValidation -fuzztime=5m
```

## Fuzz Test Details

### 1. JSON Processing (`postmark_fuzz_test.go`)

- **FuzzDoRequestJSONUnmarshal**: Tests JSON unmarshaling with malformed responses
- **FuzzAPIErrorHandling**: Tests error response parsing with various error formats
- **FuzzJSONPayloadMarshaling**: Tests request payload marshaling

**What it finds:**
- JSON injection vulnerabilities
- Unmarshaling panics with malformed data
- Memory exhaustion from large payloads
- Error message information leaks

### 2. Email Processing (`email_fuzz_test.go`)

- **FuzzEmailAddressValidation**: Tests email address handling across all fields
- **FuzzEmailHeaders**: Tests custom header processing
- **FuzzEmailAttachments**: Tests file attachment handling
- **FuzzEmailBatch**: Tests batch email processing

**What it finds:**
- Email header injection attempts
- Path traversal in attachment filenames
- Base64 decoding issues
- Email address format edge cases
- DoS via oversized payloads

### 3. Bounce Processing (`bounce_fuzz_test.go`)

- **FuzzGetBouncedTagsJSON**: Tests the custom JSON parsing logic for bounce tags
- **FuzzBounceQueryParams**: Tests URL parameter encoding
- **FuzzBounceQueryParamsInjection**: Tests for injection vulnerabilities in query params
- **FuzzBounceJSONStructure**: Tests bounce data structure parsing

**What it finds:**
- Query parameter injection
- URL encoding edge cases
- JSON structure manipulation
- XSS in query parameters

### 4. Template Processing (`templates_fuzz_test.go`)

- **FuzzTemplateValidation**: Tests template content validation
- **FuzzTemplateQueryParams**: Tests template filtering parameters
- **FuzzTemplatedEmail**: Tests templated email sending
- **FuzzTemplateBatch**: Tests batch templated email processing

**What it finds:**
- Template injection vulnerabilities
- Server-side template injection (SSTI)
- Script injection in template content
- Template model manipulation

## Security Focus Areas

The fuzz tests specifically look for:

1. **Injection Attacks**
   - Header injection (`\r\n`)
   - Script injection (`<script>`, `javascript:`)
   - Template injection (`{{`, `<%`, `${`)
   - Path traversal (`../`, `..\\`)

2. **Data Validation**
   - Email address formats (RFC 5321 compliance)
   - Base64 encoding validation
   - JSON structure validation
   - Parameter length limits

3. **Error Handling**
   - Information leaks in error messages
   - Graceful handling of malformed input
   - No panics with invalid data
   - Appropriate error boundaries

4. **DoS Prevention**
   - Large payload handling
   - Memory exhaustion protection
   - Excessive parameter counts
   - Long-running operations

## Interpreting Results

### Success
When fuzz tests pass without finding issues:
```
PASS
ok      github.com/mrz1836/postmark    5.311s
```

### Found Issues
Fuzz tests will report specific security concerns as log messages:
```
email_fuzz_test.go:65: Potential header injection detected in email field
bounce_fuzz_test.go:95: Script injection attempt in query parameters
templates_fuzz_test.go:89: Potential template injection pattern: {{exec
```

These are **expected findings** - the fuzz tests are working correctly by detecting potential security issues.

### Failures
Actual test failures indicate problems in the fuzzing logic itself:
```
FAIL: FuzzEmailAddressValidation (0.01s)
    email_fuzz_test.go:X: Unexpected panic during email processing
```

## Continuous Integration

Add fuzz testing to your CI pipeline:

```yaml
- name: Run Fuzz Tests
  run: |
    go test -fuzz=. -fuzztime=30s
```

For more comprehensive testing, run longer fuzz sessions periodically:

```yaml
- name: Extended Fuzz Testing
  run: |
    go test -fuzz=FuzzDoRequestJSONUnmarshal -fuzztime=5m
    go test -fuzz=FuzzEmailAddressValidation -fuzztime=5m
    go test -fuzz=FuzzGetBouncedTagsJSON -fuzztime=5m
    go test -fuzz=FuzzTemplateValidation -fuzztime=5m
```

## Best Practices

1. **Regular Execution**: Run fuzz tests regularly, not just before releases
2. **Seed Corpus**: The tests include comprehensive seed data for effective fuzzing
3. **Time Investment**: Longer fuzz sessions find more edge cases
4. **Multiple Targets**: Fuzz different functions to get comprehensive coverage
5. **Monitor Resources**: Long fuzz sessions can be resource-intensive

## Contributing

When adding new functionality to the library:

1. Add corresponding fuzz tests for any new input parsing
2. Include relevant seed data in the corpus
3. Test for common injection patterns
4. Validate error handling doesn't leak information
5. Ensure graceful handling of malformed input

The fuzz tests are designed to be maintainable and comprehensive, helping ensure the library remains secure and robust against malicious input.