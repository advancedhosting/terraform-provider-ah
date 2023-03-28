package ah

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAHLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAHLoadBalancerCreate,
		ReadContext:   resourceAHLoadBalancerRead,
		UpdateContext: resourceAHLoadBalancerUpdate,
		DeleteContext: resourceAHLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"datacenter": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"create_public_ip_address": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"balancing_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ip_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"private_network": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"backend_node": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_server_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"forwarding_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"request_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"communication_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"communication_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAHLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	request := &ah.LoadBalancerCreateRequest{
		Name:                  d.Get("name").(string),
		CreatePublicIPAddress: d.Get("create_public_ip_address").(bool),
	}

	datacenterAttr := d.Get("datacenter").(string)
	if _, err := uuid.Parse(datacenterAttr); err != nil {
		datacenterID, err := datacenterIDBySlug(ctx, client, d.Get("datacenter").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		request.DatacenterID = datacenterID
	} else {
		request.DatacenterID = datacenterAttr
	}

	if attr, ok := d.GetOk("ip_address"); ok {
		var ipAddressIDs []string
		for _, v := range attr.(*schema.Set).List() {
			ipAddress := v.(map[string]interface{})
			ipAddressIDs = append(ipAddressIDs, ipAddress["id"].(string))
		}
		request.IPAddressIDs = ipAddressIDs
	}

	if attr, ok := d.GetOk("private_network"); ok {
		var pnIDs []string
		for _, v := range attr.(*schema.Set).List() {
			pn := v.(map[string]interface{})
			pnIDs = append(pnIDs, pn["id"].(string))
		}
		request.PrivateNetworkIDs = pnIDs
	}

	if attr, ok := d.GetOk("backend_node"); ok {
		var bNodesRequest []ah.LBBackendNodeCreateRequest
		for _, v := range attr.(*schema.Set).List() {
			bNode := v.(map[string]interface{})
			bNodeRequest := ah.LBBackendNodeCreateRequest{
				CloudServerID: bNode["cloud_server_id"].(string),
			}
			bNodesRequest = append(bNodesRequest, bNodeRequest)
		}
		request.BackendNodes = bNodesRequest
	}

	if attr, ok := d.GetOk("forwarding_rule"); ok {
		var frsRequest []ah.LBForwardingRuleCreateRequest
		for _, v := range attr.(*schema.Set).List() {
			fr := v.(map[string]interface{})
			frRequest := makeFRCreateRequest(fr)
			frsRequest = append(frsRequest, *frRequest)
		}
		request.ForwardingRules = frsRequest
	}

	if attr, ok := d.GetOk("health_check"); ok {
		var hcsRequest []ah.LBHealthCheckCreateRequest
		for _, v := range attr.([]interface{}) {
			hc := v.(map[string]interface{})
			hcRequest := makeHCCreateRequest(hc)
			hcsRequest = append(hcsRequest, hcRequest)
		}
		request.HealthChecks = hcsRequest
	}

	lb, err := client.LoadBalancers.Create(ctx, request)

	if err != nil {
		return diag.Errorf("Error creating load balancer: %s", err)
	}
	d.SetId(lb.ID)
	if err := waitForLoadBalancerStatus(ctx, []string{"creating"}, []string{"active"}, d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
	}
	return resourceAHLoadBalancerRead(ctx, d, meta)
}

func resourceAHLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	loadBalancer, err := client.LoadBalancers.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", loadBalancer.Name)
	d.Set("state", loadBalancer.State)

	ipsAddresses := make([]map[string]interface{}, len(loadBalancer.IPAddresses))
	for i, ipAddress := range loadBalancer.IPAddresses {
		item := make(map[string]interface{})
		item["id"] = ipAddress.ID
		item["type"] = ipAddress.Type
		item["address"] = ipAddress.Address
		item["state"] = ipAddress.State
		ipsAddresses[i] = item

	}
	d.Set("ip_address", ipsAddresses)

	privateNetworks := make([]map[string]interface{}, len(loadBalancer.PrivateNetworks))
	for i, pn := range loadBalancer.PrivateNetworks {
		item := make(map[string]interface{})
		item["id"] = pn.ID
		item["state"] = pn.State
		privateNetworks[i] = item

	}
	d.Set("private_network", privateNetworks)

	backendNodes := make([]map[string]interface{}, len(loadBalancer.BackendNodes))
	for i, bn := range loadBalancer.BackendNodes {
		item := make(map[string]interface{})
		item["id"] = bn.ID
		item["cloud_server_id"] = bn.CloudServerID
		backendNodes[i] = item

	}
	d.Set("backend_node", backendNodes)

	forwardingRules := make([]map[string]interface{}, len(loadBalancer.ForwardingRules))
	for i, fr := range loadBalancer.ForwardingRules {
		item := make(map[string]interface{})
		item["id"] = fr.ID
		item["request_protocol"] = fr.RequestProtocol
		item["request_port"] = fr.RequestPort
		item["communication_protocol"] = fr.CommunicationProtocol
		item["communication_port"] = fr.CommunicationPort
		forwardingRules[i] = item

	}
	d.Set("forwarding_rule", forwardingRules)
	healthChecks := make([]map[string]interface{}, len(loadBalancer.HealthChecks))
	for i, hc := range loadBalancer.HealthChecks {
		item := make(map[string]interface{})
		item["id"] = hc.ID
		item["type"] = hc.Type
		if hc.URL != "" {
			item["url"] = hc.URL
		}
		item["interval"] = hc.Interval
		item["timeout"] = hc.Timeout
		item["unhealthy_threshold"] = hc.UnhealthyThreshold
		item["healthy_threshold"] = hc.HealthyThreshold
		item["port"] = hc.Port

		healthChecks[i] = item

	}
	d.Set("health_check", healthChecks)

	return nil

}

func resourceAHLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)

	if d.HasChange("name") {
		request := &ah.LoadBalancerUpdateRequest{
			Name: d.Get("name").(string),
		}

		err := client.LoadBalancers.Update(ctx, d.Id(), request)

		if err != nil {
			return diag.Errorf(
				"Error renaming load balancer (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("balancing_algorithm") {
		request := &ah.LoadBalancerUpdateRequest{
			BalancingAlgorithm: d.Get("balancing_algorithm").(string),
		}

		err := client.LoadBalancers.Update(ctx, d.Id(), request)

		if err != nil {
			return diag.Errorf(
				"Error changing load balancer balancing_algorithm (%s): %s", d.Id(), err)
		}

		if err := waitForLoadBalancerStatus(ctx, []string{"updating"}, []string{"active"}, d, meta); err != nil {
			return diag.Errorf(
				"Error waiting for load balancer (%s) to become ready: %s", d.Id(), err)
		}

	}

	if d.HasChange("forwarding_rule") {
		err := updateForwardingRules(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("private_network") {
		err := updatePrivateNetworks(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backend_node") {
		err := updateBackendNodes(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("health_check") {
		err := updateHealthChecks(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAHLoadBalancerRead(ctx, d, meta)
}

func resourceAHLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ah.APIClient)
	if err := client.LoadBalancers.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf(
			"Error deleting load balancer (%s): %s", d.Id(), err)
	}

	if err := waitForLoadBalancerDestroy(ctx, d, meta); err != nil {
		return diag.Errorf(
			"Error waiting for load balancer (%s) to become deleted: %s", d.Id(), err)
	}

	return nil
}

func waitForLoadBalancerStatus(ctx context.Context, pendingStatuses, targetStatuses []string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		lb, err := client.LoadBalancers.Get(context.Background(), d.Id())
		if err != nil {
			log.Printf("Error on waitForLoadBalancerStatus: %v", err)
			return nil, "", err
		}
		return lb.ID, lb.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      20 * time.Second,
		Pending:    pendingStatuses,
		Refresh:    stateRefreshFunc,
		Target:     targetStatuses,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 5 * time.Second,
	}
	_, err := stateChangeConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf(
			"error waiting for load balancer to reach desired status %s: %s", targetStatuses, err)
	}

	return nil
}

func waitForLoadBalancerDestroy(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ah.APIClient)

	stateRefreshFunc := func() (interface{}, string, error) {
		lb, err := client.LoadBalancers.Get(context.Background(), d.Id())
		if err == ah.ErrResourceNotFound {
			return d.Id(), "deleted", nil
		}
		if err != nil {
			log.Printf("Error on waitForLoadBalancerDestroy: %v", err)
			return nil, "", err
		}

		return lb.ID, lb.State, nil
	}

	stateChangeConf := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"active", "deleting"},
		Refresh:    stateRefreshFunc,
		Target:     []string{"deleted"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf(
			"error waiting for load balancer to reach desired status deleted: %s", err)
	}

	return nil
}

func updateForwardingRules(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	oldFRs, newFRs := d.GetChange("forwarding_rule")
	oldFRsList := oldFRs.(*schema.Set).List()
	newFRsList := newFRs.(*schema.Set).List()

	frsToDelete := make(map[int]interface{}, len(oldFRsList))

	for _, v := range oldFRsList {
		fr := v.(map[string]interface{})
		frsToDelete[fr["request_port"].(int)] = fr
	}

	for _, v := range newFRsList {
		newFR := v.(map[string]interface{})
		fr, ok := frsToDelete[newFR["request_port"].(int)]
		if !ok {
			if err := addForwardingRule(ctx, d, meta, newFR); err != nil {
				return err
			}
		} else {

			oldFR := fr.(map[string]interface{})

			if oldFR["request_protocol"].(string) != newFR["request_protocol"].(string) ||
				oldFR["request_port"].(int) != newFR["request_port"].(int) ||
				oldFR["communication_protocol"].(string) != newFR["communication_protocol"].(string) ||
				oldFR["communication_port"].(int) != newFR["communication_port"].(int) {

				if err := updateForwardingRule(ctx, d, meta, oldFR["id"].(string), newFR); err != nil {
					return err
				}

			}

			delete(frsToDelete, newFR["request_port"].(int))
		}

	}

	for _, v := range frsToDelete {
		fr := v.(map[string]interface{})
		if err := removeForwardingRule(ctx, d, meta, fr["id"].(string)); err != nil {
			return err
		}
	}

	return nil
}

func makeFRCreateRequest(fr map[string]interface{}) *ah.LBForwardingRuleCreateRequest {
	return &ah.LBForwardingRuleCreateRequest{
		RequestProtocol:       fr["request_protocol"].(string),
		RequestPort:           fr["request_port"].(int),
		CommunicationProtocol: fr["communication_protocol"].(string),
		CommunicationPort:     fr["communication_port"].(int),
	}
}

func addForwardingRule(ctx context.Context, d *schema.ResourceData, meta interface{}, fr map[string]interface{}) error {
	client := meta.(*ah.APIClient)

	frRequest := makeFRCreateRequest(fr)

	newFR, err := client.LoadBalancers.CreateForwardingRule(ctx, d.Id(), frRequest)
	if err != nil {
		return fmt.Errorf("error creating forwarding rule %s", err)
	}

	stateFunc := func() (result interface{}, state string, err error) {
		fr, err := client.LoadBalancers.GetForwardingRule(ctx, d.Id(), newFR.ID)
		if err != nil {
			return nil, "", err
		}
		return newFR.ID, fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "creating", "active", d); err != nil {
		return err
	}

	return nil
}

func removeForwardingRule(ctx context.Context, d *schema.ResourceData, meta interface{}, frID string) error {
	client := meta.(*ah.APIClient)
	if err := client.LoadBalancers.DeleteForwardingRule(ctx, d.Id(), frID); err != nil {
		return err
	}

	stateFunc := func() (result interface{}, state string, err error) {
		fr, err := client.LoadBalancers.GetForwardingRule(ctx, d.Id(), frID)
		if err != nil {
			if err == ah.ErrResourceNotFound {
				return frID, "deleted", nil
			}
			return nil, "", err
		}
		return frID, fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "deleting", "deleted", d); err != nil {
		return err
	}

	return nil
}

func updateForwardingRule(ctx context.Context, d *schema.ResourceData, meta interface{}, frID string, fr map[string]interface{}) error {
	if err := removeForwardingRule(ctx, d, meta, frID); err != nil {
		return err
	}

	if err := addForwardingRule(ctx, d, meta, fr); err != nil {
		return err
	}
	return nil
}

func waitForState(ctx context.Context, stateFunc resource.StateRefreshFunc, pendingState, expectedState string, d *schema.ResourceData) error {

	stateChangeConf := resource.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{pendingState},
		Refresh:    stateFunc,
		Target:     []string{expectedState},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	_, err := stateChangeConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf(
			"error waiting for state: %s", err)
	}

	return nil
}

func updatePrivateNetworks(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	oldPNs, newPNs := d.GetChange("private_network")
	oldPNsList := oldPNs.(*schema.Set).List()
	newPNsList := newPNs.(*schema.Set).List()

	pnsToDelete := make(map[string]bool, len(oldPNsList))

	for _, v := range oldPNsList {
		fr := v.(map[string]interface{})
		pnsToDelete[fr["id"].(string)] = true
	}

	for _, v := range newPNsList {
		newPN := v.(map[string]interface{})
		pnID := newPN["id"].(string)
		_, ok := pnsToDelete[pnID]
		if !ok {
			if err := addPrivateNetwork(ctx, d, meta, pnID); err != nil {
				return err
			}
		} else {
			delete(pnsToDelete, pnID)
		}

	}

	for pnID := range pnsToDelete {
		if err := removePrivateNetwork(ctx, d, meta, pnID); err != nil {
			return err
		}
	}

	return nil
}

func addPrivateNetwork(ctx context.Context, d *schema.ResourceData, meta interface{}, pnID string) error {
	client := meta.(*ah.APIClient)

	_, err := client.LoadBalancers.ConnectPrivateNetworks(ctx, d.Id(), []string{pnID})
	if err != nil {
		return fmt.Errorf("error connecting private network %s", err)
	}

	stateFunc := func() (result interface{}, state string, err error) {
		pn, err := client.LoadBalancers.GetPrivateNetwork(ctx, d.Id(), pnID)
		if err != nil {
			return nil, "", err
		}
		return pnID, pn.State, nil
	}

	if err := waitForState(ctx, stateFunc, "updating", "active", d); err != nil {
		return err
	}

	return nil
}

func removePrivateNetwork(ctx context.Context, d *schema.ResourceData, meta interface{}, pnID string) error {
	client := meta.(*ah.APIClient)
	if err := client.LoadBalancers.DisconnectPrivateNetwork(ctx, d.Id(), pnID); err != nil {
		return err
	}

	stateFunc := func() (result interface{}, state string, err error) {
		fr, err := client.LoadBalancers.GetPrivateNetwork(ctx, d.Id(), pnID)
		if err != nil {
			if err == ah.ErrResourceNotFound {
				return pnID, "deleted", nil
			}
			return nil, "", err
		}
		return pnID, fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "deleting", "deleted", d); err != nil {
		return err
	}

	return nil
}

func updateBackendNodes(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	oldBNs, newBNs := d.GetChange("backend_node")
	oldBNsList := oldBNs.(*schema.Set).List()
	newBNsList := newBNs.(*schema.Set).List()

	bnsToDelete := make(map[string]bool, len(oldBNsList))

	for _, v := range oldBNsList {
		fr := v.(map[string]interface{})
		bnsToDelete[fr["id"].(string)] = true
	}

	for _, v := range newBNsList {
		newBN := v.(map[string]interface{})
		bnID := newBN["id"].(string)
		_, ok := bnsToDelete[bnID]
		if !ok {
			if err := addBackendNode(ctx, d, meta, newBN); err != nil {
				return err
			}
		} else {
			delete(bnsToDelete, bnID)
		}

	}

	for bnID, _ := range bnsToDelete {
		if err := removeBackendNode(ctx, d, meta, bnID); err != nil {
			return err
		}
	}

	return nil
}

func addBackendNode(ctx context.Context, d *schema.ResourceData, meta interface{}, backendNode map[string]interface{}) error {
	client := meta.(*ah.APIClient)

	bns, err := client.LoadBalancers.AddBackendNodes(ctx, d.Id(), []string{backendNode["cloud_server_id"].(string)})
	if err != nil {
		return fmt.Errorf("error connecting backend node %s", err)
	}

	bnID := bns[0].ID

	stateFunc := func() (result interface{}, state string, err error) {
		bn, err := client.LoadBalancers.GetBackendNode(ctx, d.Id(), bnID)
		if err != nil {
			return nil, "", err
		}
		return bnID, bn.State, nil
	}

	if err := waitForState(ctx, stateFunc, "updating", "active", d); err != nil {
		return err
	}

	return nil
}

func removeBackendNode(ctx context.Context, d *schema.ResourceData, meta interface{}, bnID string) error {
	client := meta.(*ah.APIClient)
	if err := client.LoadBalancers.DeleteBackendNode(ctx, d.Id(), bnID); err != nil {
		return err
	}

	stateFunc := func() (result interface{}, state string, err error) {
		fr, err := client.LoadBalancers.GetBackendNode(ctx, d.Id(), bnID)
		if err != nil {
			if err == ah.ErrResourceNotFound {
				return bnID, "deleted", nil
			}
			return nil, "", err
		}
		return bnID, fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "deleting", "deleted", d); err != nil {
		return err
	}

	return nil
}

func updateHealthChecks(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	oldHCs, newHCs := d.GetChange("health_check")
	oldHCsList := oldHCs.([]interface{})
	newHCsList := newHCs.([]interface{})

	if len(oldHCsList) == 0 {
		hc := newHCsList[0].(map[string]interface{})
		if err := addHealthCheck(ctx, d, meta, hc); err != nil {
			return err
		}
		return nil
	}

	if len(newHCsList) == 0 {
		hc := oldHCsList[0].(map[string]interface{})
		if err := removeHealthCheck(ctx, d, meta, hc["id"].(string)); err != nil {
			return err
		}
		return nil
	}

	oldHC := oldHCsList[0].(map[string]interface{})
	hcID := oldHC["id"].(string)
	hc := newHCsList[0].(map[string]interface{})
	if err := updateHealthCheck(ctx, d, meta, hcID, hc); err != nil {
		return err
	}
	return nil
}

func makeHCCreateRequest(hc map[string]interface{}) ah.LBHealthCheckCreateRequest {
	hcRequest := ah.LBHealthCheckCreateRequest{
		Type: hc["type"].(string),
		Port: hc["port"].(int),
	}

	if url, ok := hc["url"]; ok {
		hcRequest.URL = url.(string)
	}
	if interval, ok := hc["interval"]; ok {
		hcRequest.Interval = interval.(int)
	}
	if timeout, ok := hc["timeout"]; ok {
		hcRequest.Timeout = timeout.(int)
	}
	if unhealthyThreshold, ok := hc["unhealthy_threshold"]; ok {
		hcRequest.UnhealthyThreshold = unhealthyThreshold.(int)
	}
	if healthyThreshold, ok := hc["healthy_threshold"]; ok {
		hcRequest.HealthyThreshold = healthyThreshold.(int)
	}
	return hcRequest
}

func addHealthCheck(ctx context.Context, d *schema.ResourceData, meta interface{}, hc map[string]interface{}) error {
	client := meta.(*ah.APIClient)

	hcRequest := makeHCCreateRequest(hc)

	newHC, err := client.LoadBalancers.CreateHealthCheck(ctx, d.Id(), &hcRequest)
	if err != nil {
		return fmt.Errorf("error creating health check: %s", err)
	}

	stateFunc := func() (result interface{}, state string, err error) {
		fr, err := client.LoadBalancers.GetHealthCheck(ctx, d.Id(), newHC.ID)
		if err != nil {
			return nil, "", err
		}
		return newHC.ID, fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "creating", "active", d); err != nil {
		return err
	}

	return nil
}

func removeHealthCheck(ctx context.Context, d *schema.ResourceData, meta interface{}, hcID string) error {
	client := meta.(*ah.APIClient)
	if err := client.LoadBalancers.DeleteHealthCheck(ctx, d.Id(), hcID); err != nil {
		return err
	}

	stateFunc := func() (result interface{}, state string, err error) {
		hc, err := client.LoadBalancers.GetHealthCheck(ctx, d.Id(), hcID)
		if err != nil {
			if err == ah.ErrResourceNotFound {
				return hcID, "deleted", nil
			}
			return nil, "", err
		}
		return hcID, hc.State, nil
	}

	if err := waitForState(ctx, stateFunc, "deleting", "deleted", d); err != nil {
		return err
	}

	return nil
}

func updateHealthCheck(ctx context.Context, d *schema.ResourceData, meta interface{}, hcID string, hc map[string]interface{}) error {
	client := meta.(*ah.APIClient)

	hcRequest := &ah.LBHealthCheckUpdateRequest{
		Type: hc["type"].(string),
		Port: hc["port"].(int),
	}

	if url, ok := hc["url"]; ok {
		hcRequest.URL = url.(string)
	}
	if interval, ok := hc["interval"]; ok {
		hcRequest.Interval = interval.(int)
	}
	if timeout, ok := hc["timeout"]; ok {
		hcRequest.Timeout = timeout.(int)
	}
	if unhealthyThreshold, ok := hc["unhealthy_threshold"]; ok {
		hcRequest.UnhealthyThreshold = unhealthyThreshold.(int)
	}
	if healthyThreshold, ok := hc["healthy_threshold"]; ok {
		hcRequest.HealthyThreshold = healthyThreshold.(int)
	}

	if err := client.LoadBalancers.UpdateHealthCheck(ctx, d.Id(), hcID, hcRequest); err != nil {
		return fmt.Errorf("error updating health check %s", err)
	}

	stateFunc := func() (result interface{}, state string, err error) {
		fr, err := client.LoadBalancers.GetHealthCheck(ctx, d.Id(), hcID)
		if err != nil {
			return nil, "", err
		}
		return hcID, fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "updating", "active", d); err != nil {
		return err
	}

	return nil
}
