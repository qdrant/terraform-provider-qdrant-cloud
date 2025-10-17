package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

// resourceAccountsRole constructs a Terraform resource for managing a role associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the CRUD functions.
func resourceAccountsRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Role resource for Qdrant Cloud (custom roles).",
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Schema:        accountsRoleSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// resourceRoleCreate performs a create operation to generate a new IAM role.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error creating role"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get a client
	client := qci.NewIAMServiceClient(apiClientConn)

	// Expand the role from Terraform configuration (no ID required for creation)
	role, err := expandRole(d, getDefaultAccountID(m), false)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	// Create the role
	var trailer metadata.MD
	resp, err := client.CreateRole(
		clientCtx,
		&qci.CreateRoleRequest{Role: role},
		grpc.Trailer(&trailer),
	)
	// Enrich prefix with request ID
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	// Inspect the results
	created := resp.GetRole()
	// Set the ID
	d.SetId(created.GetId())

	// Flatten role and store in Terraform state
	for k, v := range flattenRole(created) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", op, err))
		}
	}

	return nil
}

// resourceRoleRead performs a read operation to fetch a specific IAM role.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error reading role"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get a client
	client := qci.NewIAMServiceClient(apiClientConn)

	// Get the account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	// Fetch the role
	var trailer metadata.MD
	resp, err := client.GetRole(
		clientCtx,
		&qci.GetRoleRequest{
			AccountId: accountUUID.String(),
			RoleId:    d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	// Enrich prefix with request ID
	if err != nil {
		// Resource gone in the backend, clear state
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	// Inspect the results
	got := resp.GetRole()
	// Set the ID
	d.SetId(got.GetId())

	// Flatten role and store in Terraform state
	for k, v := range flattenRole(got) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", op, err))
		}
	}

	return nil
}

// resourceRoleUpdate performs an update operation on an existing IAM role.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error updating role"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get a client
	client := qci.NewIAMServiceClient(apiClientConn)

	// Get the account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	// Expand the role from Terraform configuration (include ID for update)
	role, err := expandRole(d, accountUUID.String(), true)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	// Update the role
	var trailer metadata.MD
	resp, err := client.UpdateRole(
		clientCtx,
		&qci.UpdateRoleRequest{Role: role},
		grpc.Trailer(&trailer),
	)
	// Enrich prefix with request ID
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	// Inspect the results
	updated := resp.GetRole()
	// Set the ID
	d.SetId(updated.GetId())

	// Flatten role and store in Terraform state
	for k, v := range flattenRole(updated) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", op, err))
		}
	}

	return nil
}

// resourceRoleDelete performs a delete operation to remove an IAM role.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error deleting role"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get a client
	client := qci.NewIAMServiceClient(apiClientConn)

	// Get the account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	// Delete the role
	var trailer metadata.MD
	_, err = client.DeleteRole(
		clientCtx,
		&qci.DeleteRoleRequest{
			AccountId: accountUUID.String(),
			RoleId:    d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	// Enrich prefix with request ID
	if err != nil {
		// Resource gone in the backend, clear state
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	// Resource gone in the backend, clear state
	d.SetId("")
	return nil
}
