package qdrant

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"api_key": {
				Type:        schema.TypeString,                                  // Data type of the API key.
				Required:    true,                                               // API key is a required field.
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_API_KEY", nil), // Default can be set via an environment variable.
				Description: "The API Key for Qdrant Cloud API operations.",     // Description of the API key usage.
			},
			"api_url": {
				Type:        schema.TypeString,                                                        // Data type of the API URL.
				Optional:    true,                                                                     // API URL is an optional field, with a default provided.
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_API_URL", "https://cloud.qdrant.io"), // Default API URL.
				Description: "The URL of the Qdrant Cloud API.",                                       // Description of the API URL.
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_ACCOUNT_ID", ""),
				Description: "Default Account Identifier for the Qdrant cloud",
			},
		},
		// ResourcesMap defines all the resources that this provider offers.
		ResourcesMap: map[string]*schema.Resource{
			"qdrant-cloud_accounts_auth_key": resourceAccountsAuthKey(), // Resource for Qdrant Cloud accounts' authorization keys.
			"qdrant-cloud_accounts_cluster":  resourceAccountsCluster(), // Resource for managing Qdrant Cloud account clusters.
		},
		// DataSourcesMap defines all the data sources that this provider offers.
		DataSourcesMap: map[string]*schema.Resource{
			"qdrant-cloud_accounts_auth_keys": dataSourceAccountsAuthKeys(), // Data source for retrieving Qdrant Cloud accounts' authorization keys.
			"qdrant-cloud_accounts_clusters":  dataSourceAccountsClusters(), // Data source for listing Qdrant Cloud clusters under an account.
			"qdrant-cloud_accounts_cluster":   dataSourceAccountsCluster(),  // Data source for retrieving details of a specific Qdrant cluster.
			"qdrant-cloud_booking_packages":   dataSourceBookingPackages(),  // Data source for Qdrant booking packages.
		},
		// ConfigureContextFunc points to the function used to configure the runtime environment of the provider.
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure initializes and configures a client using the provided schema resource data.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values.
// d: Resource data structure used to configure the client, typically provided by Terraform.
// Returns a configured client object and any diagnostic information.
// If api_key or api_url is empty, it returns an error diagnostic.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Retrieve the API key and URL from the schema resource data.
	apiKey := d.Get("api_key").(string)
	apiURL := d.Get("api_url").(string)
	var accountID string
	if aid, ok := d.GetOk("account_id"); ok {
		accountID = aid.(string)
	}
	var diags diag.Diagnostics

	// Validate that the API key is not empty, returning an error diagnostic if it is.
	if strings.TrimSpace(apiKey) == "" {
		return nil, diag.Errorf("api_key must not be empty")
	}

	// Validate that the API URL is not empty, returning an error diagnostic if it is.
	if strings.TrimSpace(apiURL) == "" {
		apiURL = "https://cloud.qdrant.io/public/v1"
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Using default URL",
			Detail:   "No API URL was provided, using default URL " + apiURL,
		})
	}

	// Create and return the client configuration structure.
	config := ProviderConfig{
		ApiKey:    apiKey,
		BaseURL:   apiURL,
		AccountID: accountID,
	}

	return &config, diags
}

// ProviderConfig holds the configuration details for creating HTTP requests to the Qdrant Cloud API.
// It encapsulates the API key, the base URL, and the HTTP client configured for API communication.
// As well as the (optional) default account ID
type ProviderConfig struct {
	ApiKey    string // ApiKey represents the authentication token used for Qdrant Cloud API access.
	BaseURL   string // BaseURL is the root URL for all API requests, typically pointing to the Qdrant Cloud API endpoint.
	AccountID string // The default Account Identifier for the Qdrant cloud, if any
}
