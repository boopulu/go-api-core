package core

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	postgresSession *sql.DB
)

// Create a random string of length `length` consisting of lowercase letters and numbers
func GenerateSecret(length int) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var b strings.Builder
	for i := 0; i < length; i++ {
		n := r.Intn(36)
		if n < 26 {
			b.WriteRune(rune(n + 'a'))
		} else {
			b.WriteRune(rune(n - 26 + '0'))
		}
	}
	return b.String()
}

// Initialize the postgres session.
// URLformat: postgres://user:password@host:port/database
func InitPostgres(url string) {
	connStr := url

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	postgresSession = db
}

// Get the postgres session
func GetPostgresSession() *sql.DB {
	return postgresSession
}

// Send a PUT request to the main api, to let it know that this microservice is running
func StartMicroservice(host_address, host_port, id, secret, address, port string) {

	requestURL := fmt.Sprintf("http://%s:%s/v1/microservice/start", host_address, host_port)

	// create a PUT request, with a body that contains the microservice ID and secret

	requestBody := []byte(fmt.Sprintf(`{"microservice_id": "%s", "microservice_secret": "%s", "microservice_address": "%s", "microservice_port": "%s"}`, id, secret, address, port))

	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(requestBody))

	if err != nil {
		panic(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:")
		panic(err.Error())
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		fmt.Println("PUT request was successful!")
	} else {
		fmt.Printf("PUT request failed with status: %s", resp.Status)
		panic("PUT request failed")
	}

	// You can read the response body if needed
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:")
		panic(err.Error())

	} else {
		fmt.Println("Response body:", string(responseBody))
	}
}