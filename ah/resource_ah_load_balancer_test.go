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

func TestAccAHLoadBalancer_Basic(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHLoadBalancerConfig_Basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "name", name),
					resource.TestCheckResourceAttr("ah_load_balancer.web", "state", "active"),
					resource.TestCheckResourceAttr("ah_load_balancer.web", "balancing_algorithm", "round_robin"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "datacenter"),

					resource.TestCheckResourceAttr("ah_load_balancer.web", "ip_address.#", "1"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "ip_address.0.id"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "ip_address.0.type"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "ip_address.0.address"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "ip_address.0.state"),

					resource.TestCheckResourceAttr("ah_load_balancer.web", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "private_network.0.id"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "private_network.0.state"),

					resource.TestCheckResourceAttr("ah_load_balancer.web", "backend_node.#", "2"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "backend_node.0.id"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "backend_node.0.cloud_server_id"),

					resource.TestCheckResourceAttr("ah_load_balancer.web", "forwarding_rule.#", "2"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "forwarding_rule.0.id"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "forwarding_rule.0.request_protocol"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "forwarding_rule.0.request_port"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "forwarding_rule.0.communication_protocol"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "forwarding_rule.0.communication_port"),

					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.#", "1"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "health_check.0.id"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "health_check.0.type"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "health_check.0.port"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "health_check.0.interval"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "health_check.0.timeout"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "health_check.0.unhealthy_threshold"),
					resource.TestCheckResourceAttrSet("ah_load_balancer.web", "health_check.0.healthy_threshold"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_DeleteForwardingRule(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_WithFR(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "forwarding_rule.#", "2"),
				),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_Empty(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "forwarding_rule.#", "0"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_Rename(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))
	newName := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_Empty(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "name", name),
				),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_Empty(newName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "name", newName),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_ChangeBalancingAlgorithm(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_Empty(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "balancing_algorithm", "round_robin"),
				),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_LeastRequests(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "balancing_algorithm", "least_requests"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_AddForwardingRule(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_Empty(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "forwarding_rule.#", "0"),
				),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_WithFR(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "forwarding_rule.#", "2"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_UpdateForwardingRule(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_WithFR(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("ah_load_balancer.web", "forwarding_rule.*", map[string]string{
						"communication_port": "80",
					}),
				),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_UpdateFR(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "forwarding_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("ah_load_balancer.web", "forwarding_rule.*", map[string]string{
						"communication_port": "9090",
					}),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_ConnectPrivateNetwork(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_WithPN(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "private_network.#", "1"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_ConnectBackendNode(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHLoadBalancerConfig_WithPN(name),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHLoadBalancerConfig_WithBackendNode(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "backend_node.#", "1"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_RemoveBackendNode(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHLoadBalancerConfig_WithBackendNode(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "backend_node.#", "1"),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHLoadBalancerConfig_WithoutBackendNode(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "backend_node.#", "0"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_ChangeBackendNode(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfigBasic() + testAccCheckAHLoadBalancerConfig_WithBackendNode(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "backend_node.#", "1"),
				),
			},
			{
				Config: datasourceConfigBasic() + testAccCheckAHLoadBalancerConfig_ChangeBackendNode(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "backend_node.#", "1"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_CreateHealthCheck(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_Empty(name),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_WithHealthCheck(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.#", "1"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_DeleteHealthCheck(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_WithHealthCheck(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.#", "1"),
				),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_Empty(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.#", "0"),
				),
			},
		},
	})
}

func TestAccAHLoadBalancer_UpdateHealthCheck(t *testing.T) {
	name := fmt.Sprintf("test-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAHLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAHLoadBalancerConfig_WithHealthCheck(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.#", "1"),
				),
			},
			{
				Config: testAccCheckAHLoadBalancerConfig_UpdateHealthCheck(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.#", "1"),
					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.0.port", "9091"),
					resource.TestCheckResourceAttr("ah_load_balancer.web", "health_check.0.unhealthy_threshold", "3"),
				),
			},
		},
	})
}

func testAccCheckAHLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ah.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ah_load_balancer" {
			continue
		}

		_, err := client.LoadBalancers.Get(context.Background(), rs.Primary.ID)

		if err != ah.ErrResourceNotFound {
			return fmt.Errorf(
				"error waiting for load balancer (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckAHLoadBalancerConfig_Basic(name string) string {
	return fmt.Sprintf(`
     resource "ah_private_network" "test" {
	   ip_range = "10.0.0.0/24"
	   name = "Test Private Network"
	 }
	 
	 resource "ah_cloud_server" "web" {
       count = 2
	   name = "cs-%[1]s-${count.index}"
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   product = "%s"
	 }

	 resource "ah_private_network_connection" "example" {
       count = 2 
	   cloud_server_id = ah_cloud_server.web[count.index].id
	   private_network_id = ah_private_network.test.id
	 }

	 resource "ah_load_balancer" "web" {
       depends_on = [
         ah_private_network_connection.example,
	   ]
	   name = "%[1]s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       private_network {
         id = ah_private_network.test.id
       }
       backend_node {
         cloud_server_id = ah_cloud_server.web[0].id
       }
       backend_node {
         cloud_server_id = ah_cloud_server.web[1].id
       }
       forwarding_rule {
         request_protocol = "tcp"
         request_port = 80
         communication_protocol = "tcp"
         communication_port = 80
       }
       forwarding_rule {
         request_protocol = "tcp"
         request_port = 8080
         communication_protocol = "tcp"
         communication_port = 8080
       }
       health_check {
         type = "tcp"
         port = "9090"
       }
	 }`, name, DatacenterName, VpsPlanID, DatacenterName)
}

func testAccCheckAHLoadBalancerConfig_Empty(name string) string {
	return fmt.Sprintf(`
	 resource "ah_load_balancer" "web" {
	   name = "%s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
	 }`, name, DatacenterName)
}

func testAccCheckAHLoadBalancerConfig_LeastRequests(name string) string {
	return fmt.Sprintf(`
	 resource "ah_load_balancer" "web" {
	   name = "%s"
	   datacenter = "%s"
	   balancing_algorithm = "least_requests"
	   instance_count = 1
       create_public_ip_address = false
	 }`, name, DatacenterName)
}

func testAccCheckAHLoadBalancerConfig_WithFR(name string) string {
	return fmt.Sprintf(`
	 resource "ah_load_balancer" "web" {
	   name = "%s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       forwarding_rule {
         request_protocol = "tcp"
         request_port = 80
         communication_protocol = "tcp"
         communication_port = 80
       }
       forwarding_rule {
         request_protocol = "tcp"
         request_port = 8080
         communication_protocol = "tcp"
         communication_port = 8080
       }
	 }`, name, DatacenterName)
}

func testAccCheckAHLoadBalancerConfig_UpdateFR(name string) string {
	return fmt.Sprintf(`
	 resource "ah_load_balancer" "web" {
	   name = "%s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       forwarding_rule {
         request_protocol = "tcp"
         request_port = 80
         communication_protocol = "tcp"
         communication_port = 9090
       }
       forwarding_rule {
         request_protocol = "tcp"
         request_port = 8080
         communication_protocol = "tcp"
         communication_port = 8080
       }
	 }`, name, DatacenterName)
}

func testAccCheckAHLoadBalancerConfig_WithPN(name string) string {
	return fmt.Sprintf(`
     resource "ah_private_network" "test" {
	   ip_range = "10.0.0.0/24"
	   name = "Test Private Network"
	 }

	 resource "ah_load_balancer" "web" {
	   name = "%s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       private_network {
         id = ah_private_network.test.id
       }
	 }`, name, DatacenterName)
}

func testAccCheckAHLoadBalancerConfig_WithBackendNode(name string) string {
	return fmt.Sprintf(`
     resource "ah_cloud_server" "web" {
	   count = 2
	   name = "cs-%[1]s-${count.index}"
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   product = "%s"
	 }
     
     resource "ah_private_network" "test" {
	   ip_range = "10.0.0.0/24"
	   name = "Test Private Network"
	 }

     resource "ah_private_network_connection" "example" {
       count = 2 
	   cloud_server_id = ah_cloud_server.web[count.index].id
	   private_network_id = ah_private_network.test.id
	 }

	 resource "ah_load_balancer" "web" {
       depends_on = [
         ah_private_network_connection.example,
	   ]
	   name = "%[1]s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       private_network {
         id = ah_private_network.test.id
       }
       backend_node {
         cloud_server_id = ah_cloud_server.web[0].id
       }
	 }`, name, DatacenterName, VpsPlanID, DatacenterName)

}

func testAccCheckAHLoadBalancerConfig_WithoutBackendNode(name string) string {
	return fmt.Sprintf(`
     resource "ah_cloud_server" "web" {
	   count = 2
	   name = "cs-%[1]s-${count.index}"
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   product = "%s"
	 }
     
     resource "ah_private_network" "test" {
	   ip_range = "10.0.0.0/24"
	   name = "Test Private Network"
	 }

     resource "ah_private_network_connection" "example" {
       count = 2 
	   cloud_server_id = ah_cloud_server.web[count.index].id
	   private_network_id = ah_private_network.test.id
	 }

	 resource "ah_load_balancer" "web" {
       depends_on = [
         ah_private_network_connection.example,
	   ]
	   name = "%[1]s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       private_network {
         id = ah_private_network.test.id
       }
	 }`, name, DatacenterName, VpsPlanID, DatacenterName)

}

func testAccCheckAHLoadBalancerConfig_ChangeBackendNode(name string) string {
	return fmt.Sprintf(`
     resource "ah_cloud_server" "web" {
	   count = 2
	   name = "cs-%[1]s-${count.index}"
	   datacenter = "%s"
	   image = "${data.ah_cloud_images.test.images.0.id}"
	   product = "%s"
	 }
     
     resource "ah_private_network" "test" {
	   ip_range = "10.0.0.0/24"
	   name = "Test Private Network"
	 }

     resource "ah_private_network_connection" "example" {
       count = 2 
	   cloud_server_id = ah_cloud_server.web[count.index].id
	   private_network_id = ah_private_network.test.id
	 }

	 resource "ah_load_balancer" "web" {
       depends_on = [
         ah_private_network_connection.example,
	   ]
	   name = "%[1]s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       private_network {
         id = ah_private_network.test.id
       }
       backend_node {
         cloud_server_id = ah_cloud_server.web[1].id
       }
	 }`, name, DatacenterName, VpsPlanID, DatacenterName)

}

func testAccCheckAHLoadBalancerConfig_WithHealthCheck(name string) string {
	return fmt.Sprintf(`
	 resource "ah_load_balancer" "web" {
	   name = "%s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       health_check {
         type = "tcp"
         port = "9090"
       }
	 }`, name, DatacenterName)
}

func testAccCheckAHLoadBalancerConfig_UpdateHealthCheck(name string) string {
	return fmt.Sprintf(`
	 resource "ah_load_balancer" "web" {
	   name = "%s"
	   datacenter = "%s"
	   balancing_algorithm = "round_robin"
	   instance_count = 1
       create_public_ip_address = false
       health_check {
         type = "tcp"
         port = "9091"
         unhealthy_threshold = 3
       }
	 }`, name, DatacenterName)
}
