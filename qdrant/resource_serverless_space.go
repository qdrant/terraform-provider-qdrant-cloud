package qdrant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	spacev1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/v1"
)

// resourceServerlessSpace constructs a Terraform resource for managing a serverless space.
func resourceServerlessSpace() *schema.Resource {
	return &schema.Resource{
		Description:   "Serverless Space Resource",
		ReadContext:   resourceServerlessSpaceRead,
		CreateContext: resourceServerlessSpaceCreate,
		UpdateContext: resourceServerlessSpaceUpdate,
		DeleteContext: resourceServerlessSpaceDelete,
		Schema:        serverlessSpaceSchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(serverlessDefaultTimeout),
			Update: schema.DefaultTimeout(serverlessDefaultTimeout),
			Delete: schema.DefaultTimeout(serverlessDefaultTimeout),
		},
	}
}

func resourceServerlessSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading serverless space"
	client, clientCtx, diags := getServiceClient(ctx, m, spacev1.NewSpaceServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	resp, err := client.GetSpace(clientCtx, &spacev1.GetSpaceRequest{
		AccountId: accountUUID.String(),
		SpaceId:   d.Id(),
	}, grpc.Trailer(&trailer))
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	space := resp.GetSpace()
	if space.GetDeletedAt() != nil {
		d.SetId("")
		return nil
	}
	if err := setFlattened(d, flattenServerlessSpace(space)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessSpaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating serverless space"
	client, clientCtx, diags := getServiceClient(ctx, m, spacev1.NewSpaceServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	space := expandServerlessSpace(d, accountUUID.String())

	var trailer metadata.MD
	backupID := d.Get(slSpaceFromBackupIDFieldName).(string)
	if backupID == "" {
		resp, err := client.CreateSpace(clientCtx, &spacev1.CreateSpaceRequest{Space: space}, grpc.Trailer(&trailer))
		if err != nil {
			return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
		}
		d.SetId(resp.GetSpace().GetId())
	} else {
		resp, err := client.CreateSpaceFromBackup(clientCtx, &spacev1.CreateSpaceFromBackupRequest{
			AccountId: accountUUID.String(),
			BackupId:  backupID,
			SpaceName: space.GetName(),
		}, grpc.Trailer(&trailer))
		if err != nil {
			return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
		}
		d.SetId(resp.GetSpace().GetId())
	}

	if err := waitForServerlessSpaceReady(ctx, client, clientCtx, accountUUID.String(), d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	if backupID != "" {
		// CreateSpaceFromBackup only accepts a name, apply the remaining configured fields with an update.
		paths := serverlessSpaceUpdateMask(func(attr string) bool {
			if attr == slSpaceNameFieldName {
				return false
			}
			_, ok := d.GetOk(attr)
			return ok
		})
		if len(paths) > 0 {
			if diags := updateServerlessSpace(ctx, d, client, clientCtx, space, paths, d.Timeout(schema.TimeoutCreate), errorPrefix); diags.HasError() {
				return diags
			}
		}
	}
	return resourceServerlessSpaceRead(ctx, d, m)
}

func resourceServerlessSpaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error updating serverless space"
	client, clientCtx, diags := getServiceClient(ctx, m, spacev1.NewSpaceServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	paths := serverlessSpaceUpdateMask(d.HasChange)
	if len(paths) > 0 {
		space := expandServerlessSpace(d, accountUUID.String())
		if diags := updateServerlessSpace(ctx, d, client, clientCtx, space, paths, d.Timeout(schema.TimeoutUpdate), errorPrefix); diags.HasError() {
			return diags
		}
	}
	return resourceServerlessSpaceRead(ctx, d, m)
}

// updateServerlessSpace updates the given paths of the space (with the ID of the resource) and waits until it is ready.
func updateServerlessSpace(
	ctx context.Context,
	d *schema.ResourceData,
	client spacev1.SpaceServiceClient,
	clientCtx context.Context,
	space *spacev1.Space,
	paths []string,
	timeout time.Duration,
	errorPrefix string,
) diag.Diagnostics {
	space.Id = d.Id()
	var trailer metadata.MD
	_, err := client.UpdateSpace(clientCtx, &spacev1.UpdateSpaceRequest{
		Space:      space,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: paths},
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := waitForServerlessSpaceReady(ctx, client, clientCtx, space.GetAccountId(), d.Id(), timeout); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessSpaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting serverless space"
	client, clientCtx, diags := getServiceClient(ctx, m, spacev1.NewSpaceServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	_, err = client.DeleteSpace(clientCtx, &spacev1.DeleteSpaceRequest{
		AccountId:     accountUUID.String(),
		SpaceId:       d.Id(),
		DeleteBackups: newPointer(d.Get(slSpaceDeleteBackupsFieldName).(bool)),
	}, grpc.Trailer(&trailer))
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	stateConf := &retry.StateChangeConf{
		Pending:      []string{serverlessWaitPending},
		Target:       []string{serverlessWaitDeleted},
		Refresh:      serverlessSpaceDeletedRefreshFunc(client, clientCtx, accountUUID.String(), d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: serverlessPollInterval,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId("")
	return nil
}

// waitForServerlessSpaceReady waits until the space is ready and has an endpoint.
func waitForServerlessSpaceReady(
	ctx context.Context,
	client spacev1.SpaceServiceClient,
	clientCtx context.Context,
	accountID, spaceID string,
	timeout time.Duration,
) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{serverlessWaitPending},
		Target:       []string{serverlessWaitReady},
		Refresh:      serverlessSpaceReadyRefreshFunc(client, clientCtx, accountID, spaceID),
		Timeout:      timeout,
		PollInterval: serverlessPollInterval,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// serverlessSpaceReadyRefreshFunc returns a StateRefreshFunc that polls GetSpace
// until the space is ready and its endpoint is populated, or the space is disabled.
func serverlessSpaceReadyRefreshFunc(
	client spacev1.SpaceServiceClient,
	ctx context.Context,
	accountID, spaceID string,
) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetSpace(ctx, &spacev1.GetSpaceRequest{
			AccountId: accountID,
			SpaceId:   spaceID,
		})
		if err != nil {
			return nil, "", err
		}
		space := resp.GetSpace()
		state := space.GetState()
		switch state.GetPhase() {
		case spacev1.SpaceStatePhase_SPACE_STATE_PHASE_DISABLED, spacev1.SpaceStatePhase_SPACE_STATE_PHASE_DELETING:
			return nil, "", fmt.Errorf("space is not usable (phase=%q reason=%q)", state.GetPhase().String(), state.GetReason())
		case spacev1.SpaceStatePhase_SPACE_STATE_PHASE_READY:
			if strings.TrimSpace(state.GetEndpoint().GetUrl()) != "" {
				return space, serverlessWaitReady, nil
			}
		}
		return space, serverlessWaitPending, nil
	}
}

// serverlessSpaceDeletedRefreshFunc returns a StateRefreshFunc that polls GetSpace until the space is gone.
func serverlessSpaceDeletedRefreshFunc(
	client spacev1.SpaceServiceClient,
	ctx context.Context,
	accountID, spaceID string,
) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetSpace(ctx, &spacev1.GetSpaceRequest{
			AccountId: accountID,
			SpaceId:   spaceID,
		})
		if err != nil {
			if isNotFound(err) {
				return struct{}{}, serverlessWaitDeleted, nil
			}
			return nil, "", err
		}
		if resp.GetSpace().GetDeletedAt() != nil {
			return resp.GetSpace(), serverlessWaitDeleted, nil
		}
		return resp.GetSpace(), serverlessWaitPending, nil
	}
}
