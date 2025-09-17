package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qch "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/hybrid/v1"
)

const (
	hcEnvIdFieldName                = "id"
	hcEnvAccountIdFieldName         = "account_id"
	hcEnvNameFieldName              = "name"
	hcEnvConfigurationFieldName     = "configuration"
	hcEnvCfgNamespaceFieldName      = "namespace"
	hcEnvCreatedAtFieldName         = "created_at"
	hcEnvLastModifiedAtFieldName    = "last_modified_at"
	hcEnvBootstrapCommandsFieldName = "bootstrap_commands"
)

func accountsHybridCloudEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvIdFieldName: {
			Type:     schema.TypeString,
			Computed: true, // mirror of d.Id()
		},
		hcEnvAccountIdFieldName: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		hcEnvNameFieldName: {
			Type:     schema.TypeString,
			Required: true,
		},
		hcEnvConfigurationFieldName: {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					hcEnvCfgNamespaceFieldName: {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},
		hcEnvCreatedAtFieldName: {
			Type:     schema.TypeString,
			Computed: true,
		},
		hcEnvLastModifiedAtFieldName: {
			Type:     schema.TypeString,
			Computed: true,
		},

		// Intentionally no `status` in schema.

		hcEnvBootstrapCommandsFieldName: {
			Type:      schema.TypeList,
			Elem:      &schema.Schema{Type: schema.TypeString},
			Computed:  true,
			Sensitive: true,
		},
	}
}

// flattenHCEnv maps API object -> Terraform state fields (excluding status).
func flattenHCEnv(env *qch.HybridCloudEnvironment) map[string]interface{} {
	out := map[string]interface{}{
		hcEnvIdFieldName:             env.GetId(),
		hcEnvAccountIdFieldName:      env.GetAccountId(),
		hcEnvNameFieldName:           env.GetName(),
		hcEnvCreatedAtFieldName:      formatTime(env.GetCreatedAt()),
		hcEnvLastModifiedAtFieldName: formatTime(env.GetLastModifiedAt()),
	}

	cfg := env.GetConfiguration()
	if cfg != nil {
		out[hcEnvConfigurationFieldName] = []interface{}{
			map[string]interface{}{
				hcEnvCfgNamespaceFieldName: cfg.GetNamespace(),
			},
		}
	} else {
		out[hcEnvConfigurationFieldName] = []interface{}{}
	}

	return out
}

// expandHCEnvForCreate builds the create payload from config; validates required fields.
func expandHCEnvForCreate(d *schema.ResourceData, defaultAccountID string) (*qch.HybridCloudEnvironment, error) {
	accountID := defaultAccountID
	if v, ok := d.GetOk(hcEnvAccountIdFieldName); ok && v.(string) != "" {
		accountID = v.(string)
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID not specified")
	}

	name := d.Get(hcEnvNameFieldName).(string)

	var namespace string
	if v, ok := d.GetOk(hcEnvConfigurationFieldName); ok {
		l := v.([]interface{})
		if len(l) > 0 && l[0] != nil {
			m := l[0].(map[string]interface{})
			if ns, ok := m[hcEnvCfgNamespaceFieldName]; ok && ns.(string) != "" {
				namespace = ns.(string)
			}
		}
	}
	if namespace == "" {
		return nil, fmt.Errorf("configuration.namespace must be set")
	}

	return &qch.HybridCloudEnvironment{
		AccountId: accountID,
		Name:      name,
		Configuration: &qch.HybridCloudEnvironmentConfiguration{
			Namespace: namespace,
		},
	}, nil
}

// expandHCEnvForUpdate builds the update payload; requires ID and resolved account ID.
func expandHCEnvForUpdate(d *schema.ResourceData, defaultAccountID string) (*qch.HybridCloudEnvironment, error) {
	accountID := defaultAccountID
	if v, ok := d.GetOk(hcEnvAccountIdFieldName); ok && v.(string) != "" {
		accountID = v.(string)
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID not specified")
	}
	id := d.Id()
	if id == "" {
		return nil, fmt.Errorf("resource ID not set")
	}

	env := &qch.HybridCloudEnvironment{
		AccountId: accountID,
		Id:        id,
		Name:      d.Get(hcEnvNameFieldName).(string),
	}

	if v, ok := d.GetOk(hcEnvConfigurationFieldName); ok {
		l := v.([]interface{})
		if len(l) > 0 && l[0] != nil {
			m := l[0].(map[string]interface{})
			if ns, ok := m[hcEnvCfgNamespaceFieldName]; ok && ns.(string) != "" {
				env.Configuration = &qch.HybridCloudEnvironmentConfiguration{
					Namespace: ns.(string),
				}
			}
		}
	}

	return env, nil
}
