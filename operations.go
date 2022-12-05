package oat

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// FilterOperationsById takes a slice of operation IDs to keep and removes non-matching operations
// in-place. The return value is the number of operations removed.
func FilterOperationsById(spec *openapi3.T, ids []string) (removedCount int) {
	s := make(map[string]struct{})
	for _, id := range ids {
		s[id] = struct{}{}
	}
	keep := func(opId string) bool {
		_, ok := s[opId]
		return ok
	}
	for pathKey, path := range spec.Paths {
		hasOperations := false
		for method, op := range path.Operations() {
			if !keep(op.OperationID) {
				path.SetOperation(method, nil)
				removedCount++
			} else {
				hasOperations = true
			}
		}
		if !hasOperations {
			delete(spec.Paths, pathKey)
		}
	}
	return
}
