# Terraform Provider for FMailer - Developer Documentation

## Project Overview

This is a Terraform provider that enables infrastructure-as-code management of email templates within the FMailer email sending service. The provider allows users to create, read, update, and delete email domain templates through Terraform configurations.

**Repository:** `github.com/skymanrm/terraform-provider-fmailer`
**Version:** v1.0.2
**Go Version:** 1.19+
**Terraform Version:** >= 1.0
**License:** MIT

## Architecture

```
Terraform CLI
     ↓
main.go (Provider Plugin Server)
     ↓
Provider Configuration (token, endpoint)
     ↓
┌─────────────────────────────┐
│ resource_domain_template.go │  ← Resource/DataSource definitions
│ - CRUD operations           │
│ - Schema validation         │
└──────────────┬──────────────┘
               ↓
        client.go (API Client)
        - HTTP communication
        - JSON serialization
        - Error handling
               ↓
        FMailer API Service
        https://api.fmailer.com
```

## Project Structure

```
/
├── .github/workflows/          # CI/CD pipelines
│   ├── ci.yml                 # Build, test, and lint
│   └── release.yml            # Automated releases
├── examples/
│   └── main.tf                # Example Terraform configuration
├── client.go                  # FMailer API client (281 lines)
├── main.go                    # Provider entry point (54 lines)
├── resource_domain_template.go # Resource/DataSource impl (327 lines)
├── .goreleaser.yml            # Multi-platform build config
├── go.mod                     # Go module definition
├── readme.md                  # User documentation
├── publishing.md              # Publishing guide
└── release.md                 # Release automation guide
```

## Core Files

### main.go (54 lines)
**Purpose:** Terraform provider plugin entry point

**Key Functions:**
- `main()`: Serves the provider as a Terraform plugin
- `Provider()`: Defines provider schema with `token` (required, sensitive) and `endpoint` (optional)
- `providerConfigure()`: Initializes the FMailer API client

**Location:** `/home/user/terraform-provider-fmailer/main.go`

### client.go (281 lines)
**Purpose:** FMailer API client implementation

**Data Structures:**
- `Client`: HTTP client with 1-minute timeout
- `DomainTemplate`: Email template with metadata
- `DomainTemplateLang`: Language-specific content (subject, body)
- `PaginatedDomainTemplateList`: API response wrapper

**API Methods:**
- `CreateDomainTemplate()`: POST `/api/domains/templates/`
- `GetDomainTemplate()`: GET `/api/domains/templates/{uuid}/`
- `UpdateDomainTemplate()`: PUT `/api/domains/templates/{uuid}/`
- `DeleteDomainTemplate()`: DELETE `/api/domains/templates/{uuid}/`
- `ListDomainTemplates()`: GET with filtering/pagination/search
- `DuplicateDomainTemplate()`: POST `/api/domains/templates/{uuid}/duplicate/`

**Authentication:** Bearer token via HTTP Authorization header

**Location:** `/home/user/terraform-provider-fmailer/client.go`

### resource_domain_template.go (327 lines)
**Purpose:** Terraform resource and data source implementations

**Resource:** `fmailer_domain_template`
- CRUD operations for email templates
- Schema validation (slug format: `^[-a-zA-Z0-9_]+$`)
- Multi-language support via nested `langs` blocks
- Attributes: `uuid`, `name`, `slug`, `domain`, `allow_copy`, `editable`

**Data Source:** `fmailer_domain_template`
- Read-only access to existing templates by UUID

**Location:** `/home/user/terraform-provider-fmailer/resource_domain_template.go`

## Key Features

### Multi-Language Template Support
Templates can include multiple language variants:

```hcl
resource "fmailer_domain_template" "example" {
  name   = "Welcome Email"
  slug   = "welcome-email"
  domain = "example.com"

  langs {
    lang     = "en"
    subject  = "Welcome!"
    body     = "Welcome to our service"
    default  = true
  }

  langs {
    lang    = "fr"
    subject = "Bienvenue!"
    body    = "Bienvenue dans notre service"
  }
}
```

### Supported Platforms
- **Operating Systems:** Linux, Windows, macOS, FreeBSD
- **Architectures:** amd64, 386, arm, arm64

## Development Workflows

### Building the Provider

```bash
# Build for current platform
go build -v ./...

# Build for all platforms (requires GoReleaser)
goreleaser build --snapshot --rm-dist
```

### Testing

```bash
# Run tests (when available)
go test -v ./...

# Run linting
golangci-lint run
```

**Note:** Currently no test files exist in the codebase.

### Local Development with Terraform

1. Build the provider:
   ```bash
   go build -o terraform-provider-fmailer
   ```

2. Create local override configuration in `~/.terraformrc`:
   ```hcl
   provider_installation {
     dev_overrides {
       "skymanrm/fmailer" = "/path/to/terraform-provider-fmailer"
     }
     direct {}
   }
   ```

3. Use the provider in your Terraform configuration (see `examples/main.tf`)

### Creating a Release

1. Ensure all changes are committed
2. Create and push a version tag:
   ```bash
   git tag v1.0.3
   git push origin v1.0.3
   ```
3. GitHub Actions will automatically:
   - Build for all platforms
   - Create GPG-signed checksums
   - Publish GitHub release with artifacts

## CI/CD Pipelines

### ci.yml - Continuous Integration
**Triggers:** Push to master, all pull requests

**Jobs:**
1. **Build & Test:** Verifies dependencies, builds binaries, runs tests
2. **Lint:** Runs golangci-lint for code quality

### release.yml - Release Automation
**Triggers:** Git tags matching `v*` pattern

**Process:**
1. Checkout with full history
2. Setup Go 1.19
3. Import GPG key from secrets
4. Run GoReleaser for multi-platform builds
5. Create GitHub release with signed artifacts

## API Client Design

The `Client` struct in `client.go` handles all HTTP communication with the FMailer API:

- **Base URL:** Configurable via provider (default: `https://api.fmailer.com`)
- **Authentication:** Bearer token in `Authorization` header
- **Timeout:** 1 minute for all requests
- **Content Type:** `application/json`
- **Error Handling:** HTTP status codes with descriptive error messages

## Schema Validation

### Slug Format
Slugs must match the regex: `^[-a-zA-Z0-9_]+$`
- Valid: `welcome-email`, `signup_template`, `email123`
- Invalid: `email template` (spaces), `email@template` (special chars)

### Required Fields
- `token`: Provider authentication token (sensitive)
- `name`: Human-readable template name
- `slug`: URL-safe template identifier
- `domain`: Email domain for the template

### Optional Fields
- `endpoint`: Custom API endpoint (defaults to production)
- `allow_copy`: Allow template duplication
- `editable`: Allow template editing

## Common Development Tasks

### Adding a New API Method

1. Add method to `Client` struct in `client.go`
2. Define request/response structures if needed
3. Implement HTTP request with proper error handling
4. Update resource/data source in `resource_domain_template.go` if needed

### Adding a New Resource/Data Source

1. Create schema definition using Terraform Plugin SDK
2. Implement CRUD functions (Create, Read, Update, Delete)
3. Add resource to provider in `main.go`
4. Update documentation in `readme.md`
5. Add example to `examples/` directory

### Debugging Tips

1. **Enable Terraform logging:**
   ```bash
   export TF_LOG=DEBUG
   terraform apply
   ```

2. **Test API client directly:**
   ```go
   client := &Client{
       BaseURL: "https://api.fmailer.com",
       Token:   "your-token",
   }
   template, err := client.GetDomainTemplate("uuid-here")
   ```

3. **Check HTTP requests:**
   - Review error messages from API responses
   - Validate JSON payloads
   - Verify authentication token

## Testing Strategy (Recommended)

Currently, the project has no test files. Consider adding:

### Unit Tests
- Test `Client` methods with mock HTTP server
- Test schema validation logic
- Test data transformations

### Acceptance Tests
- Use Terraform Plugin SDK testing framework
- Test full resource lifecycle (CRUD)
- Test import functionality
- Test error scenarios

Example test structure:
```go
func TestAccDomainTemplate_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckDomainTemplateDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccDomainTemplateConfig_basic,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckDomainTemplateExists("fmailer_domain_template.test"),
                ),
            },
        },
    })
}
```

## Dependencies

**Direct:**
- `github.com/hashicorp/terraform-plugin-sdk/v2` v2.31.0

**Standard Library:**
- `net/http`: HTTP client
- `encoding/json`: JSON serialization
- `fmt`: String formatting
- `context`: Request contexts
- `time`: Timeout handling

## Resources

- **User Documentation:** `readme.md`
- **Publishing Guide:** `publishing.md`
- **Release Guide:** `release.md`
- **Example Configuration:** `examples/main.tf`
- **Terraform Registry:** https://registry.terraform.io/providers/skymanrm/fmailer
- **Terraform Plugin SDK:** https://github.com/hashicorp/terraform-plugin-sdk

## Code Quality Standards

- Follow Go best practices and idioms
- Use meaningful variable and function names
- Add error handling for all operations
- Validate inputs before API calls
- Use contexts for cancellable operations
- Keep functions focused and single-purpose
- Document exported functions and types

## Version History

- **v1.0.2:** Debug CI improvements
- **v1.0.1:** Initial release with bug fixes
- **v1.0.0:** Initial stable release

## Contact & Support

- **Repository:** https://github.com/skymanrm/terraform-provider-fmailer
- **Issues:** GitHub Issues
- **License:** MIT (Copyright 2025 Andrey Fanyagin)

---

*This documentation is intended for developers working on the terraform-provider-fmailer project. For user-facing documentation, see readme.md.*
