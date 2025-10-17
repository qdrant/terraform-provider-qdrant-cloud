# Example: Role

This example shows how to use the Terraform Qdrant Cloud provider to manage **custom Role** resources in Qdrant Cloud.

## Prerequisites

This example uses syntax elements specific to a Terraform provider version; see the `terraform` block in the example `.tf` file for details.

## Environment variables

Please refer to the [Main README](../../README.md) for all environment variables you might need (e.g., API key, account ID, endpoint).

## How to run

```bash
terraform init
terraform plan
terraform apply
```

To remove the resources created, run:

```bash
terraform destroy
```

## Notes

- Each role requires at least one `permissions` block (e.g., `read:backups`, `write:backups` in the `Cluster` category).  
- Permissions are stored as a set; order does not affect plans.
