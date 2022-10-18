package cloudbit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeRouterRoute_Basic(t *testing.T) {
	networkName := acctest.RandomWithPrefix("test-network")
	networkCIDR := "192.168.1.0/24"
	routerName := acctest.RandomWithPrefix("test-router")

	destination := "0.0.0.0/0"
	nextHop, err := acctest.RandIpAddress(networkCIDR)
	if err != nil {
		t.Fatal(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccComputeRouterRouteConfigBasic, networkName, networkCIDR, routerName, destination, nextHop),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cloudbit_compute_router_route.foobar", "id"),
					resource.TestCheckResourceAttrSet("cloudbit_compute_router_route.foobar", "router_id"),
					resource.TestCheckResourceAttr("cloudbit_compute_router_route.foobar", "destination", destination),
					resource.TestCheckResourceAttr("cloudbit_compute_router_route.foobar", "next_hop", nextHop),
				),
			},
		},
	})
}

const testAccComputeRouterRouteConfigBasic = `
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

resource "cloudbit_compute_router_route" "foobar" {
	router_id = cloudbit_compute_router.foobar.id
	destination = "%s"
	next_hop = "%s"

	depends_on = [cloudbit_compute_router_interface.foobar]
}
`
