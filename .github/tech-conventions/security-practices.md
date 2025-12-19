# Security Practices

> Security is a first-class requirement. This document outlines security practices, vulnerability reporting, and tools used to maintain a secure codebase.

<br><br>

## 🛡️ Vulnerability Reporting

If you discover a vulnerability—no matter how small—follow our responsible disclosure process:

* **Do not** open a public issue or pull request.
* Follow the instructions in [`SECURITY.md`](../SECURITY.md).
* Include:
  * A clear, reproducible description of the issue
  * Proof‑of‑concept code or steps (if possible)
  * Any known mitigations or workarounds
* You will receive an acknowledgment within **72 hours** and status updates until the issue is resolved.

> For general hardening guidance (e.g., `govulncheck`, dependency pinning), see the [Dependency Management](dependency-management.md) section.

<br><br>

## 🔐 Security Tools

### Required Security Scans

1. **govulncheck** - Go vulnerability database scanning
   ```bash
   magex deps:audit
   ```

2. **gitleaks** - Secret detection in code
   ```bash
   gitleaks detect --source . --log-opts="--all" --verbose
   ```

3. **CodeQL** - Semantic code analysis (runs in CI)
   - Automated via `.github/workflows/codeql-analysis.yml`
   - Scans for common vulnerabilities

<br><br>

## 🚫 Security Anti-Patterns

### Never Do This:

```go
// 🚫 Never hardcode secrets
// apiKey := "1234..."

// 🚫 Never log sensitive data
// log.Printf("User password: %s", password)

// 🚫 Never use weak cryptography
hash := md5.Sum([]byte(data)) // MD5 is broken

// 🚫 Never trust user input without validation
query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)

// 🚫 Never ignore security errors
cert, _ := tls.LoadX509KeyPair(certFile, keyFile) // Always check errors!
```

### Always Do This:

```go
// ✅ Use environment variables for secrets
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    return errors.New("API_KEY environment variable not set")
}

// ✅ Sanitize logs
log.Printf("User authentication attempt for ID: %s", userID)

// ✅ Use strong cryptography
hash := sha256.Sum256([]byte(data))

// ✅ Use parameterized queries
query := "SELECT * FROM users WHERE id = ?"
rows, err := db.QueryContext(ctx, query, userInput)

// ✅ Always handle security-critical errors
cert, err := tls.LoadX509KeyPair(certFile, keyFile)
if err != nil {
    return fmt.Errorf("failed to load TLS certificate: %w", err)
}
```

<br><br>

## 🔒 Secure Coding Practices

### Input Validation
* **Validate all inputs** at trust boundaries
* **Use allowlists** over denylists when possible
* **Sanitize** before using in queries, commands, or output
* **Set limits** on input size and complexity

### Authentication & Authorization
* **Use standard libraries** for crypto operations
* **Never roll your own crypto**
* **Store passwords** using bcrypt, scrypt, or argon2
* **Implement proper session management**
* **Use constant-time comparisons** for secrets

### Error Handling
* **Don't leak sensitive info** in error messages
* **Log security events** for monitoring
* **Fail securely** - deny by default
* **Handle panics** in goroutines

<br><br>

## 📋 Security Checklist

Before committing code, verify:

- [ ] No hardcoded secrets or credentials
- [ ] All user inputs are validated
- [ ] SQL queries use parameters, not string concatenation
- [ ] File paths are sanitized before use
- [ ] Proper error handling without info leakage
- [ ] Dependencies are up to date
- [ ] Security scans pass (govulncheck, gitleaks)

<br><br>

## 🚨 Incident Response

If a security issue is found in production:

1. **Don't panic** - Follow the process
2. **Assess severity** using CVSS scoring
3. **Notify security team** immediately
4. **Create private fix** in security fork
5. **Test thoroughly** including regression tests
6. **Coordinate disclosure** with security team
7. **Release patch** with security advisory
8. **Monitor** for exploitation attempts

<br><br>

## 🏗️ Security by Design

### Principle of The Least Privilege
* Run processes with minimal permissions
* Use read-only file systems where possible
* Drop privileges after initialization
* Segment access by service boundaries

### Defense in Depth
* Multiple layers of security controls
* Don't rely on a single security measure
* Assume other defenses may fail
* Monitor and alert on anomalies

### Secure Defaults
* Deny by default, allow explicitly
* Require secure configuration
* Force HTTPS/TLS connections
* Enable security features by default

<br><br>

## 📊 OpenSSF Best Practices

We follow [OpenSSF](https://openssf.org) guidelines:

1. **Vulnerability Disclosure** - Clear security policy
2. **Dependency Maintenance** - Regular updates
3. **Security Testing** - Automated scanning
4. **Access Control** - Protected branches
5. **Build Integrity** - Signed releases
6. **Code Review** - Required for all changes

### Scorecard Compliance
Monitor security posture via:
```bash
# Check OpenSSF Scorecard
scorecard --repo=github.com/owner/repo
```

<br><br>

## 🔍 Security Reviews

### When to Request Review
* Cryptographic implementations
* Authentication/authorization changes
* Handling sensitive data
* Network-facing services
* File system operations
* Subprocess execution

### Review Focus Areas
* Input validation completeness
* Output encoding correctness
* Error handling safety
* Resource limit enforcement
* Permission boundaries
* Audit logging coverage

<br><br>

## 📚 Security Resources

### Internal Resources
* [`SECURITY.md`](../SECURITY.md) - Vulnerability reporting
* [Dependency Management](dependency-management.md) - Supply chain security

### External Resources
* [OWASP Top 10](https://owasp.org/Top10/)
* [Go Security Best Practices](https://golang.org/doc/security)
* [CWE Database](https://cwe.mitre.org/)
* [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
