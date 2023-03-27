package ah

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/google/uuid"
)

const (
	ImageName      = "centos-7-x64"
	DatacenterID   = "c54e8896-53d8-479a-8ff1-4d7d9d856a50"
	DatacenterName = "ams1"
	VpsPlanID      = "381347529"
	VpsPlanName    = "start-xs"
	VolumePlanID   = "381347560"
	VolumePlanName = "hdd2-ash1"
	ClusterID      = ""
	NodeID         = "2486b2f8-f7a6-4207-979b-9b94d93c174e"
)

// IsUUID checks if string is in valid UUID format.
func IsUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func generateHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func datacenterIDBySlug(ctx context.Context, client *ah.APIClient, datacenterSlug string) (string, error) {
	datacenters, err := client.Datacenters.List(ctx, nil)
	if err != nil {
		return "", err
	}
	for _, datacenter := range datacenters {
		if datacenter.Slug == datacenterSlug {
			return datacenter.ID, nil
		}
	}
	return "", fmt.Errorf("datacenter slug %s not found", datacenterSlug)
}
