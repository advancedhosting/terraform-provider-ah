package ah

import "github.com/google/uuid"

// IsUUID checks if string is in valid UUID format.
func IsUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
