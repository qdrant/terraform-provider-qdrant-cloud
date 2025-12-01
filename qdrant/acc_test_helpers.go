package qdrant

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	qcb "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
)

// testCheckHasListAttr verifies a list attribute has at least one element.
func testCheckHasListAttr(resourceName, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if n, ok := rs.Primary.Attributes[attr+".#"]; !ok || n == "0" {
			return fmt.Errorf("expected %s to have at least one element, got %q", attr, n)
		}
		return nil
	}
}

// testAccCheckListNonEmpty returns a TestCheckFunc that asserts the given resource list attribute is present and non-empty in state.
func testAccCheckListNonEmpty(name, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", name)
		}
		// Terraform stores list length at "<attr>.#"
		if n := rs.Primary.Attributes[attr+".#"]; n == "" || n == "0" {
			return fmt.Errorf("expected %s to be non-empty, got length=%q", attr, n)
		}
		return nil
	}
}

// testAccCheckAttrEqual asserts name1.attr1 == name2.attr2 (string compare).
func testAccCheckAttrEqual(name1, attr1, name2, attr2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, ok := s.RootModule().Resources[name1]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", name1)
		}
		rs2, ok := s.RootModule().Resources[name2]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", name2)
		}
		v1, ok := rs1.Primary.Attributes[attr1]
		if !ok {
			return fmt.Errorf("attribute %s not found on %s", attr1, name1)
		}
		v2, ok := rs2.Primary.Attributes[attr2]
		if !ok {
			return fmt.Errorf("attribute %s not found on %s", attr2, name2)
		}
		if v1 != v2 {
			return fmt.Errorf("expected %s.%s (%q) to equal %s.%s (%q)", name1, attr1, v1, name2, attr2, v2)
		}
		return nil
	}
}

// testAccWaitForBackupStatus polls the backup service until the given resource reaches the desired status or times out.
// It fails fast if the backup transitions to FAILED.
func testAccWaitForBackupStatus(resourceName string, desired qcb.BackupStatus, timeout time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}
		backupID := rs.Primary.ID
		if backupID == "" {
			return fmt.Errorf("resource %s missing ID", resourceName)
		}
		accountID := rs.Primary.Attributes["account_id"]
		if accountID == "" {
			accountID = os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
		}
		if accountID == "" {
			return fmt.Errorf("account_id not found for %s", resourceName)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		client, clientCtx, closeFn, err := newTestAccBackupServiceClient(ctx)
		if err != nil {
			return err
		}
		defer closeFn()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		var lastStatus qcb.BackupStatus
		var lastErr error
		for {
			resp, err := client.GetBackup(clientCtx, &qcb.GetBackupRequest{
				AccountId: accountID,
				BackupId:  backupID,
			})
			if err != nil {
				lastErr = err
			} else if resp.GetBackup() != nil {
				lastStatus = resp.GetBackup().GetStatus()
				switch lastStatus {
				case desired:
					return nil
				case qcb.BackupStatus_BACKUP_STATUS_FAILED:
					return fmt.Errorf("backup %s entered failed status while waiting for %s", backupID, desired)
				}
			}

			select {
			case <-ctx.Done():
				if lastErr != nil {
					return fmt.Errorf("timeout waiting for backup %s to reach %s (last error: %w)", backupID, desired, lastErr)
				}
				return fmt.Errorf("timeout waiting for backup %s to reach %s (last status: %s)", backupID, desired, lastStatus)
			case <-ticker.C:
			}
		}
	}
}

// newTestAccBackupServiceClient dials the backup service using the TF_ACC credentials and returns a client, context with auth headers, and close func.
func newTestAccBackupServiceClient(ctx context.Context) (qcb.BackupServiceClient, context.Context, func(), error) {
	apiKey := strings.TrimSpace(os.Getenv("QDRANT_CLOUD_API_KEY"))
	if apiKey == "" {
		return nil, nil, nil, fmt.Errorf("QDRANT_CLOUD_API_KEY not set")
	}
	apiURL := strings.TrimSpace(os.Getenv("QDRANT_CLOUD_API_URL"))
	if apiURL == "" {
		apiURL = "grpc.cloud.qdrant.io"
	}
	insecure := strings.EqualFold(os.Getenv("QDRANT_CLOUD_INSECURE"), "true")
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: insecure})
	conn, err := grpc.NewClient(apiURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("dial backup service: %w", err)
	}
	clientCtx := metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("apikey %s", apiKey))
	closeFn := func() { _ = conn.Close() }
	return qcb.NewBackupServiceClient(conn), clientCtx, closeFn, nil
}
