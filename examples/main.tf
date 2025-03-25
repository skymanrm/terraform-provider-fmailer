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
  endpoint = var.fmailer_endpoint
}

variable "fmailer_token" {
  type = string
  description = "FMailer API token"
  sensitive = true
}

variable "fmailer_endpoint" {
  type = string
  description = "FMailer API endpoint"
  default = "https://api.fmailer.com"
}

resource "fmailer_domain_template" "welcome_email" {
  name   = "Welcome Email"
  slug   = "welcome-email"
  domain = 123

  langs {
    lang     = "en"
    subject  = "Welcome to our service!"
    body     = <<-EOT
      Hello {{name}},

      Thank you for signing up to our service. We're excited to have you on board!

      Best regards,
      The Team
    EOT
    default  = true
  }

  langs {
    lang     = "fr"
    subject  = "Bienvenue à notre service!"
    body     = <<-EOT
      Bonjour {{name}},

      Merci de vous être inscrit à notre service. Nous sommes ravis de vous avoir à bord!

      Cordialement,
      L'équipe
    EOT
    default  = false
  }
}

output "template_id" {
  value = fmailer_domain_template.welcome_email.uuid
}

# Data source example
data "fmailer_domain_template" "existing" {
  uuid = fmailer_domain_template.welcome_email.uuid
}

output "template_name" {
  value = data.fmailer_domain_template.existing.name
}