package main

import (
"log"
"net/http"

"github.com/gorilla/mux"
)

func main() {
router := mux.NewRouter()

// Health check route
router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status": "healthy", "service": "payment-service"}`))
}).Methods("GET")

// Payment routes placeholders
router.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"message": "Payment service is running"}`))
}).Methods("POST")

log.Println("Payment service starting on :8004")
log.Fatal(http.ListenAndServe(":8004", router))
}
