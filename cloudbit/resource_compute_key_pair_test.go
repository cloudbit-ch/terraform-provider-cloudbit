package cloudbit

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeKeyPair_Basic(t *testing.T) {
	keyPairName := acctest.RandomWithPrefix("test-key-pair")
	public, _, err := acctest.RandSSHKeyPair("test-key-pair")
	if err != nil {
		t.Fatal(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccComputeKeyPairConfigBasic, keyPairName, public),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cloudbit_compute_key_pair.foobar", "id"),
					resource.TestCheckResourceAttrSet("cloudbit_compute_key_pair.foobar", "fingerprint"),
					resource.TestCheckResourceAttr("cloudbit_compute_key_pair.foobar", "name", keyPairName),
					resource.TestCheckResourceAttr("cloudbit_compute_key_pair.foobar", "public_key", public),
				),
			},
		},
	})
}

const testAccComputeKeyPairConfigBasic = `
resource "cloudbit_compute_key_pair" "foobar" {
	name        = "%s"
	public_key  = "%s"
}
`
