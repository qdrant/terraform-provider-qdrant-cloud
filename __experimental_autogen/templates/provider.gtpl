package qdrantcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

    "qdrant-terraform-automation/qdrantcloud/resources"
)

// ClientConfig holds the configuration details for creating HTTP requests to the Qdrant Cloud API.
// It encapsulates the API key, the base URL, and the HTTP client configured for API communication.
type ClientConfig struct {
	ApiKey     string       // ApiKey represents the authentication token used for Qdrant Cloud API access.
	BaseURL    string       // BaseURL is the root URL for all API requests, typically pointing to the Qdrant Cloud API endpoint.
	HTTPClient *http.Client // HTTPClient is the custom configured HTTP client used for making API requests.
}

// Provider defines and returns a Terraform resource provider for Qdrant Cloud.
// It sets up the provider schema, resources, and data sources.
// Returns a pointer to the schema.Provider object.
func Provider() *schema.Provider {
	// Schema defines the provider's configuration options.
    return &schema.Provider{
        Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,                                  // Data type of the API key.
				Required:    true,                                               // API key is a required field.
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_API_KEY", nil), // Default can be set via an environment variable.
				Description: "The API Key for Qdrant Cloud API operations.",     // Description of the API key usage.
			},
			"api_url": {
				Type:        schema.TypeString,                                                                  // Data type of the API URL.
				Optional:    true,                                                                               // API URL is an optional field, with a default provided.
				DefaultFunc: schema.EnvDefaultFunc("QDRANT_CLOUD_API_URL", "https://cloud.qdrant.io/public/v1"), // Default API URL.
				Description: "The URL of the Qdrant Cloud API.",                                                 // Description of the API URL.
			},
        },
		// ResourcesMap defines all the resources that this provider offers.
        ResourcesMap: map[string]*schema.Resource{
            {{- range .OperationGroups }}
			{{- if ne (humanize .Name | snakize) "booking" }}
			"qdrantcloud_{{ humanize .Name | snakize }}": resources.{{ pascalize .Name }}(),
			{{- end }}
			{{- end }}
        },
		// DataSourcesMap defines all the data sources that this provider offers.
        DataSourcesMap: map[string]*schema.Resource{
            {{- range .OperationGroups }}
            "qdrantcloud_{{ humanize .Name | snakize }}": resources.DataResource{{ pascalize .Name }}(),
            {{- end }}
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

    // Validate that the API key is not empty, returning an error diagnostic if it is.
    if strings.TrimSpace(apiKey) == "" {
        return nil, diag.Errorf("api_key must not be empty")
    }

    // Validate that the API URL is not empty, returning an error diagnostic if it is.
    if strings.TrimSpace(apiURL) == "" {
        return nil, diag.Errorf("api_url must not be empty")
    }

    // Configure the HTTP client with a 30-second timeout and custom transport settings.
    client := &http.Client{
        Timeout: time.Second * 30,
        Transport: &http.Transport{
            MaxIdleConns:    100,              // Max number of idle connections in the pool.
            IdleConnTimeout: 90 * time.Second, // Max time an idle connection will remain idle before closing.
        },
    }

    // Create and return the client configuration structure.
    config := ClientConfig{
        ApiKey:     apiKey,
        BaseURL:    apiURL,
        HTTPClient: client,
    }

    return config, nil
}


// newQdrantCloudRequest creates a new HTTP request with the appropriate headers set for Qdrant Cloud API interaction.
// method: HTTP method (e.g., "GET", "POST", "PUT", "DELETE")
// urlPath: The API endpoint after the base URL
// payload: An optional interface{} that, if provided, will be serialized as JSON and set as the request body.
func newQdrantCloudRequest(config ClientConfig, method, urlPath string, payload interface{}) (*http.Request, diag.Diagnostics) {
	// Construct the full URL.
	fullURL := fmt.Sprintf("%s%s", config.BaseURL, urlPath)

	// Serialize the payload to JSON if it's provided.
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return nil, diag.FromErr(fmt.Errorf("failed to encode payload to JSON: %s", err))
		}
	}

	// Create the HTTP request.
	req, err := http.NewRequest(method, fullURL, &body)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to create new HTTP request: %s", err))
	}

	// Add the required headers.
	req.Header.Add("Authorization", fmt.Sprintf("apikey %s", config.ApiKey))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// formatTime converts a time value to a standardized string format.
// t: The time value which can be of type time.Time or string.
// Returns a formatted time string in RFC3339 format if the input is of type time.Time,
// returns the input string unchanged if it is of type string, or an empty string for other types.
func formatTime(t interface{}) string {
	switch v := t.(type) {
		case time.Time:
			// Format time.Time to RFC3339 standard string format.
			return v.Format(time.RFC3339)
		case string:
			// Return string as is.
			return v
		default:
			// Return empty string for other types.
			return ""
	}
}

// ExecuteRequest executes the HTTP request and handles the response.
// client: The client configuration including the API key and the base URL.
// req: The HTTP request to execute.
func ExecuteRequest(client ClientConfig, req *http.Request) (*http.Response, diag.Diagnostics) {
	// Print the curl command equivalent of the request for debugging purposes.
	DebugPrintCurlCommand(req)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// Check for non-200 status code and handle errors
	if resp.StatusCode == http.StatusUnprocessableEntity {
		return nil, handleNon422Response(resp)
	}

	return resp, nil
}
