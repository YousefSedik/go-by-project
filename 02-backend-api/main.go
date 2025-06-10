package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const PORT uint16 = 3000

func handleOperation(w http.ResponseWriter, r *http.Request, op func(float64, float64) (float64, error)) {
	var numbers NumbersModel
	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := op(numbers.Number1, numbers.Number2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseModel{Result: result})
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	handleOperation(w, r, func(a, b float64) (float64, error) { return a + b, nil })
}

func subtractHandler(w http.ResponseWriter, r *http.Request) {
	handleOperation(w, r, func(a, b float64) (float64, error) { return a - b, nil })
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	handleOperation(w, r, func(a, b float64) (float64, error) { return a * b, nil })
}

func divideHandler(w http.ResponseWriter, r *http.Request) {
	handleOperation(w, r, func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero is not allowed")
		}
		return a / b, nil
	})
}

func main() {
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/subtract", subtractHandler)
	http.HandleFunc("/multiply", multiplyHandler)
	http.HandleFunc("/divide", divideHandler)

	fmt.Printf("Server Running on localhost:%d\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
