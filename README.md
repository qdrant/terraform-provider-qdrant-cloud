# Qdrant Cloud Terraform Provider

This is a Terraform provider for Qdrant Cloud, which is the DBaaS solution for Qdrant database, which is a vector similarity search engine with extended functionality. The provider allows you to manage your Qdrant Cloud resources using Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.7.x+
- [Go](https://golang.org/doc/install) 1.24+ (to build the provider plugin)

## Building & Using the Provider (Local Development)

### 1. Build and install locally

The fastest way to build the provider and make it discoverable by Terraform:

```bash
# Build with go build and copy into
# ~/.terraform.d/plugins/local/qdrant-cloud/qdrant-cloud/1.0/<os>_<arch>/terraform-provider-qdrant-cloud
make local-build
```

### 2. Use the local plugin in Terraform

In your Terraform project, reference the locally installed provider:

```hcl
terraform {
  required_providers {
    qdrant-cloud = {
      source  = "local/qdrant-cloud"
      version = "1.0"
    }
  }
}

provider "qdrant-cloud" {
  api_key    = "<your_api_key>"
  account_id = "<your_account_id>"
  api_url    = "<grpc_server_url>"
}

```

Then initialize:

```bash
terraform init
terraform plan
```

### 3. Use local plugin with `dev_overrides`

If you want to run the examples in `./examples` without changing their `source = "qdrant/qdrant-cloud"` blocks,  
you can override the provider installation via a Terraform CLI config file.

Create `~/.terraformrc` with:

```hcl
provider_installation {
  dev_overrides {
    "qdrant/qdrant-cloud" = "~/.terraform.d/plugins/local/qdrant-cloud/qdrant-cloud/1.0/<os>_<arch>"
  }
  direct {}
}
```

## Testing

In order to test the provider, you can run `make test`.

This will run the unit & acceptance tests in the provider.

```bash
make test \
  QDRANT_CLOUD_API_KEY="<API_KEY>" \
  QDRANT_CLOUD_ACCOUNT_ID="<ACCOUNT_ID>" \
  QDRANT_CLOUD_API_URL="grpc.development-cloud.qdrant.io" \
  SKIPPED_TESTS="TestAccDataAccountsBackupSchedule TestAccResourceClusterCreate"
```

run only unit tests

```bash

make test.unit
```

run only acceptance tests

```bash

make test.acceptance \
  QDRANT_CLOUD_API_KEY="<API_KEY>" \
  QDRANT_CLOUD_ACCOUNT_ID="<ACCOUNT_ID>" \
  QDRANT_CLOUD_API_URL="grpc.development-cloud.qdrant.io" \
  SKIPPED_TESTS="TestAccDataAccountsBackupSchedule TestAccResourceClusterCreate"
```

or run a single **unit** test without `make` (acceptance tests still require the env vars from above);

```bash
go test -v ./... -run '^TestFlattenHCEnv$'
```

to run a single **acceptance** test, export `TF_ACC=1` plus the required Qdrant credentials and target the specific `TestAcc*`:

```bash
TF_ACC=1 \
  QDRANT_CLOUD_API_KEY="<API_KEY>" \
  QDRANT_CLOUD_ACCOUNT_ID="<ACCOUNT_ID>" \
  QDRANT_CLOUD_API_URL="grpc.development-cloud.qdrant.io" \
  go test -count=1 -v ./qdrant -run '^TestAccResourceAccountsUserRoles_Update$'
```

## Releasing

In order to release the provider (available for maintainers only):

- Go to the [releases](https://github.com/qdrant/terraform-provider-qdrant-cloud/releases) on GitHib
- Edit the 'DRAFT' (see pencil)
- Typing a 'Tag', e.g. v1.3.2 (note that we are using [semantic versioning](https://semver.org/))
- Edit the 'Release Title' with the version you want to release (like v1.3.2)
- Click 'Publish release' button

This will automatically create the artifacts and place them in the assets, create a tag and invoke a webhook in the terraform provider registry.
After a while (approx 5min) the release will be available at the [terraform provider registry](https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest) as latest release.

## Contributing

If you find any issues or would like to contribute, feel free to create an issue or a pull request.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
