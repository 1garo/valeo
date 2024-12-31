package main

import (
	"errors"
	"fmt"
	"reflect"
)

type name string
type age int

// Example struct with fields that can have their own Validate methods
type MyStruct struct {
	Name  name
	Age   age
	Other float32
}

// Validate method for the Name field
func (n *name) Validate() error {
	if *n == "" {
		return errors.New("Name cannot be empty")
	}
	return nil
}

// Validate method for the Age field
func (a *age) Validate() error {
	if *a < 0 {
		return errors.New("Age cannot be negative")
	}
	return nil
}

// Generic validation function to validate all fields of a struct
func ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)

	// Ensure we have a pointer to a struct
	if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
		return errors.New("ValidateStruct requires a pointer to a struct")
	}

	val = val.Elem()
	typ := val.Type()

	// Iterate through fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Check if the field has a Validate method
		validateMethod := field.Addr().MethodByName("Validate")
		if validateMethod.IsValid() {
			// Call the Validate method
			results := validateMethod.Call(nil)
			fmt.Printf("here: %+v\n", results)
			if len(results) == 1 && !results[0].IsNil() {
				return fmt.Errorf("Validation failed for field '%s': %v", fieldType.Name, results[0].Interface())
			}
		} else {
			return fmt.Errorf("Validation failed for field '%s': Validate function not implemented", fieldType.Name)
		}
	}

	return nil
}

func Run() {
	example := &MyStruct{
		Name: "nameeee",
		Age:  5,
	}

	if err := ValidateStruct(example); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
