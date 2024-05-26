package qdrant

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// dataClusterAccountsCluster constructs a Terraform resource for
// managing the reading of a specific cluster associated with an account.
func dataClusterAccountsCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataClusterAccountsClusterRead,
		Schema:      AccountsClustersSchema(),
	}
}

// dataClusterAccountsClusters constructs a Terraform resource for
// managing the reading of all clusters associated with an account.
func dataClusterAccountsClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataClusterAccountsClustersRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: AccountsClustersSchema(),
				},
			},
		},
	}
}

// dataClusterAccountsClusterRead performs a read operation to fetch a specific cluster associated with a given account and cluster ID.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataClusterAccountsClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID := d.Get("account_id").(string)
	clusterID := d.Get("cluster_id").(string)

	apiClient, err, diagnostics, done := GetClient(m)

	if done {
		return diagnostics
	}

	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error parsing account ID: %v", err))
		if d.HasError() {
			return d
		}
	}
	response, err := apiClient.ListClustersWithResponse(ctx, accountUUID, &qc.ListClustersParams{})
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error listing clusters: %v", err))
		if d.HasError() {
			return d
		}
	}

	if response.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error listing clusters: %v", response.JSON422))
	}

	cluster := *response.JSON200

	// Update the Terraform state
	if err := d.Set("cluster", cluster); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID)
	return nil
}

// dataClusterAccountsClustersRead performs a read operation to fetch all clusters associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataClusterAccountsClustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiKey := m.(ClientConfig).ApiKey
	accountID := d.Get("account_id").(string)

	opts := qc.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", apiKey))
		return nil
	})

	apiClient, err := qc.NewClientWithResponses(m.(ClientConfig).BaseURL, opts)
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error initializing client: %v", err))
		if d.HasError() {
			return d
		}
	}

	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error parsing account ID: %v", err))
		if d.HasError() {
			return d
		}
	}

	response, err := apiClient.ListClustersWithResponse(ctx, accountUUID, &qc.ListClustersParams{})
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error listing clusters: %v", err))
		if d.HasError() {
			return d
		}
	}

	if response.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error listing clusters: %v", response.JSON422))
	}

	clusters := *response.JSON200

	// Update the Terraform state
	if err := d.Set("clusters", clusters); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
