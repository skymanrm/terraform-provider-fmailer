# Publishing Your FMailer Terraform Provider

This guide covers the steps to properly publish your Terraform Provider to the Terraform Registry.

## Prerequisites

Before you can publish your provider, ensure you have:

1. A GitHub account
2. The provider code in a public GitHub repository
3. A [Terraform Registry](https://registry.terraform.io) account linked to your GitHub account
4. GPG keys for signing releases

## Step 1: Prepare Your Repository

Your repository should follow the naming convention `terraform-provider-{NAME}`, where `{NAME}` is the provider name (in this case, `terraform-provider-fmailer`).

Ensure your repository has:

- A well-structured README.md with documentation
- Proper license file
- Complete Go modules configuration
- Examples directory

## Step 2: Versioning

Terraform providers follow [Semantic Versioning](https://semver.org/). When ready to release:

1. Tag your release with a version number prefixed with `v`, e.g., `v1.0.0`
2. Sign your release tag with GPG

```bash
git tag -a v1.0.0 -m "First release"
git push origin v1.0.0
```

For a signed tag:

```bash
git tag -s v1.0.0 -m "First release"
git push origin v1.0.0
```

## Step 3: GitHub Release

1. Go to your repository on GitHub
2. Navigate to "Releases"
3. Click "Draft a new release"
4. Select your version tag
5. Add release notes describing changes, features, and bug fixes
6. Publish the release

## Step 4: Register Your Provider

1. Log in to the [Terraform Registry](https://registry.terraform.io)
2. Go to "Publish" and select "Provider"
3. Choose your GitHub repository
4. The registry will verify your repository meets the requirements
5. Complete the provider information:
    - Namespace (usually your GitHub username)
    - Provider name (`fmailer`)
    - Category (e.g., "Communication")
    - Description
    - Documentation links

## Step 5: Provider Documentation

For optimal usability, ensure your documentation includes:

- Provider configuration options
- Resource and data source details
- Examples for each resource type
- Attribute descriptions

## Step 6: Goreleaser Configuration (Optional)

To automate builds for multiple platforms, add a `.goreleaser.yml` file:

```yaml
before:
  hooks:
    - go mod tidy
builds:
- env:
    - CGO_ENABLED=0
  mod_timestamp: '{{ .CommitTimestamp }}'
  flags:
    - -trimpath
  ldflags:
    - '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}}'
  goos:
    - freebsd
    - windows
    - linux
    - darwin
  goarch:
    - amd64
    - '386'
    - arm
    - arm64
  ignore:
    - goos: darwin
      goarch: '386'
  binary: '{{ .ProjectName }}_v{{ .Version }}'
archives:
- format: zip
  name_template: '{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}'
checksum:
  name_template: '{{ .ProjectName }}_{{ .Version }}_SHA256SUMS'
  algorithm: sha256
signs:
  - artifacts: checksum
    args:
      - "--batch"
      - "--local-user"
      - "{{ .Env.GPG_FINGERPRINT }}"
      - "--output"
      - "${signature}"
      - "--detach-sign"
      - "${artifact}"
release:
  github:
    owner: yourusername
    name: terraform-provider-fmailer
```

Then you can use GoReleaser to publish:

```bash
export GPG_FINGERPRINT=your-gpg-fingerprint
goreleaser release --rm-dist
```

## Step 7: Terraform Provider Development Program (Optional)

Consider joining the [Terraform Provider Development Program](https://www.terraform.io/docs/registry/providers/publishing.html) for additional benefits like:

- Verified provider badge
- Provider documentation hosting
- Provider development resources

## Updating Your Provider

When you need to update your provider:

1. Make your code changes
2. Update the version in your code (if needed)
3. Tag a new release following semantic versioning
4. Create a new GitHub release
5. The Terraform Registry will automatically index the new version

## Testing Before Publishing

Before final publication, test your provider locally:

1. Build the provider locally
2. Configure Terraform to use your local provider:

```hcl
# development.tfrc
provider_installation {
  dev_overrides {
    "yourusername/fmailer" = "/path/to/your/build/directory"
  }
  direct {}
}
```

3. Set the `TF_CLI_CONFIG_FILE` environment variable:

```bash
export TF_CLI_CONFIG_FILE=/path/to/development.tfrc
```

4. Test your configuration