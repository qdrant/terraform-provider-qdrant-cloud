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

	qca "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/account/v1"
	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

// resourceAccountsUserRoles constructs a Terraform resource for managing role assignments
// for a user (specified by email) within an account. The resource is non-authoritative:
// it only adds the roles specified and never removes roles that are not managed by this resource.
// On delete, it revokes only the roles managed by this resource unless keep_on_destroy is true.
func resourceAccountsUserRoles() *schema.Resource {
	return &schema.Resource{
		Description:   "User Roles resource for Qdrant Cloud (non-authoritative add-only per user; user specified by email).",
		CreateContext: resourceUserRolesCreate,
		ReadContext:   resourceUserRolesRead,
		UpdateContext: resourceUserRolesUpdate,
		DeleteContext: resourceUserRolesDelete,
		Schema:        accountsUserRolesSchema(),
	}
}

// resourceUserRolesCreate performs a create operation that assigns only the requested roles,
// without removing any existing roles that are not part of this resource.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceUserRolesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error creating user roles"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get clients
	iamClient := qci.NewIAMServiceClient(apiClientConn)
	acctClient := qca.NewAccountServiceClient(apiClientConn)

	// Determine account ID
	accountID := getDefaultAccountID(m)
	if v, ok := d.GetOk(userRolesAccountIdFieldName); ok && v.(string) != "" {
		accountID = v.(string)
	}
	if accountID == "" {
		return diag.FromErr(fmt.Errorf("%s: account ID not specified", op))
	}
	// persist computed account_id
	_ = d.Set(userRolesAccountIdFieldName, accountID)

	// Resolve user ID
	email := d.Get(userRolesUserEmailFieldName).(string)
	userID, reqID, err := resolveUserIDByEmail(clientCtx, acctClient, accountID, email)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
	}
	_ = d.Set(userRolesUserIdFieldName, userID)

	// Get current user roles
	current, reqID, err := listUserRoleIDs(clientCtx, iamClient, accountID, userID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
	}

	// Compute additions only
	desired := setToStringSlice(d.Get(userRolesRoleIdsFieldName))
	toAdd, _ := diffStringSets(desired, current)
	if len(toAdd) > 0 {
		if reqID, err = assignUserRoles(clientCtx, iamClient, accountID, userID, toAdd, nil); err != nil {
			return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
		}
	}

	// Set canonical ID
	d.SetId(fmt.Sprintf("%s/%s", accountID, userID))
	return resourceUserRolesRead(ctx, d, m)
}

// resourceUserRolesRead fetches the user identifier and ensures the canonical ID is set.
// It intentionally does not write role_ids to avoid removing externally managed roles.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceUserRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error reading user roles"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get client
	acctClient := qca.NewAccountServiceClient(apiClientConn)

	// Get account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}
	accountID := accountUUID.String()
	// persist computed account_id
	_ = d.Set(userRolesAccountIdFieldName, accountID)

	// Resolve user ID (prefer state; fallback to email)
	userID := d.Get(userRolesUserIdFieldName).(string)
	if userID == "" {
		email := d.Get(userRolesUserEmailFieldName).(string)
		userID, _, err = resolveUserIDByEmail(clientCtx, acctClient, accountID, email)
		if err != nil {
			d.SetId("")
			return nil
		}
		_ = d.Set(userRolesUserIdFieldName, userID)
	}

	// Clear state if unresolved
	if userID == "" {
		d.SetId("")
		return nil
	}

	// Set canonical ID
	d.SetId(fmt.Sprintf("%s/%s", accountID, userID))
	return nil
}

// resourceUserRolesUpdate assigns newly requested roles and revokes roles that
// were previously managed by this resource but are no longer requested.
// It never removes roles that were not managed by this resource.
//
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceUserRolesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error updating user roles"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get clients
	iamClient := qci.NewIAMServiceClient(apiClientConn)
	acctClient := qca.NewAccountServiceClient(apiClientConn)

	// Get account ID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}
	accountID := accountUUID.String()
	// persist computed account_id
	_ = d.Set(userRolesAccountIdFieldName, accountID)

	// Re-resolve if email changed or missing user_id
	email := d.Get(userRolesUserEmailFieldName).(string)
	if d.HasChange(userRolesUserEmailFieldName) || d.Get(userRolesUserIdFieldName).(string) == "" {
		userID, reqID, err := resolveUserIDByEmail(clientCtx, acctClient, accountID, email)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
		}
		_ = d.Set(userRolesUserIdFieldName, userID)
		d.SetId(fmt.Sprintf("%s/%s", accountID, userID))
	}

	userID := d.Get(userRolesUserIdFieldName).(string)
	if userID == "" {
		return diag.FromErr(fmt.Errorf("%s: resolved user_id is empty", op))
	}

	// Current roles on the user (server view)
	current, reqID, err := listUserRoleIDs(clientCtx, iamClient, accountID, userID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
	}

	// Desired roles (new config)
	newDesired := setToStringSlice(d.Get(userRolesRoleIdsFieldName))

	// Additions: (newDesired - current)
	toAdd, _ := diffStringSets(newDesired, current)

	// Deletions: only roles this resource *previously* managed and has now removed.
	// delCandidates = (oldDesired - newDesired) -> then limit to currently assigned.
	var toDelete []string
	if d.HasChange(userRolesRoleIdsFieldName) {
		oldRaw, newRaw := d.GetChange(userRolesRoleIdsFieldName)
		oldDesired := setToStringSlice(oldRaw)
		_ = newRaw // newDesired already computed from d.Get above

		delCandidates, _ := diffStringSets(oldDesired, newDesired) // old - new
		toDelete = intersectStrings(current, delCandidates)        // only if still assigned
	}

	// Single RPC with both add + delete if needed
	if len(toAdd) > 0 || len(toDelete) > 0 {
		if reqID, err = assignUserRoles(clientCtx, iamClient, accountID, userID, toAdd, toDelete); err != nil {
			return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
		}
	}

	// Finalize by refreshing computed fields
	return resourceUserRolesRead(ctx, d, m)
}

// resourceUserRolesDelete revokes only the roles managed by this resource (intersection of
// desired and current) unless keep_on_destroy is true.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceUserRolesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	const op = "error deleting user roles"

	// Skip revoke if keep_on_destroy is true
	if keep, _ := d.Get(userRolesKeepOnDestroyName).(bool); keep {
		d.SetId("")
		return nil
	}

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get clients
	iamClient := qci.NewIAMServiceClient(apiClientConn)
	acctClient := qca.NewAccountServiceClient(apiClientConn)

	// Get account ID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", op, err))
	}
	accountID := accountUUID.String()

	// Resolve user_id (prefer state; fallback to email)
	userID := d.Get(userRolesUserIdFieldName).(string)
	if userID == "" {
		email := d.Get(userRolesUserEmailFieldName).(string)
		userID, _, err = resolveUserIDByEmail(clientCtx, acctClient, accountID, email)
		if err != nil {
			d.SetId("")
			return nil
		}
	}

	// Get current roles
	current, reqID, err := listUserRoleIDs(clientCtx, iamClient, accountID, userID)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
	}

	// Compute intersection and revoke only those roles
	desired := setToStringSlice(d.Get(userRolesRoleIdsFieldName))
	toDelete := intersectStrings(current, desired)
	if len(toDelete) > 0 {
		if reqID, err = assignUserRoles(clientCtx, iamClient, accountID, userID, nil, toDelete); err != nil {
			return diag.FromErr(fmt.Errorf("%s%s: %w", op, reqID, err))
		}
	}

	// Clear state
	d.SetId("")
	return nil
}

// resolveUserIDByEmail resolves a user ID by exact email match using AccountService.ListAccountMembers.
// Returns: userID, requestID (trace), error.
func resolveUserIDByEmail(
	ctx context.Context,
	acctClient qca.AccountServiceClient,
	accountID string,
	email string,
) (string, string, error) {
	var trailer metadata.MD
	resp, err := acctClient.ListAccountMembers(
		ctx,
		&qca.ListAccountMembersRequest{AccountId: accountID},
		grpc.Trailer(&trailer),
	)
	reqID := getRequestID(trailer)
	if err != nil {
		return "", reqID, err
	}
	var matches []string
	for _, it := range resp.GetItems() {
		u := it.GetAccountMember()
		if u.GetEmail() == email {
			matches = append(matches, u.GetId())
		}
	}
	switch len(matches) {
	case 0:
		return "", reqID, fmt.Errorf("user with email %q not found in account %s", email, accountID)
	case 1:
		return matches[0], reqID, nil
	default:
		return "", reqID, fmt.Errorf("multiple users with email %q in account %s", email, accountID)
	}
}

// listUserRoleIDs returns the list of role IDs assigned to a user via IAMService.ListUserRoles.
// Returns: roleIDs, requestID (trace), error.
func listUserRoleIDs(
	ctx context.Context,
	iamClient qci.IAMServiceClient,
	accountID string,
	userID string,
) ([]string, string, error) {
	var trailer metadata.MD
	resp, err := iamClient.ListUserRoles(
		ctx,
		&qci.ListUserRolesRequest{
			AccountId: accountID,
			UserId:    userID,
		},
		grpc.Trailer(&trailer),
	)
	reqID := getRequestID(trailer)
	if err != nil {
		return nil, reqID, err
	}
	out := make([]string, 0, len(resp.GetRoles()))
	for _, r := range resp.GetRoles() {
		out = append(out, r.GetId())
	}
	return out, reqID, nil
}

// assignUserRoles wraps IAMService.AssignUserRoles.
// Returns: requestID (trace), error.
func assignUserRoles(
	ctx context.Context,
	iamClient qci.IAMServiceClient,
	accountID string,
	userID string,
	toAdd []string,
	toDel []string,
) (string, error) {
	var trailer metadata.MD
	_, err := iamClient.AssignUserRoles(
		ctx,
		&qci.AssignUserRolesRequest{
			AccountId:       accountID,
			UserId:          userID,
			RoleIdsToAdd:    toAdd,
			RoleIdsToDelete: toDel,
		},
		grpc.Trailer(&trailer),
	)
	reqID := getRequestID(trailer)
	if err != nil {
		return reqID, err
	}
	return reqID, nil
}
