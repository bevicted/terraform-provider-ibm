// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworksRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Networks: {
				Computed:    true,
				Description: "List of all networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_AccessConfig: {
							Computed:    true,
							Description: "The network communication configuration option of the network (for satellite locations only).",
							Type:        schema.TypeString,
						},
						Attr_DhcpManaged: {
							Computed:    true,
							Description: "Indicates if the network DHCP Managed.",
							Type:        schema.TypeBool,
						},
						Attr_Href: {
							Computed:    true,
							Description: "The hyper link of a network.",
							Type:        schema.TypeString,
						},
						Attr_MTU: {
							Computed:    true,
							Description: "Maximum Transmission Unit option of the network.",
							Type:        schema.TypeInt,
						},
						Attr_Name: {
							Computed:    true,
							Description: "The name of a network.",
							Type:        schema.TypeString,
						},
						Attr_NetworkID: {
							Computed:    true,
							Description: "The unique identifier of a network.",
							Type:        schema.TypeString,
						},
						Attr_Type: {
							Computed:    true,
							Description: "The type of network.",
							Type:        schema.TypeString,
						},
						Attr_VLanID: {
							Computed:    true,
							Description: "The VLAN ID that the network is connected to.",
							Type:        schema.TypeInt,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPINetworksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	networkC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_Networks, flattenNetworks(networkdata.Networks))

	return nil
}

func flattenNetworks(list []*models.NetworkReference) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			Attr_AccessConfig: i.AccessConfig,
			Attr_DhcpManaged:  i.DhcpManaged,
			Attr_Href:         *i.Href,
			Attr_MTU:          i.Mtu,
			Attr_Name:         *i.Name,
			Attr_NetworkID:    *i.NetworkID,
			Attr_Type:         *i.Type,
			Attr_VLanID:       *i.VlanID,
		}
		result = append(result, l)
	}
	return result
}
