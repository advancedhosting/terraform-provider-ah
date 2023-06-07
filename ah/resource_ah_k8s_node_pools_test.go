package ah

import (
	"context"
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccAHK8sNodePool_Basic(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("test-terraform-cluster-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHK8sNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHK8sNodePoolConfigBasic(clusterName),
				//Check: resource.ComposeTestCheckFunc(
				//	resource.TestCheckResourceAttrSet("ah_k8s_node_pool.np", "id"),
				//),
			},
			{
				ResourceName: "ah_k8s_node_pool.np",
				ImportState:  true,
			},
		},
	})
}

func TestAccAHK8sNodePool_UpdateAutoScale(t *testing.T) {
	t.Parallel()

	clusterName := fmt.Sprintf("test-terraform-cluster-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHK8sNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHK8sNodePoolConfigBasic(clusterName),
			},
		},
	})
}

func testAccCheckAHK8sNodePoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_k8s_node_pool" {
			continue
		}

		_, err := client.KubernetesClusters.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf("error removing k8s Node Pool (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHK8sNodePoolConfigBasic(clusterName string) string {
	return fmt.Sprintf(`
resource "ah_k8s_cluster" "cluster" {
	name        = "%s"
	datacenter  = "%s"
	k8s_version = "%s"

	node_pools {
		type              = "%s"
		nodes_count       = 1
		public_properties = {
			plan_id = "%s"
		}
	}
}

resource "ah_k8s_node_pool" "np" {
	cluster_id  = ah_k8s_cluster.cluster.id
	type        = "public"
	nodes_count = 1
	public_properties = {
		plan_id = "%s"
	}
}
`, clusterName, DatacenterName, K8SVersion, NodePoolType, K8sPlanID, K8sPlanID)
}
