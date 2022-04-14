package merge

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"strings"

	types "github.com/bhojpur/os/pkg/config/data"
	"github.com/bhojpur/os/pkg/config/data/convert"
	"github.com/bhojpur/os/pkg/config/data/definition"
)

func APIUpdateMerge(schema *types.Schema, schemas *types.Schemas, dest, src map[string]interface{}, replace bool) map[string]interface{} {
	result := UpdateMerge(schema, schemas, dest, src, replace)
	if s, ok := dest["status"]; ok {
		result["status"] = s
	}
	if m, ok := dest["metadata"]; ok {
		result["metadata"] = mergeMetadata(convert.ToMapInterface(m), convert.ToMapInterface(src["metadata"]))
	}
	return result
}

func UpdateMerge(schema *types.Schema, schemas *types.Schemas, dest, src map[string]interface{}, replace bool) map[string]interface{} {
	return mergeMaps("", nil, schema, schemas, replace, dest, src)
}

func isProtected(k string) bool {
	if !strings.Contains(k, "bhojpur.net/") || (isField(k) && k != "field.bhojpur.net/creatorId") {
		return false
	}
	return true
}

func isField(k string) bool {
	return strings.HasPrefix(k, "field.bhojpur.net/")
}

func mergeProtected(dest, src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return dest
	}

	result := copyMap(dest)

	for k, v := range src {
		if isProtected(k) {
			continue
		}
		result[k] = v
	}

	for k := range dest {
		if isProtected(k) || isField(k) {
			continue
		}
		if _, ok := src[k]; !ok {
			delete(result, k)
		}
	}

	return result
}

func mergeMetadata(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	result := copyMap(dest)

	labels := convert.ToMapInterface(dest["labels"])
	srcLabels := convert.ToMapInterface(src["labels"])
	labels = mergeProtected(labels, srcLabels)

	annotations := convert.ToMapInterface(dest["annotations"])
	srcAnnotation := convert.ToMapInterface(src["annotations"])
	annotations = mergeProtected(annotations, srcAnnotation)

	result["labels"] = labels
	result["annotations"] = annotations

	return result
}

func merge(field, fieldType string, parentSchema, schema *types.Schema, schemas *types.Schemas, replace bool, dest, src interface{}) interface{} {
	if isMap(field, schema, schemas) {
		return src
	}

	sm, smOk := src.(map[string]interface{})
	dm, dmOk := dest.(map[string]interface{})
	if smOk && dmOk {
		fieldType, fieldSchema := getSchema(field, fieldType, parentSchema, schema, schemas)
		return mergeMaps(fieldType, schema, fieldSchema, schemas, replace, dm, sm)
	}
	return src
}

func getSchema(field, parentFieldType string, parentSchema, schema *types.Schema, schemas *types.Schemas) (string, *types.Schema) {
	if schema == nil {
		if definition.IsMapType(parentFieldType) && parentSchema != nil {
			subType := definition.SubType(parentFieldType)
			s := schemas.Schema(subType)
			if s != nil && s.InternalSchema != nil {
				s = s.InternalSchema
			}
			return subType, s
		}
		return "", nil
	}
	fieldType := schema.ResourceFields[field].Type
	s := schemas.Schema(fieldType)
	if s != nil && s.InternalSchema != nil {
		return fieldType, s.InternalSchema
	}
	return fieldType, s
}

func isMap(field string, schema *types.Schema, schemas *types.Schemas) bool {
	if schema == nil {
		return false
	}
	f := schema.ResourceFields[field]
	mapType := definition.IsMapType(f.Type)
	if !mapType {
		return false
	}

	subType := definition.SubType(f.Type)
	return schemas.Schema(subType) == nil
}

func mergeMaps(fieldType string, parentSchema, schema *types.Schema, schemas *types.Schemas, replace bool, dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	result := copyMapReplace(schema, dest, replace)
	for k, v := range src {
		result[k] = merge(k, fieldType, parentSchema, schema, schemas, replace, dest[k], v)
	}
	return result
}

func copyMap(src map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range src {
		result[k] = v
	}
	return result
}

func copyMapReplace(schema *types.Schema, src map[string]interface{}, replace bool) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range src {
		if replace && schema != nil {
			f := schema.ResourceFields[k]
			if f.Update {
				continue
			}
		}
		result[k] = v
	}
	return result
}
