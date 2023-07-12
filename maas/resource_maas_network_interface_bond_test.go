package maas_test

import (
	"terraform-provider-maas/maas/testutils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestACCResourceMaasNetworkInterfaceBond(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testutils.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestACCResourceMaasNetworkInterfaceBondConfig_basic,
			},
			{
				Config: TestACCResourceMaasNetworkInterfaceBondConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("maas_network_interface_bond.test", "bond_updelay", "11"),
				),
			},
		},
	})
}

const TestACCResourceMaasNetworkInterfaceBondConfig_basic = `
resource "maas_network_interface_bond" "test" {
	machine = "mq4s3r"
	name = "bond0"
	accept_ra = false
	bond_downdelay = 1
	bond_lacp_rate = "slow"
	bond_miimon = 10
	bond_mode = "802.3ad"
	bond_num_grat_arp = 1
	bond_updelay = 1
	bond_xmit_hash_policy = "layer2"
	mtu = 9000
	parents = ["enp109s0f0", "enp109s0f1"]
  }
  `

const TestACCResourceMaasNetworkInterfaceBondConfig_update = `
resource "maas_network_interface_bond" "test" {
	machine = "mq4s3r"
	name = "bond0"
	accept_ra = false
	bond_downdelay = 2
	bond_lacp_rate = "slow"
	bond_miimon = 11
	bond_mode = "802.3ad"
	bond_num_grat_arp = 11
	bond_updelay = 11
	bond_xmit_hash_policy = "layer2"
	mtu = 9000
	parents = ["enp109s0f0", "enp109s0f1"]
  }
  `