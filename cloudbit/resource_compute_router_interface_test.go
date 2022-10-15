package cloudbit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeRouterInterface_Basic(t *testing.T) {
	networkName := acctest.RandomWithPrefix("test-network")
	networkCIDR := "192.168.1.0/24"
	routerName := acctest.RandomWithPrefix("test-router")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccComputeRouterInterfaceConfigBasic, networkName, networkCIDR, routerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cloudbit_compute_router_interface.foobar", "id"),
					resource.TestCheckResourceAttrSet("cloudbit_compute_router_interface.foobar", "router_id"),
					resource.TestCheckResourceAttrSet("cloudbit_compute_router_interface.foobar", "network_id"),
					resource.TestCheckResourceAttrSet("cloudbit_compute_router_interface.foobar", "private_ip"),
				),
			},
		},
	})
}

const testAccComputeRouterInterfaceConfigBasic = `
locals {
	location_id = 1
}

resource "cloudbit_compute_network" "foobar" {
	name        = "%s"
	location_id = local.location_id

	cidr = "%s"
}

resource "cloudbit_compute_router" "foobar" {
	name        = "%s"
	location_id = local.location_id

	public = false
}

resource "cloudbit_compute_router_interface" "foobar" {
	router_id = cloudbit_compute_router.foobar.id
	network_id = cloudbit_compute_network.foobar.id
}
`
