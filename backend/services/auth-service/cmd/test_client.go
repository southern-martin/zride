package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type LoginRequest struct {
	Code string `json:"code"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         interface{} `json:"user"`
}

func main() {
	// Wait for the auth service to start
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:8081"
	
	// Test health check
	fmt.Println("Testing health check...")
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		log.Printf("Health check failed: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Println("✓ Health check passed")
	} else {
		fmt.Printf("✗ Health check failed with status: %d\n", resp.StatusCode)
		os.Exit(1)
	}
	
	// Test login endpoint (will fail with mock data, but should return proper error)
	fmt.Println("Testing login endpoint...")
	loginReq := LoginRequest{Code: "mock_zalo_code"}
	jsonData, _ := json.Marshal(loginReq)
	
	resp, err = http.Post(baseURL + "/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Login test failed: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	// We expect this to fail since it's a mock code, but the endpoint should respond
	if resp.StatusCode == 400 || resp.StatusCode == 503 || resp.StatusCode == 500 {
		fmt.Println("✓ Login endpoint is responsive (expected failure with mock data)")
	} else {
		fmt.Printf("✗ Unexpected response from login endpoint: %d\n", resp.StatusCode)
	}
	
	fmt.Println("Auth service basic tests completed!")
}