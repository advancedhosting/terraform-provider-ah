package ah

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	K8sPlanID    = "381347758"
	resourceName = "ah_k8s_cluster.ah_test_cluster"
)

func TestAccAHK8sCluster_Basic(t *testing.T) {
	name := fmt.Sprintf("test-terraform-cluster-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHK8sClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHK8sClusterConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "datacenter"),
					resource.TestCheckResourceAttrSet(resourceName, "private_network"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "number"),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttr(resourceName, "k8s_version", K8SVersion),

					resource.TestCheckResourceAttrSet(resourceName, "node_pools.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "node_pools.0.name"),
					resource.TestCheckResourceAttr(resourceName, "node_pools.0.type", WorkerPoolType),
					resource.TestCheckResourceAttr(resourceName, "node_pools.0.nodes_count", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "node_pools.0.labels.%"),
					resource.TestCheckResourceAttr(resourceName, "node_pools.0.public_properties.plan_id", K8sPlanID),
				),
			},
		},
	})
}

func TestAccAHK8sCluster_UpdateName(t *testing.T) {
	name := fmt.Sprintf("test-terraform-cluster-%s", acctest.RandString(5))
	newName := fmt.Sprintf("test-terraform-cluster-%s", acctest.RandString(5))

	var beforeID, afterID string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHK8sClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHK8sClusterConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHK8sClusterExists(resourceName, &beforeID),
				),
			},
			{
				Config: testAccCheckAHK8sClusterConfigUpdateName(newName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHK8sClusterExists(resourceName, &afterID),
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "number"),
					testAccCheckAHResourceNoRecreated(t, &beforeID, &afterID),
					resource.TestCheckResourceAttr(resourceName, "node_pools.0.public_properties.plan_id", K8sPlanID),
				),
			},
		},
	})
}

func testAccCheckAHK8sClusterExists(n string, clusterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("тot found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no k8s cluster ID is set")
		}

		*clusterID = rs.Primary.ID
		return nil
	}
}

func testAccCheckAHK8sClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_k8s_cluster" {
			continue
		}

		_, err := client.KubernetesClusters.Get(context.Background(), rs.Primary.ID)

		if !errors.Is(err, ah.ErrResourceNotFound) {
			return fmt.Errorf("error removing k8s cluster (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHK8sClusterConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "ah_k8s_cluster" "ah_test_cluster" {
	name                      = "%s"
    datacenter                = "%s"
    k8s_version               = "%s"
    node_pools {
			type              = "%s"
			nodes_count       = 1
			labels            = {
				"labels.websa.com/technologies": "terraform",
				"labels.websa.com/terraform": "default-node-pool",
			}
			public_properties = {
				plan_id = "%s"
			}
		}
	}
    `, name, DatacenterName, K8SVersion, WorkerPoolType, K8sPlanID)
}

func testAccCheckAHK8sClusterConfigUpdateName(name string) string {
	return fmt.Sprintf(`
resource "ah_k8s_cluster" "ah_test_cluster" {
	name                      = "%s"
    datacenter                = "%s"
    k8s_version               = "%s"
    node_pools {
			type              = "%s"
			nodes_count       = 1
			labels            = {
				"labels.websa.com/technologies": "terraform",
				"labels.websa.com/terraform": "default-node-pool",
			}
			public_properties = {
				plan_id = "%s"
			}
		}
	}
    `, name, DatacenterName, K8SVersion, WorkerPoolType, K8sPlanID)
}
