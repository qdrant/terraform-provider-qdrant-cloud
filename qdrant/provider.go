package qdrant

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

// Provider defines and returns a Terraform resource provider for Qdrant Cloud.
// It sets up the provider schema, resources, and data sources.
// Returns a pointer to the schema.Provider object.
func Provider() *schema.Provider {
	return &schema.Provider{
		// Schema defines the provider's configuration options.
		Schema: map[string]*schema.Schema{
			"auth": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_AUTH", authModeAPIKey),
				Description: "Authentication mode: api_key (default, CI/automation) or cli (local: runs `qcloud auth token --json` after `qcloud auth login`). Env: QDRANT_CLOUD_AUTH.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{authModeAPIKey, authModeCLI}, false)),
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_API_KEY", nil),
				Description: "The API Key for Qdrant Cloud API operations. Required when auth is api_key (default). Ignored when auth is cli.",
				Sensitive:   true,
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_API_URL", "grpc.cloud.qdrant.io"),
				Description: "The URL of the Qdrant Cloud API (gRPC endpoint host:port).",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_ACCOUNT_ID", ""),
				Description: "Default Account Identifier for the Qdrant cloud",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow insecure gRPC connections. This is useful for development environments with self-signed certificates. Defaults to false.",
			},
		},
		// ResourcesMap defines all the resources that this provider offers.
		ResourcesMap: map[string]*schema.Resource{
			"qdrant-cloud_accounts_auth_key":                 resourceAccountsAuthKey(),                // Resource for Qdrant Cloud accounts' authorization keys.
			"qdrant-cloud_accounts_database_api_key_v2":      resourceAccountsAuthKeyV2(),              // Resource for Qdrant Cloud accounts' authorization keys v2.
			"qdrant-cloud_accounts_cluster":                  resourceAccountsCluster(),                // Resource for managing Qdrant Cloud account clusters.
			"qdrant-cloud_accounts_backup_schedule":          resourceAccountsBackupSchedule(),         // Resource for managing Qdrant Cloud account backup schedules (for a cluster).
			"qdrant-cloud_accounts_manual_backup":            resourceAccountsManualBackup(),           // Resource for managing Qdrant Cloud account manual backup (for a cluster).
			"qdrant-cloud_accounts_hybrid_cloud_environment": resourceAccountsHybridCloudEnvironment(), // Resource for managing Qdrant Cloud account hybrid cloud environments.
			"qdrant-cloud_accounts_role":                     resourceAccountsRole(),                   // Resource for managing Qdrant Cloud account roles.
			"qdrant-cloud_accounts_user_roles":               resourceAccountsUserRoles(),              // Resource for managing role assignments for a user (by email) within an account.
			"qdrant-cloud_serverless_space":                  resourceServerlessSpace(),                // Resource for managing Qdrant Cloud serverless spaces.
			"qdrant-cloud_serverless_space_api_key":          resourceServerlessSpaceAPIKey(),          // Resource for managing API keys of a Qdrant Cloud serverless space.
			"qdrant-cloud_serverless_backup_schedule":        resourceServerlessBackupSchedule(),       // Resource for managing backup schedules of a Qdrant Cloud serverless space.
			"qdrant-cloud_serverless_backup":                 resourceServerlessBackup(),               // Resource for managing manual backups of a Qdrant Cloud serverless space.
		},
		// DataSourcesMap defines all the data sources that this provider offers.
		DataSourcesMap: map[string]*schema.Resource{
			"qdrant-cloud_accounts_auth_keys":            dataSourceAccountsAuthKeys(),          // Data source for retrieving Qdrant Cloud accounts' authorization keys.
			"qdrant-cloud_accounts_database_api_keys_v2": dataSourceAccountsAuthKeysV2(),        // Data source for retrieving Qdrant Cloud accounts' authorization keys v2.
			"qdrant-cloud_accounts_clusters":             dataSourceAccountsClusters(),          // Data source for listing Qdrant Cloud clusters under an account.
			"qdrant-cloud_accounts_cluster":              dataSourceAccountsCluster(),           // Data source for retrieving details of a specific Qdrant cluster.
			"qdrant-cloud_booking_packages":              dataSourceBookingPackages(),           // Data source for Qdrant booking packages.
			"qdrant-cloud_accounts_backup_schedules":     dataSourceAccountsBackupSchedules(),   // Data source for listing Qdrant Cloud backup schedules under an account and cluster.
			"qdrant-cloud_accounts_backup_schedule":      dataSourceAccountsBackupSchedule(),    // Data source for retrieving Qdrant Cloud accounts' backup schedules (for a cluster).
			"qdrant-cloud_accounts_members":              dataSourceAccountsMembers(),           // Data source for listing Qdrant Cloud account members.
			"qdrant-cloud_accounts_roles":                dataSourceAccountsRoles(),             // Data source for listing Qdrant Cloud account roles (system and custom).
			"qdrant-cloud_serverless_space":              dataSourceServerlessSpace(),           // Data source for retrieving a Qdrant Cloud serverless space.
			"qdrant-cloud_serverless_spaces":             dataSourceServerlessSpaces(),          // Data source for listing Qdrant Cloud serverless spaces.
			"qdrant-cloud_serverless_space_api_keys":     dataSourceServerlessSpaceAPIKeys(),    // Data source for listing the API keys of a Qdrant Cloud serverless space.
			"qdrant-cloud_serverless_backup_schedule":    dataSourceServerlessBackupSchedule(),  // Data source for retrieving a Qdrant Cloud serverless backup schedule.
			"qdrant-cloud_serverless_backup_schedules":   dataSourceServerlessBackupSchedules(), // Data source for listing Qdrant Cloud serverless backup schedules.
			"qdrant-cloud_serverless_backups":            dataSourceServerlessBackups(),         // Data source for listing Qdrant Cloud serverless backups.
			"qdrant-cloud_serverless_cloud_regions":      dataSourceServerlessCloudRegions(),    // Data source for listing cloud regions available for serverless spaces.
		},
		// ConfigureContextFunc points to the function used to configure the runtime environment of the provider.
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure initializes and configures a client using the provided schema resource data.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values.
// d: Resource data structure used to configure the client, typically provided by Terraform.
// Returns a configured client object and any diagnostic information.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	apiURL := d.Get("api_url").(string)
	authMode := strings.TrimSpace(d.Get("auth").(string))
	if authMode == "" {
		authMode = authModeAPIKey
	}
	var accountID string
	if aid, ok := d.GetOk("account_id"); ok {
		accountID = aid.(string)
	}
	insecure := d.Get("insecure").(bool)
	var diags diag.Diagnostics

	if strings.TrimSpace(apiURL) == "" {
		apiURL = "grpc.cloud.qdrant.io"
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Using default URL",
			Detail:   "No API URL was provided, using default URL " + apiURL,
		})
	}

	config := ProviderConfig{
		AuthMode:  authMode,
		ApiKey:    apiKey,
		BaseURL:   apiURL,
		AccountID: accountID,
		Insecure:  insecure,
	}

	switch authMode {
	case authModeCLI:
		config.CLITokens = newCLITokenSource(defaultQcloudBinary, apiURL)
	case authModeAPIKey:
		if strings.TrimSpace(apiKey) == "" {
			return nil, diag.Errorf("api_key must not be empty when auth is %q (set auth = %q after qcloud auth login for local use)", authModeAPIKey, authModeCLI)
		}
	default:
		return nil, diag.Errorf("unsupported auth %q (want %q or %q)", authMode, authModeAPIKey, authModeCLI)
	}

	return &config, diags
}

// ProviderConfig holds the configuration details for creating gRPC requests to the Qdrant Cloud API.
type ProviderConfig struct {
	AuthMode  string          // AuthMode is "api_key" (default) or "cli".
	ApiKey    string          // ApiKey is used when AuthMode is api_key.
	BaseURL   string          // BaseURL is the gRPC API endpoint host:port.
	AccountID string          // Default account ID, if any.
	Insecure  bool            // Insecure allows insecure TLS for development.
	CLITokens *cliTokenSource // CLITokens resolves Bearer tokens when AuthMode is cli.
}
