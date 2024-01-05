# Terraform Provider for MAAS

This repository contains the source code for the Terraform MAAS provider.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.4.x
- [Go](https://golang.org/doc/install) >= 1.20

## Build The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider with:

    ```sh
    make build
    ```

1. (Optional): Install the freshly built provider with:

    ```sh
    make install
    ```

## Usage

### Provider Configuration

The provider accepts the following config options:

- **api_key**: [MAAS API key](https://maas.io/docs/snap/3.0/cli/maas-cli#heading--log-in-required).
- **api_url**: URL for the MAAS API server (eg: <http://127.0.0.1:5240/MAAS>).
- **api_version**: MAAS API version used. It is optional and it defaults to `2.0`.

#### `maas`

```hcl
provider "maas" {
  api_version = "2.0"
  api_key = "YOUR MAAS API KEY"
  api_url = "http://<MAAS_SERVER>[:MAAS_PORT]/MAAS"
}
```

### Resources and Data Sources Configuration

The [docs](/docs) section contains details about each supported Terraform resource and data source.

### Release process

1. Create a new branch from `master` as `release-vX.X.X`
   ```bash
   git branch release-vX.X.X master
   git push -u origin release-vX.X.X
   ```
2. Update the `CHANGELOG.md` with your release version, date and change details and push the changes to the new branch.
3. Update `Makefile` with the new release version.
4. Raise a PR on Github, title of the PR should be in the following format
   `Release vX.X.X`
5. Merge the PR into master, taking a note of the merge commit which is created
6. Go to [Releases](https://github.com/maas/terraform-provider-maas/releases) on Github
7. Click [Draft a new release](https://github.com/maas/terraform-provider-maas/releases/new)
8. On `Target` choose the latest merge commit you want to release for
9. Set the `tag` to create a new tag for the version in the format "vX.X.X"
10. Set the `release title` to the release version, for example `v1.0.1`
11. Copy and paste the relevant CHANGELOG.md entries to the release description
12. Click `Publish release`
13. The new version should be available on the [Releases](https://github.com/maas/terraform-provider-maas/releases) page

## Additional Info

### Testing

Unit tests run with every pull request and merge to master. The end to end tests run on a nightly basis against a hosted MAAS deployment, results can be found [here](https://raw.githubusercontent.com/maas/maas-terraform-e2e-tests/main/results.json?token=GHSAT0AAAAAAB3FX6R5C67Q4LH7ADOO5O3IY4ODCNA) and are checked on each PR, with a warning if failed.   
