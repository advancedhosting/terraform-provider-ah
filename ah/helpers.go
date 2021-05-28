package ah

import (
	"crypto/sha1"
	"fmt"
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
