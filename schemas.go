package oat

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/getkin/kin-openapi/openapi3"
)

const SchemaRefPrefix = "#/components/schemas/"

// SchemaCleaner identifies components in a schema which are never referenced from any operation,
// directly or indirectly. The zero-value of SchemaCleaner is ready to use.
type SchemaCleaner struct {
	// vis identifies components of the spec which have already been inspected. The unsafe.Pointer
	// keys bypass type safety, allowing for arbitrary pointer types to be stored.
	vis map[unsafe.Pointer]struct{}

	// referenced identifies names of components (without the #/components/schema prefix) which
	// have been referenced by some chain of $ref references, starting at an operation.
	referenced map[string]struct{}
}

// visited is called before any component is examined by a "dfs" function. It returns true if the
// component has already been visited or if the component is nil, otherwise it returns false.
func (c *SchemaCleaner) visited(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
	if c.vis == nil {
		c.vis = make(map[unsafe.Pointer]struct{})
	}
	_, ok := c.vis[ptr]
	if !ok {
		c.vis[ptr] = struct{}{}
	}
	return ok
}

func (c *SchemaCleaner) recordRef(spec *openapi3.T, ref string) {
	if ref == "" {
		return
	}
	if c.referenced == nil {
		c.referenced = make(map[string]struct{})
	}
	name := strings.TrimPrefix(ref, SchemaRefPrefix)
	sr := spec.Components.Schemas[name]
	if sr == nil {
		fmt.Printf("schema ref not found: %s\n", ref)
	}
	if _, ok := c.referenced[name]; !ok {
		c.referenced[name] = struct{}{}
		c.dfsSchemaRef(spec, sr)
	}
}

func (c *SchemaCleaner) dfsSchema(spec *openapi3.T, s *openapi3.Schema) {
	if c.visited(unsafe.Pointer(s)) {
		return
	}
	var srs openapi3.SchemaRefs
	srs = append(srs, s.OneOf...)
	srs = append(srs, s.AnyOf...)
	srs = append(srs, s.AllOf...)
	srs = append(srs, s.Not, s.Items, s.AdditionalProperties)
	for _, p := range s.Properties {
		srs = append(srs, p)
	}
	c.dfsSchemaRefs(spec, srs)
}

func (c *SchemaCleaner) dfsSchemaRef(spec *openapi3.T, sr *openapi3.SchemaRef) {
	if c.visited(unsafe.Pointer(sr)) {
		return
	}
	c.recordRef(spec, sr.Ref)
	c.dfsSchema(spec, sr.Value)
}

func (c *SchemaCleaner) dfsSchemaRefs(spec *openapi3.T, srs openapi3.SchemaRefs) {
	for _, sr := range srs {
		c.dfsSchemaRef(spec, sr)
	}
}

func (c *SchemaCleaner) dfsMediaType(spec *openapi3.T, m *openapi3.MediaType) {
	if c.visited(unsafe.Pointer(m)) {
		return
	}
	c.dfsSchemaRef(spec, m.Schema)
}

func (c *SchemaCleaner) dfsContent(spec *openapi3.T, ct openapi3.Content) {
	for _, m := range ct {
		c.dfsMediaType(spec, m)
	}
}

func (c *SchemaCleaner) dfsResponse(spec *openapi3.T, r *openapi3.Response) {
	if c.visited(unsafe.Pointer(r)) {
		return
	}
	c.dfsContent(spec, r.Content)
}

func (c *SchemaCleaner) dfsResponseRef(spec *openapi3.T, rr *openapi3.ResponseRef) {
	if c.visited(unsafe.Pointer(rr)) {
		return
	}
	c.recordRef(spec, rr.Ref)
	c.dfsResponse(spec, rr.Value)
}

func (c *SchemaCleaner) dfsResponses(spec *openapi3.T, rs openapi3.Responses) {
	for _, rr := range rs {
		c.dfsResponseRef(spec, rr)
	}
}

func (c *SchemaCleaner) dfsParameter(spec *openapi3.T, p *openapi3.Parameter) {
	if c.visited(unsafe.Pointer(p)) {
		return
	}
	c.dfsSchemaRef(spec, p.Schema)
}

func (c *SchemaCleaner) dfsParameterRef(spec *openapi3.T, pr *openapi3.ParameterRef) {
	if c.visited(unsafe.Pointer(pr)) {
		return
	}
	c.recordRef(spec, pr.Ref)
	c.dfsParameter(spec, pr.Value)
}

func (c *SchemaCleaner) dfsParameters(spec *openapi3.T, ps openapi3.Parameters) {
	for _, pr := range ps {
		c.dfsParameterRef(spec, pr)
	}
}

// Clean removes all unused components in-place, respecting the slice of names to keep. The return
// value is the number of components removed.
func (c *SchemaCleaner) Clean(spec *openapi3.T, excepts []string) (removedCount int) {
	for _, path := range spec.Paths {
		for _, op := range path.Operations() {
			c.dfsParameters(spec, op.Parameters)
			c.dfsResponses(spec, op.Responses)
		}
	}
	for _, e := range excepts {
		if sr, ok := spec.Components.Schemas[e]; ok {
			c.recordRef(spec, e)
			c.dfsSchemaRef(spec, sr)
		}
	}
	for name := range spec.Components.Schemas {
		if _, ok := c.referenced[name]; !ok {
			delete(spec.Components.Schemas, name)
			removedCount++
		}
	}
	return
}
