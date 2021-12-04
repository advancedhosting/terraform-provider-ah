package ah

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sort"
	"strconv"
	"strings"
)

func filterPlans(d *schema.ResourceData, plans []map[string]interface{}) []map[string]interface{} {
	filters := buildPlanFilters(d)
	if filters == nil {
		return plans
	}

	var filteredRecords []map[string]interface{}

	for _, plan := range plans {
		if checkPlanFilter(filters, plan) {
			filteredRecords = append(filteredRecords, plan)
		}
	}
	return filteredRecords
}

func buildPlanFilters(d *schema.ResourceData) map[string]map[string]bool {
	buildFunc := func(set *schema.Set) map[string]map[string]bool {
		var filters = make(map[string]map[string]bool, len(set.List()))
		for _, v := range set.List() {
			m := v.(map[string]interface{})
			values := m["values"].([]interface{})
			var filterValues = make(map[string]bool, len(values))
			for _, e := range values {
				filterValues[e.(string)] = true
			}

			key := m["key"].(string)

			filters[key] = filterValues

		}
		return filters
	}

	if v, ok := d.GetOk("filter"); ok {
		return buildFunc(v.(*schema.Set))
	}
	return nil
}

func checkPlanFilter(filters map[string]map[string]bool, plan map[string]interface{}) bool {
	for key, values := range filters {
		v := plan[key]
		if !filterPlanValue(v, values) {
			return false
		}
	}

	return true
}

func filterPlanValue(v interface{}, values map[string]bool) bool {
	switch t := v.(type) {
	case int:
		_, ok := values[strconv.Itoa(t)]
		return ok
	case string:
		_, ok := values[t]
		return ok
	case bool:
		var vBool string
		if t {
			vBool = "true"
		} else {
			vBool = "false"
		}
		_, ok := values[vBool]
		return ok
	default:
		panic("type is not supported")
	}
}

type PlanSorting struct {
	key       string
	direction string
}

func buildPlanSort(d *schema.ResourceData) []PlanSorting {
	buildFunc := func(set *schema.Set) []PlanSorting {
		var sortArr = make([]PlanSorting, len(set.List()))
		for i, v := range set.List() {
			m := v.(map[string]interface{})
			sortArr[i] = PlanSorting{m["key"].(string), m["direction"].(string)}
		}
		return sortArr
	}

	if v, ok := d.GetOk("sort"); ok {
		return buildFunc(v.(*schema.Set))
	}
	return nil
}

func sortPlans(d *schema.ResourceData, plans []map[string]interface{}) {
	sortArr := buildPlanSort(d)
	if sortArr == nil {
		return
	}

	sort.Slice(plans, func(i, j int) bool {
		for _, s := range sortArr {
			if s.direction == "desc" {
				i, j = j, i
			}

			cmp := comparePlanValues(plans[i][s.key], plans[j][s.key])
			if cmp != 0 {
				return cmp < 0
			}
		}

		return true
	})
}

func comparePlanValues(v1, v2 interface{}) int {
	switch t := v1.(type) {
	case int:
		v := v2.(int)
		if t == v {
			return 0
		} else if t < v {
			return -1
		} else {
			return 1
		}
	case string:
		return strings.Compare(t, v2.(string))
	case bool:
		if t == v2.(bool) {
			return 0
		} else if !t {
			return -1
		} else {
			return 1
		}
	default:
		panic("type is not supported")
	}
}
