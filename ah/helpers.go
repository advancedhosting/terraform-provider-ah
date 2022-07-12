package ah

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/google/uuid"
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
