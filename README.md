# git-credential-github-app

[Git credential helper](https://git-scm.com/docs/gitcredentials) to use GitHub App key authentication, either locally or via an imported [Google Cloud KMS key](https://cloud.google.com/kms/docs/key-management-service), for `git` operations (e.g. `git clone`).

## Installation

Download one or all helpers from the [releases page](https://github.com/ericnorris/git-credential-github-app/releases) to a location in your `PATH`.

## Prerequisites

Besides installation, you must:

1. [Set up a GitHub App](https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app).
2. Retrieve your application's **Client ID**, which you can find in your application's [settings page](https://docs.github.com/en/apps/maintaining-github-apps/modifying-a-github-app-registration#navigating-to-your-github-app-settings).
3. Determine the application's **Installation ID**. In the "Install App" section of your application's settings page, you can find this in the URL of the installation (e.g. `https://github.com/organizations/<organization>/settings/installations/<id>` for an application installed to an organization).

## Usage

These binaries are not meant for direct use; instead, you should configure `git` to call them as a credential helper. You may also combine these with other credential helpers, as they will only respond to requests for "github.com" credentials.

Each example configuration below will include a `cache` helper as well to prevent unnecessary network traffic that may result in hitting rate limits.

### `github-app-credential-helper`

For a private key stored locally on disk, use the following git configuration:

```
[credential]
  helper = "cache --timeout=1800"
  helper = "github-app --private-key <path> -client-id <application client ID> --installation-id <id>"
```

### `github-app-credential-helper-gcpkms`

For a private key imported into Google Cloud KMS, use the following git configuration:

```
[credential]
  helper = "cache --timeout=1800"
  helper = "github-app-gcpkms --kms-key <id> -client-id <application client ID> --installation-id <id>"
```

`github-app-credential-helper-gcpkms` uses Google's [application default credentials](https://cloud.google.com/docs/authentication/application-default-credentials) pattern to authenticate with Google Cloud.

The caller must have `roles/cloudkms.viewer` and `roles/cloudkms.signer` on the KMS resource containing the imported key.

## Why Google Cloud KMS?

While `github-app-credential-helper` works fine with a private key, you must now distribute that key securely, as well as rotate it on some cadence.

If you instead [import the key into Google Cloud KMS](https://cloud.google.com/kms/docs/key-import) and throw away the key, you may use Google Cloud IAM to allow your infrastructure to _use_ the key without exposing the underlying key material.

This means that `github-app-credential-helper-gcpkms` and IAM is all you need to allow your infrastructure to perform `git` operations with GitHub.
