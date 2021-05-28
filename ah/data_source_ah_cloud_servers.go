package ah

import (
	"context"
	"fmt"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAHCloudServers() *schema.Resource {
	allowedFilterKeys := []string{"id", "name", "state", "vcpu", "ram", "disk"}
	allowedSortingKeys := []string{"id", "state", "created_at", "vcpu", "ram", "disk"}
	return &schema.Resource{
		Read: dataSourceAHCloudServersRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(allowedFilterKeys, false),
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"sort": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(allowedSortingKeys, false),
						},
						"direction": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "desc",
							ValidateFunc: validation.StringInSlice([]string{"asc", "desc"}, false),
						},
					},
				},
			},
			"cloud_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backups": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"use_password": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"assignment_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"primary": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"reverse_dns": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"private_networks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildListFilter(set *schema.Set) []ah.FilterInterface {
	var filters []ah.FilterInterface
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		filter := &ah.InFilter{
			Keys:   []string{m["key"].(string)},
			Values: filterValues,
		}

		filters = append(filters, filter)
	}
	return filters
}

func buildListSorting(set *schema.Set) []*ah.Sorting {
	var sortings []*ah.Sorting
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		sorting := &ah.Sorting{
			Key:   m["key"].(string),
			Order: m["direction"].(string),
		}

		sortings = append(sortings, sorting)
	}
	return sortings
}

func dataSourceAHCloudServersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	var listFilters []ah.FilterInterface
	if v, ok := d.GetOk("filter"); ok {
		listFilters = buildListFilter(v.(*schema.Set))
	}

	var listSortings []*ah.Sorting
	if v, ok := d.GetOk("sort"); ok {
		listSortings = buildListSorting(v.(*schema.Set))
	}

	cloudServers, err := allCloudServers(client, listFilters, listSortings)
	if err != nil {
		return err
	}

	if err = cloudServersSchema(cloudServers, d, meta); err != nil {
		return err
	}
	return nil
}

func allCloudServers(client *ah.APIClient, listFilters []ah.FilterInterface, listSortings []*ah.Sorting) ([]ah.Instance, error) {
	options := &ah.ListOptions{
		Meta: &ah.ListMetaOptions{
			Page: 1,
		},
	}

	if len(listFilters) > 0 {
		options.Filters = listFilters
	}

	if len(listSortings) > 0 {
		options.Sortings = listSortings
	}

	var cloudServers []ah.Instance

	for {
		servers, meta, err := client.Instances.List(context.Background(), options)

		if err != nil {
			return nil, fmt.Errorf("Error list instances: %s", err)
		}

		cloudServers = append(cloudServers, servers...)
		if meta.IsLastPage() {
			break
		}

		options.Meta.Page++
	}

	return cloudServers, nil
}

func cloudServersSchema(instances []ah.Instance, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)
	cloudServers := make([]map[string]interface{}, len(instances))
	var ids string
	for i, instance := range instances {
		cloudServer := map[string]interface{}{
			"id":           instance.ID,
			"name":         instance.Name,
			"datacenter":   instance.Datacenter.ID,
			"product":      instance.ProductID,
			"state":        instance.State,
			"vcpu":         instance.Vcpu,
			"ram":          instance.RAM,
			"disk":         instance.Disk,
			"created_at":   instance.CreatedAt,
			"image":        instance.Image.ID,
			"backups":      instance.SnapshotBySchedule,
			"use_password": instance.UseSSHPassword,
		}
		ids += instance.ID

		var privateNetworks []map[string]string
		for _, instancePrivateNetwork := range instance.PrivateNetworks {
			item := make(map[string]string)
			item["id"] = instancePrivateNetwork.PrivateNetwork.ID
			item["ip"] = instancePrivateNetwork.IP
			privateNetworks = append(privateNetworks, item)
		}
		if len(privateNetworks) > 0 {
			cloudServer["private_networks"] = privateNetworks
		}

		var volumesIDs []string
		if len(instance.Volumes) > 0 {
			for _, volume := range instance.Volumes {
				volumesIDs = append(volumesIDs, volume.ID)
			}
		}
		if len(volumesIDs) > 0 {
			cloudServer["volumes"] = volumesIDs
		}

		var ips []map[string]interface{}
		for _, instanceIPAddress := range instance.IPAddresses {
			item := make(map[string]interface{})
			item["assignment_id"] = instanceIPAddress.ID
			item["ip_address"] = instanceIPAddress.Address
			item["primary"] = instance.PrimaryInstanceIPAddressID == instanceIPAddress.ID
			ipAddress, err := client.IPAddresses.Get(context.Background(), instanceIPAddress.IPAddressID)
			if err != nil {
				return err
			}
			item["type"] = ipAddress.Type
			item["reverse_dns"] = ipAddress.ReverseDNS
			ips = append(ips, item)
		}
		if len(ips) > 0 {
			cloudServer["ips"] = ips
		}

		cloudServers[i] = cloudServer

	}

	if err := d.Set("cloud_servers", cloudServers); err != nil {
		return fmt.Errorf("unable to set cloud_servers attribute: %s", err)
	}
	d.SetId(generateHash(ids))

	return nil
}
