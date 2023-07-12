package maas

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maas/gomaasclient/client"
	"github.com/maas/gomaasclient/entity"
)

func resourceMaasNetworkInterfaceBridge() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource to manage MAAS network Bridges.",
		CreateContext: resourceMaasNetworkInterfaceBridgeCreate,
		ReadContext:   resourceMaasNetworkInterfaceBridgeRead,
		UpdateContext: resourceMaasNetworkInterfaceBridgeUpdate,
		DeleteContext: resourceMaasNetworkInterfaceBridgeDelete,

		Schema: map[string]*schema.Schema{
			"machine": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "List of MAAS machines' identifiers (system ID, hostname, or FQDN) that will be tagged with the new tag.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the interface.",
			},
			"parent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parent interface name for this bridge interface.",
			},
			"accept_ra": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Accept router advertisements. (IPv6 only).",
			},
			"bridge_fd": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Set bridge forward delay to time seconds. (Default: 15).",
			},
			"bridge_stp": {
				Type:     schema.TypeBool,
				Optional: true,
				// Computed:    true,
				Description: "Turn spanning tree protocol on or off. (Default: False).",
			},
			"bridge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of bridge to create. Possible values are: ``standard``, ``ovs``.",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "MAC address of the interface.",
			},
			"mtu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum transmission unit.",
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Tags for the interface.",
			},
			"vlan": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "VLAN the interface is connected to.",
			},
		},
	}
}

func resourceMaasNetworkInterfaceBridgeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	machine, err := getMachine(client, d.Get("machine").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	parentID, err := findInterfaceParent(client, machine.SystemID, d.Get("parent").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	params := getNetworkInterfaceBridgeParams(d, parentID)
	networkInterface, err := client.NetworkInterfaces.CreateBridge(machine.SystemID, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v", networkInterface.ID))

	return resourceMaasNetworkInterfaceBridgeRead(ctx, d, m)
}

func resourceMaasNetworkInterfaceBridgeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	machine, err := getMachine(client, d.Get("machine").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	networkInterface, err := client.NetworkInterface.Get(machine.SystemID, id)
	if err != nil {
		return diag.FromErr(err)
	}

	p := networkInterface.Params.(map[string]interface{})
	if _, ok := p["bridge_fd"]; ok {
		d.Set("bridge_fd", int64(p["bridge_fd"].(float64)))
	}
	if _, ok := p["bridge_stp"]; ok {
		d.Set("bridge_stp", p["bridge_stp"].(bool))
	}
	if _, ok := p["bridge_type"]; ok {
		d.Set("bridge_type", p["bridge_type"].(string))
	}
	if _, ok := p["accept-ra"]; ok {
		d.Set("accept_ra", p["accept-ra"].(bool))
	}

	tfState := map[string]interface{}{
		"name":        networkInterface.Name,
		"mac_address": networkInterface.MACAddress,
		"mtu":         networkInterface.EffectiveMTU,
		"tags":        networkInterface.Tags,
		"vlan":        fmt.Sprintf("%v", networkInterface.VLAN.ID),
	}
	if err := setTerraformState(d, tfState); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMaasNetworkInterfaceBridgeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	machine, err := getMachine(client, d.Get("machine").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	parentID, err := findInterfaceParent(client, machine.SystemID, d.Get("parent").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	params := getNetworkInterfaceBridgeUpdateParams(d, parentID)
	_, err = client.NetworkInterface.Update(machine.SystemID, id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMaasNetworkInterfaceBridgeRead(ctx, d, m)
}

func resourceMaasNetworkInterfaceBridgeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	machine, err := getMachine(client, d.Get("machine").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := client.NetworkInterface.Delete(machine.SystemID, id); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getNetworkInterfaceBridgeParams(d *schema.ResourceData, parentID int) *entity.NetworkInterfaceBridgeParams {
	return &entity.NetworkInterfaceBridgeParams{
		MACAddress: d.Get("mac_address").(string),
		Name:       d.Get("name").(string),
		Tags:       strings.Join(convertToStringSlice(d.Get("tags").(*schema.Set).List()), ","),
		VLAN:       d.Get("vlan").(int),
		MTU:        d.Get("mtu").(int),
		AcceptRA:   d.Get("accept_ra").(bool),
		Parents:    []int{parentID},
		BridgeType: d.Get("bridge_type").(string),
		BridgeSTP:  d.Get("bridge_stp").(bool),
		BridgeFD:   d.Get("bridge_fd").(int),
	}
}

func getNetworkInterfaceBridgeUpdateParams(d *schema.ResourceData, parentID int) *entity.NetworkInterfaceUpdateParams {
	return &entity.NetworkInterfaceUpdateParams{
		MACAddress: d.Get("mac_address").(string),
		Name:       d.Get("name").(string),
		Tags:       strings.Join(convertToStringSlice(d.Get("tags").(*schema.Set).List()), ","),
		VLAN:       d.Get("vlan").(int),
		MTU:        d.Get("mtu").(int),
		AcceptRA:   d.Get("accept_ra").(bool),
		Parents:    []int{parentID},
		BridgeType: d.Get("bridge_type").(string),
		BridgeSTP:  d.Get("bridge_stp").(bool),
		BridgeFD:   d.Get("bridge_fd").(int),
	}
}

func findInterfaceParent(client *client.Client, machineSystemID string, parent string) (int, error) {
	networkInterface, err := getNetworkInterface(client, machineSystemID, parent)
	if err != nil {
		return 0, err
	}

	return networkInterface.ID, nil
}