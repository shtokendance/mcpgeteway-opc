package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// PrepareImport applies an optional title and normalizes Apifox OpenAPI schema references.
func PrepareImport(data []byte, title string) ([]byte, error) {
	var spec map[string]any
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse specification: %w", err)
	}
	info, ok := spec["info"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("OpenAPI info is required")
	}
	if title = strings.TrimSpace(title); title != "" {
		info["title"] = title
	}
	name, _ := info["title"].(string)
	if strings.TrimSpace(name) == "默认模块" {
		return nil, fmt.Errorf("OpenAPI title cannot be 默认模块; please provide a different name")
	}
	if version, _ := spec["openapi"].(string); strings.HasPrefix(version, "3.0") {
		normalizeSchemaRefs(spec, false, false)
	}
	return json.Marshal(spec)
}

func normalizeSchemaRefs(value any, inSchema, inComponents bool) {
	switch node := value.(type) {
	case map[string]any:
		if inSchema {
			if ref, ok := node["$ref"].(string); ok && len(node) == 2 {
				if _, ok := node["description"]; ok {
					delete(node, "$ref")
					node["allOf"] = []any{map[string]any{"$ref": ref}}
				}
			}
		}
		for key, child := range node {
			if key == "example" || key == "examples" || key == "default" || strings.HasPrefix(key, "x-") {
				continue
			}
			normalizeSchemaRefs(child, inSchema || key == "schema" || (inComponents && key == "schemas"), key == "components")
		}
	case []any:
		for _, child := range node {
			normalizeSchemaRefs(child, inSchema, false)
		}
	}
}
