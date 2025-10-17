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

// resourceAccountsRole constructs a Terraform resource for managing roles.
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

// resourceRoleCreate creates an role and writes the returned fields into state.
func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error creating role"

	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	client := qci.NewIAMServiceClient(apiClientConn)

	role, err := expandRoleForCreate(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	var trailer metadata.MD
	resp, err := client.CreateRole(
		clientCtx,
		&qci.CreateRoleRequest{Role: role},
		grpc.Trailer(&trailer),
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	created := resp.GetRole()
	d.SetId(created.GetId())

	for k, v := range flattenRole(created) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", op, err))
		}
	}

	return nil
}

// resourceRoleRead reads an role by ID and refreshes Terraform state.
func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error reading role"

	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	client := qci.NewIAMServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	var trailer metadata.MD
	resp, err := client.GetRole(
		clientCtx,
		&qci.GetRoleRequest{
			AccountId: accountUUID.String(),
			RoleId:    d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	got := resp.GetRole()
	d.SetId(got.GetId())

	for k, v := range flattenRole(got) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", op, err))
		}
	}

	return nil
}

// resourceRoleUpdate updates an role and refreshes Terraform state.
func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error updating role"

	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	client := qci.NewIAMServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	role := &qci.Role{
		Id:          d.Id(),
		AccountId:   accountUUID.String(),
		Name:        d.Get(roleNameFieldName).(string),
		Description: d.Get(roleDescriptionFieldName).(string),
		RoleType:    qci.RoleType_ROLE_TYPE_CUSTOM,
		Permissions: expandPermissions(d.Get(rolePermissionsFieldName)),
	}

	var trailer metadata.MD
	resp, err := client.UpdateRole(
		clientCtx,
		&qci.UpdateRoleRequest{Role: role},
		grpc.Trailer(&trailer),
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	updated := resp.GetRole()
	d.SetId(updated.GetId())

	for k, v := range flattenRole(updated) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", op, err))
		}
	}

	return nil
}

// resourceRoleDelete deletes an role and clears its Terraform state.
func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error deleting role"

	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	client := qci.NewIAMServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}

	var trailer metadata.MD
	_, err = client.DeleteRole(
		clientCtx,
		&qci.DeleteRoleRequest{
			AccountId: accountUUID.String(),
			RoleId:    d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, getRequestID(trailer), err))
	}

	d.SetId("")
	return nil
}
