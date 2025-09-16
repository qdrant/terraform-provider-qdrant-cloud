package qdrant

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	authv2 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/auth/v2"
)

// resourceAccountsAuthKeyV2 constructs a Terraform resource for managing a v2 API key.
// Returns a schema.Resource pointer configured with schema definitions and the CRUD functions.
func resourceAccountsAuthKeyV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Account Database API Key Resource (v2)",
		ReadContext:   resourceAPIKeyV2Read,
		CreateContext: resourceAPIKeyV2Create,
		UpdateContext: nil, // Not available in the public API
		DeleteContext: resourceAPIKeyV2Delete,
		Schema:        accountsAuthKeyV2ResourceSchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected <cluster_id>/<api_key_id>", d.Id())
				}
				d.Set("cluster_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

// resourceAPIKeyV2Read performs a read operation to fetch a specific v2 API key.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error getting API key (v2)"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := authv2.NewDatabaseApiKeyServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	clusterID := d.Get(authKeysV2ClusterIDFieldName).(string)
	apiKeyID := d.Id()

	var trailer metadata.MD
	resp, err := client.ListDatabaseApiKeys(clientCtx, &authv2.ListDatabaseApiKeysRequest{
		AccountId: accountUUID.String(),
		ClusterId: clusterID,
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Process the correct one, if any
	for _, apiKey := range resp.GetItems() {
		if apiKey.GetId() == apiKeyID {
			for k, v := range flattenAuthKeyV2(apiKey, false) {
				if err := d.Set(k, v); err != nil {
					return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
				}
			}
			d.SetId(apiKeyID)
			return nil
		}
	}
	// If the key is not found, it might have been deleted manually.
	// Remove it from the state.
	d.SetId("")
	return nil
}

// resourceAPIKeyV2Create performs a create operation to generate a new v2 API key.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceAPIKeyV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating API Key (v2)"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := authv2.NewDatabaseApiKeyServiceClient(apiClientConn)

	apiKey, err := expandAuthKeyV2(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	var trailer metadata.MD
	resp, err := client.CreateDatabaseApiKey(clientCtx, &authv2.CreateDatabaseApiKeyRequest{
		DatabaseApiKey: apiKey,
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	for k, v := range flattenAuthKeyV2(resp.GetDatabaseApiKey(), true) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	d.SetId(resp.GetDatabaseApiKey().GetId())
	return nil
}

// resourceAPIKeyV2Delete performs a delete operation to remove a v2 API key.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyV2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting API Key (v2)"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := authv2.NewDatabaseApiKeyServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	clusterID := d.Get(authKeysV2ClusterIDFieldName).(string)
	apiKeyID := d.Id()

	var trailer metadata.MD
	_, err = client.DeleteDatabaseApiKey(clientCtx, &authv2.DeleteDatabaseApiKeyRequest{
		AccountId:        accountUUID.String(),
		ClusterId:        clusterID,
		DatabaseApiKeyId: apiKeyID,
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	d.SetId("")
	return nil
}

// expandAuthKeyV2 expands the Terraform resource data into a v2 API key object.
func expandAuthKeyV2(d *schema.ResourceData, accountID string) (*authv2.DatabaseApiKey, error) {
	if v, ok := d.GetOk(authKeysV2AccountIDFieldName); ok {
		accountID = v.(string)
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID not specified")
	}

	key := &authv2.DatabaseApiKey{
		AccountId: accountID,
		ClusterId: d.Get(authKeysV2ClusterIDFieldName).(string),
		Name:      d.Get(authKeysV2NameFieldName).(string),
	}

	if v, ok := d.GetOk(authKeysV2ExpiresAtFieldName); ok {
		key.ExpiresAt = parseTime(v.(string))
	}

	var rules []*authv2.AccessRule
	if v, ok := d.GetOk(authKeysV2GlobalAccessRuleFieldName); ok {
		if len(v.([]interface{})) > 0 {
			ruleMap := v.([]interface{})[0].(map[string]interface{})
			accessTypeStr := ruleMap[authKeysV2AccessTypeFieldName].(string)
			accessTypeValue, ok := authv2.GlobalAccessRuleAccessType_value[accessTypeStr]
			if !ok {
				return nil, fmt.Errorf("invalid global access type: %s", accessTypeStr)
			}
			accessType := authv2.GlobalAccessRuleAccessType(accessTypeValue)
			rules = append(rules, &authv2.AccessRule{
				Scope: &authv2.AccessRule_GlobalAccess{
					GlobalAccess: &authv2.GlobalAccessRule{AccessType: accessType},
				},
			})
		}
	}
	if v, ok := d.GetOk(authKeysV2CollectionAccessRulesFieldName); ok {
		for _, rawRule := range v.([]interface{}) {
			ruleMap := rawRule.(map[string]interface{})
			accessTypeStr := ruleMap[authKeysV2AccessTypeFieldName].(string)
			accessTypeValue, ok := authv2.CollectionAccessRuleAccessType_value[accessTypeStr]
			if !ok {
				return nil, fmt.Errorf("invalid collection access type: %s", accessTypeStr)
			}
			accessType := authv2.CollectionAccessRuleAccessType(accessTypeValue)
			collectionRule := &authv2.CollectionAccessRule{
				CollectionName: ruleMap[authKeysV2CollectionNameFieldName].(string),
				AccessType:     accessType,
			}
			rules = append(rules, &authv2.AccessRule{
				Scope: &authv2.AccessRule_CollectionAccess{CollectionAccess: collectionRule},
			})
		}
	}
	if len(rules) == 0 {
		// Insert the default
		rules = append(rules, &authv2.AccessRule{
			Scope: &authv2.AccessRule_GlobalAccess{
				GlobalAccess: &authv2.GlobalAccessRule{AccessType: authv2.GlobalAccessRuleAccessType_GLOBAL_ACCESS_RULE_ACCESS_TYPE_MANAGE},
			},
		})
	}
	key.AccessRules = rules

	return key, nil
}
