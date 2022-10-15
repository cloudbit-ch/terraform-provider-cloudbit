package cloudbit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeVolume_Basic(t *testing.T) {
	volumeName := acctest.RandomWithPrefix("test-volume")
	volumeSize := acctest.RandIntRange(1, 20)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccComputeVolumeConfigBasic, volumeName, volumeSize),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cloudbit_compute_volume.foobar", "id"),
					resource.TestCheckResourceAttr("cloudbit_compute_volume.foobar", "name", volumeName),
					resource.TestCheckResourceAttr("cloudbit_compute_volume.foobar", "location_id", "1"),
					resource.TestCheckResourceAttr("cloudbit_compute_volume.foobar", "size", fmt.Sprint(volumeSize)),
					resource.TestCheckResourceAttrSet("cloudbit_compute_volume.foobar", "serial_number"),
					resource.TestCheckNoResourceAttr("cloudbit_compute_volume.foobar", "restore_from_snapshot_id"),
				),
			},
		},
	})
}

const testAccComputeVolumeConfigBasic = `
resource "cloudbit_compute_volume" "foobar" {
	name        = "%s"
	location_id = 1

	size = %d
}
`
