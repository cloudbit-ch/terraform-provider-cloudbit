package cloudbit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeRouter_Basic(t *testing.T) {
	routerName := acctest.RandomWithPrefix("test-router")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccComputeRouterConfigBasic, "foobar_public", routerName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cloudbit_compute_router.foobar_public", "id"),
					resource.TestCheckResourceAttr("cloudbit_compute_router.foobar_public", "name", routerName),
					resource.TestCheckResourceAttr("cloudbit_compute_router.foobar_public", "location_id", "1"),
					resource.TestCheckResourceAttr("cloudbit_compute_router.foobar_public", "public", "true"),
					resource.TestCheckResourceAttrSet("cloudbit_compute_router.foobar_public", "public_ip"),
				),
			},
			{
				Config: fmt.Sprintf(testAccComputeRouterConfigBasic, "foobar_private", routerName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cloudbit_compute_router.foobar_private", "id"),
					resource.TestCheckResourceAttr("cloudbit_compute_router.foobar_private", "name", routerName),
					resource.TestCheckResourceAttr("cloudbit_compute_router.foobar_private", "location_id", "1"),
					resource.TestCheckResourceAttr("cloudbit_compute_router.foobar_private", "public", "false"),
					resource.TestCheckNoResourceAttr("cloudbit_compute_router.foobar_private", "public_ip"),
				),
			},
		},
	})
}

const testAccComputeRouterConfigBasic = `
resource "cloudbit_compute_router" "%s" {
	name        = "%s"
	location_id = 1

	public = %t
}
`
