package ah

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const (
	ImageName      = "centos-7-x64"
	DatacenterID   = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	DatacenterName = "ams1"
	VpsPlanID      = "381347529"
	VpsPlanName    = "start-xs"
	VpsUpgPlanName = "start-m"
	VpsUpgPlanID   = "381347841"
	VolumePlanID   = "381347560"
	VolumePlanName = "hdd2-ash1"
	ClusterID      = ""
	NodeID         = "2486b2f8-f7a6-4207-979b-9b94d93c174e"
	K8SVersion     = "v1.19.0"
	NodePoolName   = "KNP100243"
	NodePoolType   = "public"
)

// IsUUID checks if string is in valid UUID format.
func IsUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func generateHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func datacenterIDBySlug(ctx context.Context, client *ah.APIClient, datacenterSlug string) (string, error) {
	datacenters, err := client.Datacenters.List(ctx, nil)
	if err != nil {
		return "", err
	}
	for _, datacenter := range datacenters {
		if datacenter.Slug == datacenterSlug {
			return datacenter.ID, nil
		}
	}
	return "", fmt.Errorf("datacenter slug %s not found", datacenterSlug)
}

func kubernetesVersion(ctx context.Context, client *ah.APIClient, k8sVersion string) (string, error) {
	versions, err := client.KubernetesClusters.GetKubernetesClustersVersions(ctx)
	if err != nil {
		return "", err
	}
	for _, version := range versions {
		if version == k8sVersion {
			return k8sVersion, nil
		}
	}
	return "", fmt.Errorf("kubernetes version %s not found", k8sVersion)
}

func testAccCheckAHResourceNoRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *beforeID == "" {
			t.Fatalf("Old ID has not been set")
		}
		if *afterID == "" {
			t.Fatalf("New ID has not been set")
		}
		if *beforeID != *afterID {
			t.Fatalf("Resource has been recreated, old ID: %s, new ID: %s", *beforeID, *afterID)
		}
		return nil
	}
}

func testAccCheckAHResourceRecreated(t *testing.T, beforeID, afterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *beforeID == "" {
			t.Fatalf("Old ID has not been set")
		}
		if *afterID == "" {
			t.Fatalf("New ID has not been set")
		}
		if *beforeID == *afterID {
			t.Fatalf("Resource hasn't been recreated, ID: %s", *beforeID)
		}
		return nil
	}
}

func datasourceConfigBasic() string {
	return fmt.Sprintf(`
	data "ah_cloud_images" "test" {
		filter {
			key = "slug"
			values = ["%s"]
		  }
	}`, ImageName)
}
