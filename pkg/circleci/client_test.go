package circleci

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
)

func newClient() (*Client, error) {
	if apiToken := os.Getenv("CIRCLECI_API_TOKEN"); apiToken != "" {
		client := NewClient()
		if err := client.Authenticate(apiToken); err != nil {
			return nil, err
		}
		return client, nil
	}
	return nil, errors.New("testing of pkg/circleci need env CIRCLECI_API_TOKEN set")
}

func TestNewClient(t *testing.T) {
	if _, err := newClient(); err != nil {
		t.Errorf(fmt.Sprintf("%+v", errors.Wrap(err, "failed to authenticate client")))
	}
}
