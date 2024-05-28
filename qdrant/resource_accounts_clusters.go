package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

const (
	clusterFieldTemplate = "Cluster Resource %s field"

	clusterIdentifierFieldName                  = "id"
	clusterCreatedAtFieldName                   = "created_at"
	clusterAccountIDFieldName                   = "account_id"
	clusterNameFieldName                        = "name"
	clusterCloudProviderFieldName               = "cloud_provider"
	clusterCloudRegionFieldName                 = "cloud_region"
	clusterCloudRegionAvailabilityZoneFieldName = "cloud_region_az"
	clusterVersionFieldName                     = "version"
)

// resourceAccountsClusters constructs a Terraform resource for
// managing the creation, reading, and deletion of clusters associated with an account.
func resourceAccountsClusters() *schema.Resource {
	return &schema.Resource{
		Description:   "Account Cluster Resource",
		ReadContext:   resourceClusterRead,
		CreateContext: resourceClusterCreate,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Schema: map[string]*schema.Schema{
			clusterIdentifierFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the cluster"),
				Type:        schema.TypeString,
				Computed:    true,
			},
			clusterCreatedAtFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Timestamp then the cluster is created"),
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner_id": { // TODO: Remove?
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the owner"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			clusterAccountIDFieldName: { // If set here, overrides account ID in provider
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the account"),
				Type:        schema.TypeString,
				Optional:    true,
			},
			clusterNameFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Name of the cluster"),
				Type:        schema.TypeString,
				Required:    true,
			},
			clusterCloudProviderFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud provider where the cluster resides"),
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Cross provider migration isn't supported
			},
			clusterCloudRegionFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region where the cluster resides"),
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Cross region migration isn't supported
			},
			clusterCloudRegionAvailabilityZoneFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region availability zone where the cluster resides"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"cloud_region_setup": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region setup of the cluster"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"private_region_id": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the Private Region"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"current_configuration_id": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"encryption_key_id": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"marked_for_deletion_at": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			clusterVersionFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Version of the qdrant cluster"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"url": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "TODO",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"configuration": {
				Description: "TODO",
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Description: "TODO",
					Schema: map[string]*schema.Schema{
						"num_nodes_max": {
							Description: "TODO",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"num_nodes": {
							Description: "TODO",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"node_configuration": {
							Description: "TODO",
							Type:        schema.TypeSet,
							Required:    true,
							ForceNew:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Description: "TODO",
								Schema: map[string]*schema.Schema{
									"package_id": {
										Description: "TODO",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"resources": {
				Description: "TODO",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Description: "TODO",
					Type:        schema.TypeString,
				},
			},
			"total_extra_disk": {
				Description: "TODO",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

// resourceClusterRead reads the specific cluster's data from the API.
func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}

	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := apiClient.GetClusterWithResponse(ctx, accountID, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error reading cluster: %v", resp.JSON422))
	}

	if err := d.Set("clusters", resp.JSON200); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	accountID := d.Get("account_id").(string)

	cluster := expandClusterIn(d)
	accID, err := uuid.Parse(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := apiClient.CreateClusterWithResponse(ctx, accID, cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error creating cluster: %v", resp.JSON422))
	}

	clusterOut := resp.JSON200
	if clusterOut == nil {
		return diag.FromErr(fmt.Errorf("error creating cluster: no cluster returned"))
	}

	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(clusterOut) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	accountUUID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterUUID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("configuration") {
		conf := expandClusterConfigurationIn(d.Get("configuration").(*schema.Set))

		resp, err := apiClient.UpdateClusterWithResponse(ctx, accountUUID, clusterUUID, qc.PydanticClusterPatchIn{
			Configuration: &qc.PydanticClusterConfigurationPatchIn{
				NumNodes: &conf.NumNodes,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if resp.JSON422 != nil {
			return diag.FromErr(fmt.Errorf("error updating cluster: %v", resp.JSON422))
		}
	}
	return nil

}

// resourceClusterDelete performs a delete operation to remove a cluster associated with an account.
func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBackups := true
	params := qc.DeleteClusterParams{
		DeleteBackups: &deleteBackups,
	}
	resp, err := apiClient.DeleteClusterWithResponse(ctx, accountID, clusterID, &params)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error deleting cluster: %v", resp.JSON422))
	}

	d.SetId("")
	return nil
}

func expandClusterIn(d *schema.ResourceData) qc.ClusterIn {
	accountID := d.Get("account_id").(string)

	name := d.Get("name")
	cloudProvider := d.Get("cloud_provider")
	cloudRegion := d.Get("cloud_region")

	cluster := qc.ClusterIn{
		Name:          name.(string),
		CloudProvider: qc.ClusterInCloudProvider(cloudProvider.(string)),
		CloudRegion:   qc.ClusterInCloudRegion(cloudRegion.(string)),
		AccountId:     &accountID,
	}
	if v, ok := d.GetOk("configuration"); ok {
		configuration := expandClusterConfigurationIn(v.(*schema.Set))
		cluster.Configuration = *configuration
	}
	return cluster
}

func expandClusterConfigurationIn(v *schema.Set) *qc.ClusterConfigurationIn {
	if len(v.List()) == 0 {
		return nil
	}

	m := v.List()[0].(map[string]interface{})
	config := qc.ClusterConfigurationIn{}

	if v, ok := m["num_nodes_max"]; ok {
		config.NumNodesMax = v.(int)
	}
	if v, ok := m["num_nodes"]; ok {
		config.NumNodes = v.(int)
	}

	if v, ok := m["node_configuration"]; ok {
		nodeConfig := expandNodeConfigurationIn(v.(*schema.Set))
		if nodeConfig != nil {
			config.NodeConfiguration = *nodeConfig
		}
	}

	return &config
}

func expandNodeConfigurationIn(v *schema.Set) *qc.NodeConfiguration {
	if len(v.List()) == 0 {
		return nil
	}

	m := v.List()[0].(map[string]interface{})
	config := qc.NodeConfiguration{}

	if v, ok := m["package_id"]; ok {
		config.PackageId = v.(string)
	}

	return &config
}

// flattenCluster creates a map from a cluster for easy storage on terraform.
func flattenCluster(cluster *qc.ClusterOut) map[string]interface{} {
	result := map[string]interface{}{
		clusterIdentifierFieldName: cluster.Id,
		"name":                     cluster.Name,
		"cloud_provider":           cluster.CloudProvider,
		"cloud_region":             cluster.CloudRegion,
		"configuration":            flattenClusterConfiguration(cluster.Configuration),
	}
	return result
}

// flattenClusterConfiguration creates a map from a cluster for easy storage on terraform.
func flattenClusterConfiguration(clusterConfig *qc.ClusterConfigurationOut) map[string]interface{} {
	result := map[string]interface{}{
		"id":            clusterConfig.Id,
		"num_nodes":     clusterConfig.NumNodes,
		"num_nodes_max": clusterConfig.NumNodesMax,
	}
	return result
}
