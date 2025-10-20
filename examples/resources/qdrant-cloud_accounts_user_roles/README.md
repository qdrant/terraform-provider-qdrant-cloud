# Example: User Role Assignments

This example shows how to use the Terraform Qdrant Cloud provider to manage **role assignments for users** within an account.

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

- This resource manages **role assignments for a single user (by email)**.  
- It is **non-authoritative** — existing roles outside of this resource are not removed.  
- On destroy, only the roles listed in `role_ids` are revoked (unless `keep_on_destroy = true`).  
- The user must already exist in the account.
