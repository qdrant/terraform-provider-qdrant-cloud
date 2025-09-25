package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qcb "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
	qccl "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
	qcco "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
)

const (
	backupFieldTemplate = "Backup Schema %s field"

	// Writable fields.
	backupAccountIdFieldName = "account_id"
	backupClusterIdFieldName = "cluster_id"

	// Read-only fields (per proto).
	backupIdFieldName              = "id"
	backupCreatedAtFieldName       = "created_at"
	backupDeletedAtFieldName       = "deleted_at"
	backupNameFieldName            = "name"
	backupStatusFieldName          = "status"
	backupDurationFieldName        = "backup_duration"
	backupScheduleIdFieldName      = "backup_schedule_id"
	backupClusterInfoFieldName     = "cluster_info"
	bClusterInfoNameField          = "name"
	bClusterInfoCloudProviderField = "cloud_provider_id"
	bClusterInfoRegionField        = "cloud_provider_region_id"
	bClusterCfgField               = "configuration"
	bClusterCfgLastModifiedAtField = "last_modified_at"
	bClusterCfgNumberOfNodesField  = "number_of_nodes"
	bClusterCfgVersionField        = "version"
	bClusterCfgPackageIdField      = "package_id"
	bClusterCfgServiceTypeField    = "service_type"
	bClusterCfgRebalanceStratField = "rebalance_strategy"
	bClusterCfgDbConfigField       = "database_configuration"
	bDbCfgServiceField             = "service"
	bDbCfgServiceApiKeyField       = "api_key"
	bDbCfgServiceApiKeyNameField   = "name"
	bDbCfgServiceApiKeyKeyField    = "key"
	bDbCfgServiceJwtRbacField      = "jwt_rbac"
	bDbCfgInferenceField           = "inference"
	bDbCfgInferenceEnabledField    = "enabled"
)

// accountsBackupSchema defines the Terraform schema for a Backup resource.
// Writable fields are clearly separated from read-only fields for clarity and safety.
func accountsBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Writable
		backupAccountIdFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Account ID"),
			Optional:    true,
			Computed:    true,
			Type:        schema.TypeString,
		},
		backupClusterIdFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Cluster ID"),
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		// Read-only
		backupIdFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "ID"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupCreatedAtFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Creation timestamp"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupDeletedAtFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Deletion timestamp (if applicable)"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupNameFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Auto-generated backup name"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupStatusFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Backup status"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupDurationFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Backup duration (e.g., 36s)"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupScheduleIdFieldName: {
			Description: fmt.Sprintf(backupFieldTemplate, "Backup Schedule ID that produced this backup (if any)"),
			Type:        schema.TypeString,
			Computed:    true,
		},

		// cluster_info (read-only nested)
		backupClusterInfoFieldName: {
			Description: "Cluster metadata captured at backup time (read-only).",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsBackupClusterInfoSchema()},
		},
	}
}

func accountsBackupClusterInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		bClusterInfoNameField: {
			Description: "Cluster name.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterInfoCloudProviderField: {
			Description: "Cloud provider (e.g., aws, gcp, azure).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterInfoRegionField: {
			Description: "Cloud provider region (e.g., eu-central-1).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterCfgField: {
			Description: "Cluster configuration details (read-only).",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsBackupClusterConfigurationSchema()},
		},
	}
}

func accountsBackupClusterConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		bClusterCfgLastModifiedAtField: {
			Description: "Configuration last modified timestamp.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterCfgNumberOfNodesField: {
			Description: "Number of nodes.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		bClusterCfgVersionField: {
			Description: "Qdrant version.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterCfgPackageIdField: {
			Description: "Package ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterCfgServiceTypeField: {
			Description: "Service type (enum as string).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterCfgRebalanceStratField: {
			Description: "Rebalance strategy (enum as string).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		bClusterCfgDbConfigField: {
			Description: "Database configuration at backup time.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsBackupDatabaseConfigurationSchema()},
		},
	}
}

func accountsBackupDatabaseConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		bDbCfgServiceField: {
			Description: "Service settings.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsBackupDatabaseServiceSchema()},
		},
		bDbCfgInferenceField: {
			Description: "Inference settings.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsBackupInferenceSchema()},
		},
	}
}

func accountsBackupDatabaseServiceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		bDbCfgServiceApiKeyField: {
			Description: "Service API key (name/key).",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				bDbCfgServiceApiKeyNameField: {
					Description: "API key name.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				bDbCfgServiceApiKeyKeyField: {
					Description: "API key identifier.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			}},
		},
		bDbCfgServiceJwtRbacField: {
			Description: "JWT RBAC enabled.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
}

func accountsBackupInferenceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		bDbCfgInferenceEnabledField: {
			Description: "Inference enabled.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
}

// flattenBackup maps API object -> Terraform state fields (including read-only).
func flattenBackup(b *qcb.Backup) map[string]interface{} {
	if b == nil {
		return map[string]interface{}{}
	}
	out := map[string]interface{}{
		backupIdFieldName:         b.GetId(),
		backupAccountIdFieldName:  b.GetAccountId(),
		backupClusterIdFieldName:  b.GetClusterId(),
		backupNameFieldName:       b.GetName(),
		backupStatusFieldName:     b.GetStatus().String(),
		backupDurationFieldName:   formatDuration(b.GetBackupDuration()),
		backupScheduleIdFieldName: b.GetBackupScheduleId(),
		backupClusterInfoFieldName: flattenBackupClusterInfo(
			b.GetClusterInfo(),
		),
	}
	if ts := b.GetCreatedAt(); ts != nil {
		out[backupCreatedAtFieldName] = formatTime(ts)
	}
	if ts := b.GetDeletedAt(); ts != nil {
		out[backupDeletedAtFieldName] = formatTime(ts)
	}
	return out
}

func flattenBackupClusterInfo(ci *qcb.ClusterInfo) []interface{} {
	if ci == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		bClusterInfoNameField:          ci.GetName(),
		bClusterInfoCloudProviderField: ci.GetCloudProviderId(),
		bClusterInfoRegionField:        ci.GetCloudProviderRegionId(),
		bClusterCfgField:               flattenBackupClusterConfiguration(ci.GetConfiguration()),
	}
	return []interface{}{m}
}

func flattenBackupClusterConfiguration(cfg *qccl.ClusterConfiguration) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		bClusterCfgNumberOfNodesField:  int(cfg.GetNumberOfNodes()),
		bClusterCfgVersionField:        cfg.GetVersion(),
		bClusterCfgPackageIdField:      cfg.GetPackageId(),
		bClusterCfgServiceTypeField:    cfg.GetServiceType().String(),
		bClusterCfgRebalanceStratField: cfg.GetRebalanceStrategy().String(),
		bClusterCfgDbConfigField:       flattenBackupDatabaseConfiguration(cfg.GetDatabaseConfiguration()),
	}
	if ts := cfg.GetLastModifiedAt(); ts != nil {
		m[bClusterCfgLastModifiedAtField] = formatTime(ts)
	}
	return []interface{}{m}
}

func flattenBackupDatabaseConfiguration(db *qccl.DatabaseConfiguration) []interface{} {
	if db == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			bDbCfgServiceField:   flattenBackupDatabaseService(db.GetService()),
			bDbCfgInferenceField: flattenBackupInference(db.GetInference()),
		},
	}
}

func flattenBackupDatabaseService(svc *qccl.DatabaseConfigurationService) []interface{} {
	if svc == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			bDbCfgServiceApiKeyField:  flattenBackupAPIKey(svc.GetApiKey()),
			bDbCfgServiceJwtRbacField: svc.GetJwtRbac(),
		},
	}
}

func flattenBackupAPIKey(k *qcco.SecretKeyRef) []interface{} {
	if k == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			bDbCfgServiceApiKeyNameField: k.GetName(),
			bDbCfgServiceApiKeyKeyField:  k.GetKey(),
		},
	}
}

func flattenBackupInference(inf *qccl.DatabaseConfigurationInference) []interface{} {
	if inf == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			bDbCfgInferenceEnabledField: inf.GetEnabled(),
		},
	}
}

// expandBackup builds the API object from TF config for Create.
func expandBackup(d *schema.ResourceData, defaultAccountID string) (*qcb.Backup, error) {
	accountID := defaultAccountID
	if v, ok := d.GetOk(backupAccountIdFieldName); ok && v.(string) != "" {
		accountID = v.(string)
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID not specified")
	}
	clusterID := d.Get(backupClusterIdFieldName).(string)
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID must be set")
	}

	return &qcb.Backup{
		AccountId: accountID,
		ClusterId: clusterID,
		// read-only fields: (id, created_at, name, status, deleted_at,
		// backup_duration, backup_schedule_id, cluster_info)
	}, nil
}
