# Qdrant Cloud Terraform Provider

This is a Terraform provider for Qdrant Cloud, which is the DBaaS solution for Qdrant database, which is a vector similarity search engine with extended functionality. The provider allows you to manage your Qdrant Cloud resources using Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.7.x+
- [Go](https://golang.org/doc/install) 1.22+ (to build the provider plugin)
- [`swagger-codegen`](https://swagger.io/tools/swagger-codegen/)
  `brew install swagger-codegen`

## Building The Provider

Clone the repository:

```bash
git clone git@github.com:<your_org>/terraform-provider-qdrant-cloud.git
```

Enter the provider directory and build the provider:

```bash
cd terraform-provider-qdrant-cloud
go build
```

## Using The Provider

If you're building the provider, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). After placing it into your plugins directory, run `terraform init` to initialize the provider.

Here is an example of how to use this provider:

```hcl
provider "qdrant" {
  alias = "qdrant_cloud"
  api_key = "<your_api_key>"
}
```

Replace `<your_api_key>` with your actual API key and URL.

## Developing The Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.21+ is required). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```bash
$ make build
...
$ $GOPATH/bin/terraform-provider-qdrant-cloud
...
```

### Updating the generated Go Client to interact with the Qdrant public API
This assumes that you cloned the qdrant-cloud-cluster-api to the same base-path as this repo (terraform-provider-qdrant-cloud)

```bash
make update-go-client
```

## Testing

In order to test the provider, you can run `make test`.

```bash
$ make test
```

This will run the unit tests in the provider.

## Contributing

If you find any issues or would like to contribute, feel free to create an issue or a pull request.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
