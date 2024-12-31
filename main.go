package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Request represents the JSON request body for the /hello route
type Request struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Age  int    `json:"age" validate:"gte=0,lte=120"`
}

func (r *Response) Validate() error {
	return nil
}

// Response represents a generic JSON response
type Response struct {
	Message string   `json:"message"`
	Error   []string `json:"error,omitempty"`
}

var validate *validator.Validate

func main() {
	// Initialize the validator
	Run()
	validate = validator.New()

	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	name := query.Get("name")
	//ageParam := query.Get("age")
	var req Request
	req.Name = name
	req.Age = -1

	if err := validate.Struct(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: formatValidationError(err),
		})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	})
}

func respondJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func formatValidationError(err error) []string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldErr := range validationErrors {
			messages = append(messages, fmt.Sprintf("Field '%s' failed on '%s' tag", fieldErr.Field(), fieldErr.Tag()))
		}
		return messages
	}
	return []string{"Invalid input"}
}
