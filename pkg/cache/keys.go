package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// StringSliceKey returns a hash string of a string slice.
// It alse sorts the slice before hashing.
// This is useful for caching long string arrays
func StringSliceKey(slice []string) string {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	joinedString := strings.Join(slice, ",")
	hasher := sha256.New()
	hasher.Write([]byte(joinedString))
	hashSum := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	return hashString
}
