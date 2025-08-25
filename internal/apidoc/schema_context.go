package apidoc

import (
	"fmt"
	"strings"
)

// SchemaGenerationContext holds pre-processed data for efficient schema generation
type SchemaGenerationContext struct {
	collection         CollectionInfo
	fieldSchemas       map[string]*FieldSchema
	systemFieldSchemas map[string]*FieldSchema
	fieldMapper        SchemaMapper
	includeExamples    bool
	includeSystem      bool
}

// NewSchemaGenerationContext creates a context for efficient schema generation
func NewSchemaGenerationContext(
	collection CollectionInfo,
	fieldMapper SchemaMapper,
	includeExamples bool,
	includeSystem bool,
) (*SchemaGenerationContext, error) {
	if collection.Name == "" {
		return nil, fmt.Errorf("collection name is required")
	}

	ctx := &SchemaGenerationContext{
		collection:         collection,
		fieldSchemas:       make(map[string]*FieldSchema),
		systemFieldSchemas: fieldMapper.GetSystemFieldSchemas(),
		fieldMapper:        fieldMapper,
		includeExamples:    includeExamples,
		includeSystem:      includeSystem,
	}

	// Pre-process all field schemas in a single pass
	if err := ctx.preprocessFieldSchemas(); err != nil {
		return nil, err
	}

	return ctx, nil
}

// preprocessFieldSchemas processes all collection fields once and caches the results
func (ctx *SchemaGenerationContext) preprocessFieldSchemas() error {
	// Process collection fields
	for _, field := range ctx.collection.Fields {
		fieldSchema, err := ctx.fieldMapper.MapFieldToSchema(field)
		if err != nil {
			// Use fallback schema if mapping fails
			fieldSchema = ctx.fieldMapper.GetFallbackSchema(field.Type)
		}
		ctx.fieldSchemas[field.Name] = fieldSchema
	}
	return nil
}

// GenerateCollectionSchema generates the main collection schema using pre-processed data
func (ctx *SchemaGenerationContext) GenerateCollectionSchema() (*CollectionSchema, error) {
	schema := &CollectionSchema{
		Type:       "object",
		Properties: make(map[string]*FieldSchema),
		Required:   []string{},
	}

	// Add system fields if enabled
	if ctx.includeSystem {
		for fieldName, fieldSchema := range ctx.systemFieldSchemas {
			schema.Properties[fieldName] = fieldSchema
			if fieldSchema.Required {
				schema.Required = append(schema.Required, fieldName)
			}
		}
	}

	// Add pre-processed collection fields
	for fieldName, fieldSchema := range ctx.fieldSchemas {
		schema.Properties[fieldName] = fieldSchema
		if fieldSchema.Required {
			schema.Required = append(schema.Required, fieldName)
		}
	}

	// Generate example if enabled
	if ctx.includeExamples {
		schema.Example = ctx.generateCollectionExample(schema)
	}

	return schema, nil
}

// GenerateCreateSchema generates the create schema using pre-processed data
func (ctx *SchemaGenerationContext) GenerateCreateSchema() (*CollectionSchema, error) {
	schema := &CollectionSchema{
		Type:       "object",
		Properties: make(map[string]*FieldSchema),
		Required:   []string{},
	}

	// Only add collection fields (no system fields for create operations)
	for _, field := range ctx.collection.Fields {
		// Skip system fields in create schema, except for password and email fields in auth collections
		if field.System && field.Type != "password" && !(ctx.collection.Type == "auth" && field.Type == "email") {
			continue
		}

		// Also skip common system fields by name (created, updated, id)
		if field.Name == "created" || field.Name == "updated" || field.Name == "id" {
			continue
		}

		// Use pre-processed field schema
		fieldSchema, exists := ctx.fieldSchemas[field.Name]
		if !exists {
			// Fallback if somehow not pre-processed
			var err error
			fieldSchema, err = ctx.fieldMapper.MapFieldToSchema(field)
			if err != nil {
				fieldSchema = ctx.fieldMapper.GetFallbackSchema(field.Type)
			}
		}

		schema.Properties[field.Name] = fieldSchema
		if fieldSchema.Required {
			schema.Required = append(schema.Required, field.Name)
		}
	}

	// For auth collections, add passwordConfirm field if password field exists
	if ctx.collection.Type == "auth" {
		if _, hasPassword := schema.Properties["password"]; hasPassword {
			passwordConfirmSchema := &FieldSchema{
				Type:        "string",
				Format:      "password",
				Description: "Password confirmation (must match password)",
				Required:    true,
			}
			schema.Properties["passwordConfirm"] = passwordConfirmSchema
			schema.Required = append(schema.Required, "passwordConfirm")
		}
	}

	// Special case: For auth collections, ensure email field is included if it exists
	if ctx.collection.Type == "auth" {
		for _, field := range ctx.collection.Fields {
			if field.Type == "email" && field.Name == "email" {
				// Use pre-processed field schema
				fieldSchema, exists := ctx.fieldSchemas[field.Name]
				if !exists {
					// Fallback if somehow not pre-processed
					var err error
					fieldSchema, err = ctx.fieldMapper.MapFieldToSchema(field)
					if err != nil {
						fieldSchema = &FieldSchema{
							Type:        "string",
							Format:      "email",
							Description: "Email address",
							Required:    field.Required,
							Example:     "user@example.com",
						}
					}
				}
				schema.Properties["email"] = fieldSchema
				if fieldSchema.Required {
					// Check if "email" is already in the required slice
					found := false
					for _, req := range schema.Required {
						if req == "email" {
							found = true
							break
						}
					}
					if !found {
						schema.Required = append(schema.Required, "email")
					}
				}
				break
			}
		}
	}

	// Generate example if enabled
	if ctx.includeExamples {
		schema.Example = ctx.generateCreateExample(schema)
	}

	return schema, nil
}

// GenerateUpdateSchema generates the update schema using pre-processed data
func (ctx *SchemaGenerationContext) GenerateUpdateSchema() (*CollectionSchema, error) {
	schema := &CollectionSchema{
		Type:       "object",
		Properties: make(map[string]*FieldSchema),
		Required:   []string{}, // No required fields for updates
	}

	// Add collection fields (no system fields for update operations)
	for _, field := range ctx.collection.Fields {
		// Skip system fields in update schema, except for password and email fields in auth collections
		if field.System && field.Type != "password" && !(ctx.collection.Type == "auth" && field.Type == "email") {
			continue
		}

		// Also skip common system fields by name (created, updated, id)
		if field.Name == "created" || field.Name == "updated" || field.Name == "id" {
			continue
		}

		// Use pre-processed field schema
		fieldSchema, exists := ctx.fieldSchemas[field.Name]
		if !exists {
			// Fallback if somehow not pre-processed
			var err error
			fieldSchema, err = ctx.fieldMapper.MapFieldToSchema(field)
			if err != nil {
				fieldSchema = ctx.fieldMapper.GetFallbackSchema(field.Type)
			}
		}

		// Make field optional for updates
		fieldSchema.Required = false
		schema.Properties[field.Name] = fieldSchema
	}

	// Special case: For auth collections, ensure email field is included if it exists
	if ctx.collection.Type == "auth" {
		for _, field := range ctx.collection.Fields {
			if field.Type == "email" && field.Name == "email" {
				// Use pre-processed field schema
				fieldSchema, exists := ctx.fieldSchemas[field.Name]
				if !exists {
					// Fallback if somehow not pre-processed
					var err error
					fieldSchema, err = ctx.fieldMapper.MapFieldToSchema(field)
					if err != nil {
						fieldSchema = &FieldSchema{
							Type:        "string",
							Format:      "email",
							Description: "Email address",
							Required:    false, // Always optional for updates
							Example:     "user@example.com",
						}
					}
				}
				fieldSchema.Required = false // Ensure it's optional for updates
				schema.Properties["email"] = fieldSchema
				break
			}
		}
	}

	// Generate example if enabled
	if ctx.includeExamples {
		schema.Example = ctx.generateUpdateExample(schema)
	}

	return schema, nil
}

// GenerateListResponseSchema generates the list response schema
func (ctx *SchemaGenerationContext) GenerateListResponseSchema() (map[string]any, error) {
	// Generate the item schema using our optimized method
	itemSchema, err := ctx.GenerateCollectionSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to generate item schema: %w", err)
	}

	// Create the list response schema
	listSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page": map[string]any{
				"type":        "integer",
				"description": "Current page number",
				"example":     1,
			},
			"perPage": map[string]any{
				"type":        "integer",
				"description": "Number of items per page",
				"example":     30,
			},
			"totalItems": map[string]any{
				"type":        "integer",
				"description": "Total number of items",
				"example":     100,
			},
			"totalPages": map[string]any{
				"type":        "integer",
				"description": "Total number of pages",
				"example":     4,
			},
			"items": map[string]any{
				"type":        "array",
				"description": fmt.Sprintf("Array of %s records", ctx.collection.Name),
				"items":       itemSchema,
			},
		},
		"required": []string{"page", "perPage", "totalItems", "totalPages", "items"},
	}

	// Add example if enabled
	if ctx.includeExamples {
		listSchema["example"] = ctx.generateListResponseExample(itemSchema)
	}

	return listSchema, nil
}

// generateCollectionExample generates an example object for a collection schema
func (ctx *SchemaGenerationContext) generateCollectionExample(schema *CollectionSchema) map[string]any {
	example := make(map[string]any)

	for fieldName, fieldSchema := range schema.Properties {
		if fieldSchema.Example != nil {
			example[fieldName] = fieldSchema.Example
		} else {
			// Generate a basic example based on type
			example[fieldName] = ctx.generateBasicExample(fieldSchema, fieldName)
		}
	}

	return example
}

// generateCreateExample generates an example object for create operations
func (ctx *SchemaGenerationContext) generateCreateExample(schema *CollectionSchema) map[string]any {
	example := make(map[string]any)

	for fieldName, fieldSchema := range schema.Properties {
		// Only include required fields and some optional fields in create examples
		if fieldSchema.Required || ctx.shouldIncludeInCreateExample(fieldName, fieldSchema) {
			if fieldSchema.Example != nil {
				example[fieldName] = fieldSchema.Example
			} else {
				example[fieldName] = ctx.generateBasicExample(fieldSchema, fieldName)
			}
		}
	}

	return example
}

// generateUpdateExample generates an example object for update operations
func (ctx *SchemaGenerationContext) generateUpdateExample(schema *CollectionSchema) map[string]any {
	example := make(map[string]any)

	// Include a subset of fields for update examples
	count := 0
	maxFields := 3 // Limit update examples to a few fields

	for fieldName, fieldSchema := range schema.Properties {
		if count >= maxFields {
			break
		}

		if ctx.shouldIncludeInUpdateExample(fieldName, fieldSchema) {
			if fieldSchema.Example != nil {
				example[fieldName] = fieldSchema.Example
			} else {
				example[fieldName] = ctx.generateBasicExample(fieldSchema, fieldName)
			}
			count++
		}
	}

	return example
}

// generateListResponseExample generates an example for list responses
func (ctx *SchemaGenerationContext) generateListResponseExample(itemSchema *CollectionSchema) map[string]any {
	// Generate a couple of item examples
	items := []any{}
	if itemSchema.Example != nil {
		items = append(items, itemSchema.Example)

		// Create a second example with slight variations
		if secondExample := ctx.createVariationExample(itemSchema.Example); secondExample != nil {
			items = append(items, secondExample)
		}
	}

	return map[string]any{
		"page":       1,
		"perPage":    30,
		"totalItems": len(items),
		"totalPages": 1,
		"items":      items,
	}
}

// generateBasicExample generates a basic example value based on field schema
func (ctx *SchemaGenerationContext) generateBasicExample(fieldSchema *FieldSchema, fieldName string) any {
	switch fieldSchema.Type {
	case "string":
		if fieldSchema.Format == "email" {
			return "user@example.com"
		} else if fieldSchema.Format == "uri" {
			return "https://example.com"
		} else if fieldSchema.Format == "date-time" {
			return "2024-01-01T12:00:00Z"
		} else if len(fieldSchema.Enum) > 0 {
			return fieldSchema.Enum[0]
		}
		return fmt.Sprintf("example_%s", fieldName)
	case "number", "integer":
		if fieldSchema.Minimum != nil {
			return *fieldSchema.Minimum + 1
		} else if fieldSchema.Maximum != nil {
			return *fieldSchema.Maximum - 1
		}
		return 42
	case "boolean":
		return true
	case "array":
		return []any{}
	case "object":
		return map[string]any{}
	default:
		return fmt.Sprintf("example_%s", fieldName)
	}
}

// shouldIncludeInCreateExample determines if a field should be included in create examples
func (ctx *SchemaGenerationContext) shouldIncludeInCreateExample(fieldName string, fieldSchema *FieldSchema) bool {
	// Always include required fields
	if fieldSchema.Required {
		return true
	}

	// Include some common optional fields
	commonFields := []string{"name", "title", "description", "email", "status", "active", "enabled", "password"}
	for _, common := range commonFields {
		if strings.Contains(strings.ToLower(fieldName), common) {
			return true
		}
	}

	// Include password fields (they're important for auth collections)
	if fieldSchema.Format == "password" {
		return true
	}

	// Include relation fields (they're important for understanding the API)
	// Note: We can't use the global isRelationFieldSchema function here, so we inline the logic
	if fieldSchema.Description != "" {
		if strings.Contains(fieldSchema.Description, "Related record ID") ||
			strings.Contains(fieldSchema.Description, "Relation field") {
			return true
		}
	}

	// Include boolean fields (they're usually simple and helpful)
	if fieldSchema.Type == "boolean" {
		return true
	}

	return false
}

// shouldIncludeInUpdateExample determines if a field should be included in update examples
func (ctx *SchemaGenerationContext) shouldIncludeInUpdateExample(fieldName string, fieldSchema *FieldSchema) bool {
	// Include commonly updated fields
	commonFields := []string{"name", "title", "description", "status", "active", "enabled"}
	for _, common := range commonFields {
		if strings.Contains(strings.ToLower(fieldName), common) {
			return true
		}
	}

	// Include text and boolean fields as they're commonly updated
	if fieldSchema.Type == "string" || fieldSchema.Type == "boolean" {
		return true
	}

	return false
}

// createVariationExample creates a variation of an example for list responses
func (ctx *SchemaGenerationContext) createVariationExample(original any) any {
	if originalMap, ok := original.(map[string]any); ok {
		variation := make(map[string]any)
		for key, value := range originalMap {
			variation[key] = ctx.createVariationValue(value, key)
		}
		return variation
	}
	return nil
}

// createVariationValue creates a variation of a single value
func (ctx *SchemaGenerationContext) createVariationValue(value any, key string) any {
	switch v := value.(type) {
	case string:
		if strings.Contains(v, "example_") {
			return strings.Replace(v, "example_", "sample_", 1)
		} else if v == "user@example.com" {
			return "admin@example.com"
		} else if v == "https://example.com" {
			return "https://sample.com"
		}
		return v + "_2"
	case int:
		return v + 1
	case float64:
		return v + 1.0
	case bool:
		return !v
	default:
		return v
	}
}
