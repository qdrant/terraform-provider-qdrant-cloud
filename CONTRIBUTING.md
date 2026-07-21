## Welcome and thank you for your interest in contributing to this project ❤️🚀

If you find any issues or would like to contribute, feel free to create an issue or a pull request. In the next sections, you will find all the needed information in order to configure your environment and test your changes.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.7.x+
- [Go](https://golang.org/doc/install) 1.26+ (to build the provider plugin)

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

## Validation guidelines

Treat the Qdrant Cloud API as the source of truth and keep validation there by default.

Do:

- Add Terraform-side validation only when the rule is stable, static, and gives users useful feedback during `terraform plan` and speeds-up terraform execution. Protobuf enums are the common case.
- Use Terraform schema types and declarative behaviors such as `ConflictsWith`, `ExactlyOneOf`, and `AtLeastOneOf` before writing a custom validator.
- Use SDKv2's `ValidateDiagFunc` and `helper/validation` package. This provider uses `terraform-plugin-sdk/v2`.
- Derive allowed values from generated protobuf stubs.
- Exclude the zero-valued generated enum entry, such as `STORAGE_TIER_TYPE_UNSPECIFIED`. It represents "not set," not a valid explicit Terraform value.
- Sort generated enum names before using them in descriptions or validation errors so the output is deterministic.
- Attach validators only to user-configurable resource fields. Do not validate computed-only data-source fields. The latter is not under control by infrastructure-as-code and leaves terraform code in potentially stuck and broken state.
- Keep validators pure, deterministic, and offline. Do not read files, environment variables, time, provider state, or remote services.
- Test accepted values, rejected values, and excluded generated enum entries.

Don't:

- Copy allowed values into a hand-maintained slice.
- Duplicate dynamic, account-specific, regional, or cross-field API rules in the provider.
- Make API calls from schema validators.
- Use `terraform-plugin-framework` validators in an SDKv2 schema. Use them only after migrating that resource or data source to the Plugin Framework.

### Good: derive a stable enum from generated code

```go
func protoEnumNames(enumNames map[int32]string) []string {
	values := make([]string, 0, len(enumNames))
	for number, name := range enumNames {
		if number != 0 {
			values = append(values, name)
		}
	}
	sort.Strings(values)
	return values
}

validStorageTiers := protoEnumNames(commonv1.StorageTierType_name)
storageTierType.ValidateDiagFunc = validation.ToDiagFunc(
	validation.StringInSlice(validStorageTiers, false),
)
```

This remains aligned with the API when the generated protobuf dependency changes.

### Bad: maintain a second enum list

```go
storageTierType.ValidateDiagFunc = validation.ToDiagFunc(
	validation.StringInSlice([]string{
		"STORAGE_TIER_TYPE_COST_OPTIMISED",
		"STORAGE_TIER_TYPE_BALANCED",
		"STORAGE_TIER_TYPE_PERFORMANCE",
	}, false),
)
```

This duplicates the API contract and can silently drift from it.

### Good: leave dynamic rules to the API

Define `package_id` as a required string and let the API decide whether that package is available for the selected account, cloud, and region.

### Bad: duplicate dynamic API state

Do not hardcode currently available package IDs or query the API from a schema validator. Availability changes independently of the provider, so such validation becomes stale or makes planning depend on live API state.

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
