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

	qch "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/hybrid/v1"
)

// resourceAccountsHybridCloudEnvironment constructs a Terraform resource for managing a hybrid cloud environment associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the CRUD functions.
func resourceAccountsHybridCloudEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "Hybrid Cloud Environment Resource",
		CreateContext: resourceHCEnvCreate,
		ReadContext:   resourceHCEnvRead,
		UpdateContext: resourceHCEnvUpdate,
		DeleteContext: resourceHCEnvDelete,
		Schema:        accountsHybridCloudEnvironmentSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			// If version changed, mark bootstrap_commands as changing to a computed value
			if d.HasChange(hcEnvBootstrapCommandsVersionFieldName) {
				if err := d.SetNewComputed(hcEnvBootstrapCommandsFieldName); err != nil {
					return err
				}
				return nil
			}

			// If version did NOT change, suppress any spurious diff on bootstrap_commands
			if d.HasChange(hcEnvBootstrapCommandsFieldName) {
				old, _ := d.GetChange(hcEnvBootstrapCommandsFieldName)
				if err := d.SetNew(hcEnvBootstrapCommandsFieldName, old); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

// resourceHCEnvCreate performs a create operation to generate a new hybrid cloud environment.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceHCEnvCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating hybrid cloud environment"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get a client
	client := qch.NewHybridCloudServiceClient(apiClientConn)

	// Expand the hybrid cloud environment
	env, err := expandHCEnv(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	// Create the hybrid cloud environment
	var trailer metadata.MD
	resp, err := client.CreateHybridCloudEnvironment(
		clientCtx,
		&qch.CreateHybridCloudEnvironmentRequest{HybridCloudEnvironment: env},
		grpc.Trailer(&trailer),
	)
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	created := resp.GetHybridCloudEnvironment()
	// Set the ID
	d.SetId(created.GetId())

	for k, v := range flattenHCEnv(created) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	v := d.Get(hcEnvBootstrapCommandsVersionFieldName).(int)

	switch v {
	case -1:
		// Explicit opt-out: do NOT generate. Leave -1 in state.
	case 0:
		// Default 0 on create -> bump to 1
		_ = d.Set(hcEnvBootstrapCommandsVersionFieldName, 1)
	default:
		// v > 0: keep as-is
	}

	// Generate only when version > 0 (respect explicit -1 and 0)
	if ver, _ := d.Get(hcEnvBootstrapCommandsVersionFieldName).(int); ver > 0 {
		if ds := setHCEnvBootstrapCommands(client, clientCtx, d, m, "error getting bootstrap commands"); ds.HasError() {
			return ds
		}
	}

	return nil
}

// resourceHCEnvRead performs a read operation to fetch a specific hybrid cloud environment.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceHCEnvRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading hybrid cloud environment"

	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get a client
	client := qch.NewHybridCloudServiceClient(apiClientConn)

	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	// Fetch the hybrid cloud environment
	var trailer metadata.MD
	resp, err := client.GetHybridCloudEnvironment(
		clientCtx,
		&qch.GetHybridCloudEnvironmentRequest{
			AccountId:                accountUUID.String(),
			HybridCloudEnvironmentId: d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	errorPrefix += getRequestID(trailer)

	// Inspect the results
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			// Resource gone in the backend, clear state
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	env := resp.GetHybridCloudEnvironment()
	// Set the ID
	d.SetId(env.GetId())

	// Flatten hybrid cloud environment and store in Terraform state
	for k, v := range flattenHCEnv(env) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	// Ensure local-only fields exist in state so import/refresh doesn't show spurious diffs.
	// We never generate in Read; we only seed zero values.
	if _, ok := d.GetOk(hcEnvBootstrapCommandsVersionFieldName); !ok {
		_ = d.Set(hcEnvBootstrapCommandsVersionFieldName, 0)
	}
	if _, ok := d.GetOk(hcEnvBootstrapCommandsFieldName); !ok {
		_ = d.Set(hcEnvBootstrapCommandsFieldName, []interface{}{})
	}

	return nil
}

// resourceHCEnvUpdate performs an update operation on a hybrid cloud environment.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceHCEnvUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	changedConfigOrName := d.HasChange(hcEnvNameFieldName) || d.HasChange(hcEnvConfigurationFieldName)
	changedVersion := d.HasChange(hcEnvBootstrapCommandsVersionFieldName)

	// Nothing changed: just refresh (doesn't touch bootstrap fields)
	if !changedConfigOrName && !changedVersion {
		return resourceHCEnvRead(ctx, d, m)
	}

	// 1) Apply config/name changes (if any)
	if changedConfigOrName {
		errorPrefix := "error updating hybrid cloud environment"
		apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
		if diags.HasError() {
			return diags
		}
		client := qch.NewHybridCloudServiceClient(apiClientConn)

		env, err := expandHCEnv(d, getDefaultAccountID(m))
		if err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}

		var trailer metadata.MD
		_, err = client.UpdateHybridCloudEnvironment(
			clientCtx,
			&qch.UpdateHybridCloudEnvironmentRequest{HybridCloudEnvironment: env},
			grpc.Trailer(&trailer),
		)
		errorPrefix += getRequestID(trailer)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	// 2) If version changed, rotate/clear explicitly — this is the ONLY place we generate.
	if changedVersion {
		apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
		if diags.HasError() {
			return diags
		}
		client := qch.NewHybridCloudServiceClient(apiClientConn)

		newV := d.Get(hcEnvBootstrapCommandsVersionFieldName).(int)
		switch {
		case newV > 0:
			// Regenerate
			if ds := setHCEnvBootstrapCommands(client, clientCtx, d, m, "error getting bootstrap commands"); ds.HasError() {
				return ds
			}
		default: // newV == 0 or -1
			// Clear in state
			_ = d.Set(hcEnvBootstrapCommandsFieldName, []interface{}{})
		}
	}

	// 3) Final refresh of server-side fields (won’t mutate bootstrap fields)
	return resourceHCEnvRead(ctx, d, m)
}

// resourceHCEnvDelete performs a delete operation to remove a hybrid cloud environment.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceHCEnvDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting hybrid cloud environment"
	// Get a client connection and context
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	// Get a client
	client := qch.NewHybridCloudServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Delete the hybrid cloud environment
	var trailer metadata.MD
	_, err = client.DeleteHybridCloudEnvironment(
		clientCtx,
		&qch.DeleteHybridCloudEnvironmentRequest{
			AccountId:                accountUUID.String(),
			HybridCloudEnvironmentId: d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	// Inspect the results
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Resource gone in the backend, clear state
	d.SetId("")
	return nil
}

// setHCEnvBootstrapCommands performs a read operation to fetch the bootstrap commands for a hybrid cloud environment.
// client: The gRPC client for the hybrid cloud service.
// clientCtx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// errorPrefix: A string to prefix error messages with.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func setHCEnvBootstrapCommands(
	client qch.HybridCloudServiceClient,
	clientCtx context.Context,
	d *schema.ResourceData,
	m interface{},
	errorPrefix string,
) diag.Diagnostics {
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Fetch the bootstrap commands
	var trailer metadata.MD
	resp, err := client.GenerateBootstrapCommands(
		clientCtx,
		&qch.GenerateBootstrapCommandsRequest{
			AccountId:                accountUUID.String(),
			HybridCloudEnvironmentId: d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	if err != nil {
		// Soft-fail on common "not ready yet" or permission issues
		if st, ok := status.FromError(err); ok && (st.Code() == codes.FailedPrecondition || st.Code() == codes.PermissionDenied) {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "Bootstrap commands not available",
				Detail:   "The environment may not be ready yet or your credentials lack permission. Re-run plan/apply later to refresh.",
			}}
		}
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Flatten bootstrap commands and store in Terraform state
	cmds := resp.GetCommands()
	values := make([]interface{}, len(cmds))
	for i, c := range cmds {
		values[i] = c
	}
	if err := d.Set(hcEnvBootstrapCommandsFieldName, values); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}
