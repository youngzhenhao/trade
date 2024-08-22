package btldb

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"trade/middleware"
)

type QueryParams map[string]interface{}

func parseTagField(tag string) string {
	// Parse the gorm tag of the struct field and extract the column part
	tagParts := strings.Split(tag, ";")
	for _, part := range tagParts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}

// GenericQuery fetches records directly using the middleware's DB instance.
// It filters records according to the provided params and returns a slice of found models or an error.
func GenericQuery[T any](model *T, params QueryParams) ([]*T, error) {
	var results []*T
	// Obtain the reflection type object of the model
	modelType := reflect.TypeOf(model).Elem()
	// Start with a base model and apply filters.
	query := middleware.DB.Model(model)
	// Validate each key in the query parameters against the model's fields
	for key, value := range params {
		field, ok := modelType.FieldByName(key)
		if !ok {
			log.Printf("Invalid query field: %s is not a field of the model", key)
			return nil, fmt.Errorf("invalid query field: %s", key)
		}
		columnName := parseTagField(string(field.Tag.Get("gorm")))
		if columnName == "" {
			// Use the field name as a fallback to the column name
			columnName = key
		}
		query = query.Where(columnName+" = ?", value)
	}
	// Execute the query and find the results.
	if err := query.Find(&results).Error; err != nil {
		log.Printf("Error querying database: %v", err)
		return nil, err
	}
	// Optionally print the results (assuming model has a method String() string to print its details)
	//for _, result := range results {
	//	fmt.Println(result)
	//}
	return results, nil
}

func GenericQueryByObject[T any](condition *T) ([]*T, error) {
	var results []*T
	// Start with a base model and apply filters using the condition instance where only non-zero values are considered.
	if err := middleware.DB.Where(condition).Find(&results).Error; err != nil {
		fmt.Printf("Error querying database: %v\n", err)
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	// Optionally print the results (assuming model has a method String() string to print its details)
	//for _, result := range results {
	//	fmt.Println(result)
	//}
	return results, nil
}
