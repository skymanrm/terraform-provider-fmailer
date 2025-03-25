# Terraform Provider for FMailer

This Terraform provider allows you to manage resources within the FMailer email sending service through your Terraform configuration.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
```sh
git clone https://github.com/yourusername/terraform-provider-fmailer
```

2. Enter the repository directory
```sh
cd terraform-provider-fmailer
```

3. Build the provider
```sh
go build -o terraform-provider-fmailer
```

## Installing the Provider

### Terraform 0.13+

1. Build the provider as described above
2. Create a `~/.terraform.d/plugins/registry.terraform.io/yourusername/fmailer/1.0.0/linux_amd64/` directory (adjust the OS/architecture as needed)
3. Copy the provider binary to this directory
4. Use the provider in your Terraform configuration

### Terraform Configuration Example

```hcl
terraform {
  required_providers {
    fmailer = {
      source = "yourusername/fmailer"
      version = "1.0.0"
    }
  }
}

provider "fmailer" {
  token = var.fmailer_token
  endpoint = "https://api.fmailer.com" # optional, default value
}

# Create a domain template
resource "fmailer_domain_template" "example" {
  name   = "Welcome Email"
  slug   = "welcome-email"
  domain = 123

  langs {
    lang     = "en"
    subject  = "Welcome to our service!"
    body     = "Hello {{name}}, thank you for signing up."
    default  = true
  }

  langs {
    lang     = "fr"
    subject  = "Bienvenue à notre service!"
    body     = "Bonjour {{name}}, merci de vous être inscrit."
    default  = false
  }
}

# Use a domain template data source
data "fmailer_domain_template" "existing" {
  uuid = "123e4567-e89b-12d3-a456-426614174000"
}

output "template_name" {
  value = data.fmailer_domain_template.existing.name
}
```

## Authentication

The provider requires a token for authentication which can be provided in several ways:

1. Via the provider configuration:
```hcl
provider "fmailer" {
  token = "your-token-here"
}
```

2. Via the environment variable:
```sh
export FMAILER_TOKEN="your-token-here"
```

## Resources

### `fmailer_domain_template`

A domain template resource allows you to create, update, and delete email templates.

#### Example Usage

```hcl
resource "fmailer_domain_template" "welcome" {
  name   = "Welcome Email"
  slug   = "welcome-email"
  domain = 123

  langs {
    lang     = "en"
    subject  = "Welcome to our service!"
    body     = "Hello {{name}}, thank you for signing up."
    default  = true
  }
}
```

#### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the template.
* `slug` - (Required) The slug for the template. Must only contain alphanumeric characters, hyphens, and underscores.
* `domain` - (Required) The ID of the domain this template belongs to.
* `allow_copy` - (Optional) Whether this template can be copied. Defaults to `true`.
* `editable` - (Optional) Whether this template can be edited. Defaults to `true`.
* `langs` - (Optional) A list of language-specific template contents. Each language block supports:
* `lang` - (Required) The language code (e.g., "en", "fr").
* `subject` - (Required) The email subject line.
* `body` - (Required) The email body content.
* `default` - (Optional) Whether this is the default language. Defaults to `false`.

#### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - The unique identifier for the template.
* `created_at` - The timestamp when the template was created.
* `updated_at` - The timestamp when the template was last updated.

## Data Sources

### `fmailer_domain_template`

Retrieve information about an existing domain template.

#### Example Usage

```hcl
data "fmailer_domain_template" "existing" {
  uuid = "123e4567-e89b-12d3-a456-426614174000"
}

output "template_name" {
  value = data.fmailer_domain_template.existing.name
}
```

#### Argument Reference

The following arguments are supported:

* `uuid` - (Required) The unique identifier of the template.

#### Attribute Reference

The following attributes are exported:

* `name` - The name of the template.
* `slug` - The slug for the template.
* `domain` - The ID of the domain this template belongs to.
* `allow_copy` - Whether this template can be copied.
* `editable` - Whether this template can be edited.
* `created_at` - The timestamp when the template was created.
* `updated_at` - The timestamp when the template was last updated.
* `langs` - A list of language-specific template contents:
* `lang` - The language code.
* `subject` - The email subject line.
* `body` - The email body content.
* `default` - Whether this is the default language.

## Development

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

### Building

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `build` command:
```sh
go build -o terraform-provider-fmailer
```

## License

[MIT](LICENSE)