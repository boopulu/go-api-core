package core

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
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


type RPCError struct {
	Code string
	Desc string
}


func ParseRPCError(input string) RPCError {
	// Define a regular expression to match the code and desc fields
	re := regexp.MustCompile(`code\s*=\s*([^ ]+)\s*desc\s*=\s*([^ ]+)`)

	// Find the matches in the input string
	matches := re.FindStringSubmatch(input)

	// Check if there are enough submatches (i.e., if both code and desc were found)
	if len(matches) != 3 {
		// Handle the case where the input format is invalid
		return RPCError{
			Code: "Unknown",
			Desc: "Invalid input format",
		}
	}

	// Extract the code and desc from the submatches
	code := matches[1]
	desc := matches[2]

	return RPCError{
		Code: code,
		Desc: desc,
	}
}

func RpcCodeToHTTPStatus(codeStr string) (int, error) {
	switch codeStr {
	case "OK":
		return http.StatusOK, nil // HTTP 200 OK
	case "Canceled":
		return http.StatusRequestTimeout, nil // HTTP 408 Request Timeout
	case "Unknown":
		return http.StatusInternalServerError, nil // HTTP 500 Internal Server Error
	case "InvalidArgument":
		return http.StatusBadRequest, nil // HTTP 400 Bad Request
	case "DeadlineExceeded":
		return http.StatusGatewayTimeout, nil // HTTP 504 Gateway Timeout
	case "NotFound":
		return http.StatusNotFound, nil // HTTP 404 Not Found
	case "AlreadyExists":
		return http.StatusConflict, nil // HTTP 409 Conflict
	case "PermissionDenied":
		return http.StatusForbidden, nil // HTTP 403 Forbidden
	case "Unauthenticated":
		return http.StatusUnauthorized, nil // HTTP 401 Unauthorized
	case "ResourceExhausted":
		return http.StatusTooManyRequests, nil // HTTP 429 Too Many Requests
	case "FailedPrecondition":
		return http.StatusPreconditionFailed, nil // HTTP 412 Precondition Failed
	case "Aborted":
		return http.StatusConflict, nil // HTTP 409 Conflict
	case "OutOfRange":
		return http.StatusBadRequest, nil // HTTP 400 Bad Request
	case "Unimplemented":
		return http.StatusNotImplemented, nil // HTTP 501 Not Implemented
	case "Internal":
		return http.StatusInternalServerError, nil // HTTP 500 Internal Server Error
	case "Unavailable":
		return http.StatusServiceUnavailable, nil // HTTP 503 Service Unavailable
	case "DataLoss":
		return http.StatusBadGateway, nil // HTTP 502 Bad Gateway
	default:
		return http.StatusInternalServerError, fmt.Errorf("unrecognized gRPC code: %s", codeStr)
	}
}
