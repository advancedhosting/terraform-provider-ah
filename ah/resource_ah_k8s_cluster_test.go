package ah

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	K8sPlanID = "381347758"
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
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "id"),
					resource.TestCheckResourceAttr("ah_k8s_cluster.test", "name", name),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "nodes_count"),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "plan"),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "state"),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "created_at"),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "number"),
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
					testAccCheckAHK8sClusterExists("ah_k8s_cluster.test", &beforeID),
				),
			},
			{
				Config: testAccCheckAHK8sClusterConfigUpdateName(newName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAHK8sClusterExists("ah_k8s_cluster.test", &afterID),
					resource.TestCheckResourceAttr("ah_k8s_cluster.test", "name", newName),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "state"),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "created_at"),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "number"),
					resource.TestCheckResourceAttrSet("ah_k8s_cluster.test", "plan"),
					testAccCheckAHResourceNoRecreated(t, &beforeID, &afterID),
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

		_, err := client.Clusters.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("error removing k8s cluster (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHK8sClusterConfigBasic(name string) string {
	return fmt.Sprintf(`
	resource "ah_k8s_cluster" "test" {
	  datacenter = "%s"
	  name = "%s"
      plan = "%s"
	  nodes_count = 1
	}`, DatacenterName, name, K8sPlanID)
}

func testAccCheckAHK8sClusterConfigUpdateName(name string) string {
	return fmt.Sprintf(`
	resource "ah_k8s_cluster" "test" {
	  datacenter = "%s"
	  name = "%s"
      plan = "%s"
	  nodes_count = 1
	}`, DatacenterName, name, K8sPlanID)
}
