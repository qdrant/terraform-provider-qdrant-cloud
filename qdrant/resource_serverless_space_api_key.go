package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	spaceauthv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/auth/v1"
)

// resourceServerlessSpaceAPIKey constructs a Terraform resource for managing API keys of a serverless space.
// API keys cannot be updated, so all configurable fields force a new key.
func resourceServerlessSpaceAPIKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Serverless Space API Key Resource",
		ReadContext:   resourceServerlessSpaceAPIKeyRead,
		CreateContext: resourceServerlessSpaceAPIKeyCreate,
		UpdateContext: nil, // Not available in the public API
		DeleteContext: resourceServerlessSpaceAPIKeyDelete,
		Schema:        serverlessSpaceAPIKeySchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: importSpaceChildState(slSpaceAPIKeySpaceIDFieldName, "api_key_id"),
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(serverlessDefaultTimeout),
		},
	}
}

func resourceServerlessSpaceAPIKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading serverless space API key"
	client, clientCtx, diags := getServiceClient(ctx, m, spaceauthv1.NewSpaceApiKeyServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// There is no Get endpoint, so list the keys of the space and find ours.
	var trailer metadata.MD
	resp, err := client.ListSpaceApiKeys(clientCtx, &spaceauthv1.ListSpaceApiKeysRequest{
		AccountId: accountUUID.String(),
		SpaceId:   d.Get(slSpaceAPIKeySpaceIDFieldName).(string),
	}, grpc.Trailer(&trailer))
	if err != nil {
		if isNotFound(err) {
			// The space is gone, and so is the key.
			d.SetId("")
			return nil
		}
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	for _, key := range resp.GetItems() {
		if key.GetId() != d.Id() {
			continue
		}
		if key.GetDeletedAt() != nil || key.GetRevokedAt() != nil {
			break
		}
		if err := setFlattened(d, flattenServerlessSpaceAPIKey(key)); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
		return nil
	}
	d.SetId("")
	return nil
}

func resourceServerlessSpaceAPIKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating serverless space API key"
	client, clientCtx, diags := getServiceClient(ctx, m, spaceauthv1.NewSpaceApiKeyServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	key := expandServerlessSpaceAPIKey(d, accountUUID.String())
	var trailer metadata.MD
	resp, err := client.CreateSpaceApiKey(clientCtx, &spaceauthv1.CreateSpaceApiKeyRequest{
		SpaceApiKey: key,
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	created := resp.GetSpaceApiKey()
	d.SetId(created.GetId())
	// The secret is only returned on creation.
	if err := d.Set(slSpaceAPIKeyKeyFieldName, created.GetKey()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending:      []string{serverlessWaitPending},
		Target:       []string{serverlessWaitReady},
		Refresh:      serverlessSpaceAPIKeyReadyRefreshFunc(client, clientCtx, key.GetAccountId(), key.GetSpaceId(), created.GetId()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: serverlessPollInterval,
	}
	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := setFlattened(d, flattenServerlessSpaceAPIKey(result.(*spaceauthv1.SpaceApiKey))); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessSpaceAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting serverless space API key"
	client, clientCtx, diags := getServiceClient(ctx, m, spaceauthv1.NewSpaceApiKeyServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	_, err = client.DeleteSpaceApiKey(clientCtx, &spaceauthv1.DeleteSpaceApiKeyRequest{
		AccountId:     accountUUID.String(),
		SpaceId:       d.Get(slSpaceAPIKeySpaceIDFieldName).(string),
		SpaceApiKeyId: d.Id(),
	}, grpc.Trailer(&trailer))
	if err != nil && !isNotFound(err) {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	d.SetId("")
	return nil
}

// serverlessSpaceAPIKeyReadyRefreshFunc returns a StateRefreshFunc that polls the API keys of the space
// until the given key is ready, or fails if it is disabled or gone.
func serverlessSpaceAPIKeyReadyRefreshFunc(
	client spaceauthv1.SpaceApiKeyServiceClient,
	ctx context.Context,
	accountID, spaceID, keyID string,
) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ListSpaceApiKeys(ctx, &spaceauthv1.ListSpaceApiKeysRequest{
			AccountId: accountID,
			SpaceId:   spaceID,
		})
		if err != nil {
			return nil, "", err
		}
		for _, key := range resp.GetItems() {
			if key.GetId() != keyID {
				continue
			}
			state := key.GetState()
			switch state.GetPhase() {
			case spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_READY:
				return key, serverlessWaitReady, nil
			case spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_DISABLED,
				spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_DELETING:
				return nil, "", fmt.Errorf("API key is not usable (phase=%q reason=%q)", state.GetPhase().String(), state.GetReason())
			default:
				return key, serverlessWaitPending, nil
			}
		}
		return nil, "", fmt.Errorf("API key %s not found in space %s", keyID, spaceID)
	}
}
