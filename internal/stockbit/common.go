package stockbit

import (
	"fmt"
	"net/http"
)

func (s *stockbit) handleError(response *http.Response, username string) error {
	if response.StatusCode == http.StatusUnauthorized {
		panic(fmt.Sprintf("token unauthorized: %s", username))
	}
	return fmt.Errorf("failed to get response: %d", response.StatusCode)
}
