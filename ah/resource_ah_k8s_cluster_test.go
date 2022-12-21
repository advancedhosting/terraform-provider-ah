package ah

import (
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

import (
	"context"
	//	"fmt"
	//	"testing"
	//
	//	"github.com/advancedhosting/advancedhosting-api-go/ah"
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//func TestAccAHK8sCluster_Basic(t *testing.T) {
//resource.Test(t, resource.TestCase{
//	PreCheck:          func() { testAccPreCheck(t) },
//	ProviderFactories: testAccProviderFactories,
//	CheckDestroy:      testAccCheckAHK8sClusterDestroy,
//	Steps: []resource.TestStep{
//		{
//			Config: testAccCheckAHK8sClusterConfigBasic(),
//			Check: resource.ComposeTestCheckFunc(
//				resource.TestCheckResourceAttrSet("ah_k8s_cluster.web", "id"),
//				resource.TestCheckResourceAttr("ah_k8s_cluster.web", "name", "Test K8s Cluster"),
//				resource.TestCheckResourceAttr("ah_k8s_cluster.web", "count", "1"),
//				resource.TestCheckResourceAttr("ah_k8s_cluster.web", "plan_id", "381347529"),
//				resource.TestCheckResourceAttrSet("ah_k8s_cluster.web", "datacenter"),
//				resource.TestCheckResourceAttrSet("ah_k8s_cluster.web", "state"),
//				resource.TestCheckResourceAttrSet("ah_k8s_cluster.web", "created_at"),
//				resource.TestCheckResourceAttrSet("ah_k8s_cluster.web", "number"),
//			),
//		},
//	},
//})
//}

//func TestAccAHK8sCluster_UpdateName(t *testing.T) {
//	var beforeID, afterID string
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: testAccProviderFactories,
//		CheckDestroy:      testAccCheckAHK8sClusterDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccCheckAHK8sClusterConfigBasic(),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckAHK8sClusterExists("ah_k8s_cluster.test", &beforeID),
//				),
//			},
//			{
//				Config: testAccCheckAHK8sClusterConfigUpdateName(),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckAHK8sClusterExists("ah_k8s_cluster.test", &afterID),
//					resource.TestCheckResourceAttr("ah_k8s_cluster.test", "name", "New K8s Cluster"),
//					testAccCheckAHResourceNoRecreated(t, beforeID, afterID),
//				),
//			},
//		},
//	})
//}

func testAccCheckAHK8sClusterExists(n string, clusterID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No k8s cluster ID is set")
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
			return fmt.Errorf("Error removing k8s cluster (%s): %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHK8sClusterConfigBasic() string {
	return `
	resource "ah_k8s_cluster" "test" {
	  datacenter = "ams1"
	  name = "Test K8s Cluster"
      plan = "381347758"
	  nodes_count = 1
	}`
}

func testAccCheckAHK8sClusterConfigUpdateName() string {
	return `
	resource "ah_k8s_cluster" "test" {
	  datacenter = "ams1"
	  name = "New K8s Cluster"
      plan = "381347758"
	  nodes_count = 1
	}`
}
