package qdrant

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"reflect"
	"strings"
	qc "terraform-provider-qdrant-cloud/v1/internal/client"
	"time"
)

type ProviderConfig struct {
	ApiKey    string // ApiKey represents the authentication token used for Qdrant Cloud API access.
	BaseURL   string // BaseURL is the root URL for all API requests, typically pointing to the Qdrant Cloud API endpoint.
	AccountID string // The default Account Identifier for the Qdrant cloud, if any
}

// getClient creates a client from the provided interface
// This client already contains the Authorization needed to invoke the API
// Returns: The client to call the backend API, TF Diagnostics
func getClient(m interface{}) (*qc.ClientWithResponses, diag.Diagnostics) {
	clientConfig, ok := m.(*ProviderConfig)
	if !ok {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: provided interface cannot be casted to ClientConfig"))
	}
	optsCallback := qc.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", clientConfig.ApiKey))
		return nil
	})
	apiClient, err := qc.NewClientWithResponses(clientConfig.BaseURL, optsCallback)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: %v", err))
	}
	return apiClient, nil
}

func getAccountUUID(d *schema.ResourceData, m interface{}) (uuid.UUID, error) {
	accountID := d.Get("account_id").(string)
	if accountID == "" {
		config := m.(*ProviderConfig)
		accountID = config.AccountID
	}
	if accountID == "" {
		return uuid.Nil, fmt.Errorf("account_id is not set")
	}
	return uuid.Parse(accountID)
}

func getError(err interface{}) string {
	if httpError, ok := err.(*qc.HTTPValidationError); ok && httpError.Detail != nil {
		details := make([]string, len(*httpError.Detail))
		for i, detail := range *httpError.Detail {
			details[i] = fmt.Sprintf("%s: %s", detail.Loc, detail.Msg)
		}
		return fmt.Sprintf("Validation error: %s", strings.Join(details, "; "))
	}
	return fmt.Sprintf("%v", err)
}

func getErrorMessage(prefix string, resp *http.Response) string {
	if resp == nil {
		return prefix + ": No response"
	}
	return fmt.Sprintf("%s: [Status: %d] - %s", prefix, resp.StatusCode, resp.Status)
}

func flattenValue(v reflect.Value, result map[string]interface{}, prefix string) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			fieldValue := v.Field(i)

			if field.PkgPath != "" {
				continue
			}

			tag := field.Tag.Get("json")
			name := strings.Split(tag, ",")[0]
			if name == "" {
				name = strings.ToLower(field.Name)
			}

			newPrefix := prefix + name

			if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				result[newPrefix] = fieldValue.Interface().(time.Time).Format(time.RFC3339)
			} else {
				flattenValue(fieldValue, result, newPrefix+".")
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			flattenValue(v.MapIndex(key), result, prefix+key.String()+".")
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			flattenValue(v.Index(i), result, prefix+fmt.Sprintf("%d.", i))
		}
	default:
		result[strings.TrimSuffix(prefix, ".")] = v.Interface()
	}
}

func flattenCreateAuthKey(key interface{}) map[string]interface{} {
	return flatten(key)
}

func flattenGetAuthKey(key interface{}) map[string]interface{} {
	return flatten(key)
}

func flattenGetAuthKeys(keys []qc.GetApiKeyOut) []interface{} {
	result := make([]interface{}, len(keys))
	for i, key := range keys {
		result[i] = flattenGetAuthKey(key)
	}
	return result
}

func flattenResponse(response interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if m, ok := response.(map[string]interface{}); ok {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	apiURL := d.Get("api_url").(string)
	var accountID string
	if aid, ok := d.GetOk("account_id"); ok {
		accountID = aid.(string)
	}
	var diags diag.Diagnostics

	if strings.TrimSpace(apiKey) == "" {
		return nil, diag.Errorf("api_key must not be empty")
	}

	if strings.TrimSpace(apiURL) == "" {
		apiURL = "https://cloud.qdrant.io/public/v0"
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Using default URL",
			Detail:   "No API URL was provided, using default URL " + apiURL,
		})
	}

	config := &ProviderConfig{
		ApiKey:    apiKey,
		BaseURL:   apiURL,
		AccountID: accountID,
	}

	return config, diags
}

func flattenClustersData(clusters []qc.ClusterOut) []interface{} {
	if clusters == nil {
		return make([]interface{}, 0)
	}

	var clustersList []interface{}
	for _, cluster := range clusters {
		clusterMap := flatten(cluster)

		// Обработка вложенных структур
		if cluster.Configuration != nil {
			clusterMap["configuration"] = []interface{}{flatten(*cluster.Configuration)}
		}
		if cluster.State != nil {
			clusterMap["state"] = []interface{}{flatten(*cluster.State)}
		}
		if cluster.Resources != nil {
			clusterMap["resources"] = []interface{}{flatten(*cluster.Resources)}
		}

		clustersList = append(clustersList, clusterMap)
	}

	return clustersList
}

func flattenClusterConfiguration(config qc.ClusterConfigurationOut) map[string]interface{} {
	configMap := flatten(config)

	if config.NodeConfiguration != nil {
		configMap["node_configuration"] = []interface{}{flatten(*config.NodeConfiguration)}
	}

	return configMap
}

func flattenNodeConfiguration(nodeConfig qc.NodeConfiguration) map[string]interface{} {
	nodeConfigMap := flatten(nodeConfig)

	if nodeConfig.Package != nil {
		nodeConfigMap["package"] = []interface{}{flatten(*nodeConfig.Package)}
	}

	return nodeConfigMap
}

func flatten(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		switch field.Kind() {
		case reflect.Struct:
			result[fieldName] = flatten(field.Interface())
		case reflect.Slice, reflect.Array:
			sliceResult := make([]interface{}, field.Len())
			for j := 0; j < field.Len(); j++ {
				sliceResult[j] = flatten(field.Index(j).Interface())
			}
			result[fieldName] = sliceResult
		default:
			result[fieldName] = field.Interface()
		}
	}

	return result
}
