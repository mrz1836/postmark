# Technical Conventions

> This directory contains modular, portable technical conventions that can be adopted across projects. Each file focuses on a specific aspect of software development.

<br><br>

## 📚 Convention Categories

### 🚀 Core Development

**[Go Essentials](go-essentials.md)**
Non-negotiable Go development practices including context-first design, interface philosophy, goroutine discipline, error handling, and performance guidelines.

**[Testing Standards](testing-standards.md)**
Comprehensive testing guidelines covering unit tests, table-driven tests, fuzz testing, code coverage requirements, and testing best practices.

**[Commenting & Documentation](commenting-documentation.md)**
Standards for code comments, function documentation, package-level docs, Markdown formatting, and maintaining clear, purposeful documentation.

<br><br>

### 🔄 Version Control & Collaboration

**[Commit & Branch Conventions](commit-branch-conventions.md)**
Git commit message format, branch naming standards, and version control best practices for maintaining clean repository history.

**[Pull Request Guidelines](pull-request-guidelines.md)**
Structured approach to creating and reviewing pull requests, including required sections, review etiquette, and merging strategies.

**[Release Workflow & Versioning](release-versioning.md)**
Semantic versioning practices, release tooling with goreleaser, changelog management, and automated release processes.

<br><br>

### 🏷️ Project Management

**[Labeling Conventions](labeling-conventions.md)**
GitHub label system for categorizing issues and PRs, including standard labels, usage guidelines, and automated labeling.

<br><br>

### 🔧 Infrastructure & Quality

**[CI & Validation](ci-validation.md)**
Continuous integration setup, automated checks, required status checks, and troubleshooting CI failures.

**[Dependency Management](dependency-management.md)**
Go modules management, security scanning, version control practices, and maintaining healthy dependencies.

**[Security Practices](security-practices.md)**
Security-first development, vulnerability reporting, security tools, and following OpenSSF best practices.

**[Pre-commit Hooks](pre-commit.md)**
Pure Go pre-commit framework with 17x faster execution than Python alternatives, providing automated code quality checks.

**[GitHub Workflows Development](github-workflows.md)**
Creating and maintaining GitHub Actions workflows with security, reliability, and performance in mind.

<br><br>

### 🏗️ Build & Project Setup

**[MAGE-X Build Automation](mage-x.md)**
Zero-boilerplate build automation system with 240+ built-in commands that replaces Makefiles. Includes installation, configuration, command reference, and migration guide.

**[Governance Documents](governance-documents.md)**
Essential project files including LICENSE, README, CODE_OF_CONDUCT, SECURITY, and other governance standards.

<br><br>

## 🎯 Using These Conventions

### For New Projects
1. Copy this entire `tech-conventions` directory to your `.github` folder
2. Review each file and adjust project-specific details
3. Reference these conventions in your main documentation

### For Existing Projects
1. Review current practices against these conventions
2. Adopt conventions incrementally
3. Update team documentation to reference relevant sections

### Customization
* These conventions are designed to be forked and modified
* Maintain the structure and formatting for consistency
* Document any project-specific deviations

<br><br>

## 📋 Quick Reference

| Need to...             | See...                                                      |
|------------------------|-------------------------------------------------------------|
| Write Go code          | [Go Essentials](go-essentials.md)                           |
| Create tests           | [Testing Standards](testing-standards.md)                   |
| Document code          | [Commenting & Documentation](commenting-documentation.md)   |
| Make commits           | [Commit & Branch Conventions](commit-branch-conventions.md) |
| Open a PR              | [Pull Request Guidelines](pull-request-guidelines.md)       |
| Tag a release          | [Release Workflow & Versioning](release-versioning.md)      |
| Label issues/PRs       | [Labeling Conventions](labeling-conventions.md)             |
| Fix CI issues          | [CI & Validation](ci-validation.md)                         |
| Add dependencies       | [Dependency Management](dependency-management.md)           |
| Report security issues | [Security Practices](security-practices.md)                 |
| Create workflows       | [GitHub Workflows Development](github-workflows.md)         |
| Set up build automation| [MAGE-X Build Automation](mage-x.md)                        |
| Add governance docs    | [Governance Documents](governance-documents.md)             |

<br><br>

## 🤝 Contributing

These conventions evolve with best practices. To propose changes:

1. Fork the repository
2. Create a feature branch (`feat/improve-testing-standards`)
3. Make your changes with clear reasoning
4. Submit a PR with the "documentation" label

Remember: Good conventions are discovered through practice, not prescribed in theory.
