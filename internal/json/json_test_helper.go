package json

import (
	"fmt"
	"testing"
)

func WithFailingJSONMarshal(t *testing.T, runner func(*testing.T)) {
	Marshal = func(_ interface{}) ([]byte, error) {
		return nil, fmt.Errorf("failed to marshal")
	}

	defer func() {
		Marshal = DefaultMarshal
	}()

	runner(t)
}

func WithFailingJSONUnmarshal(t *testing.T, runner func(*testing.T)) {
	Unmarshal = func(_ []byte, _ interface{}) error {
		return fmt.Errorf("failed to unmarshal")
	}

	defer func() {
		Unmarshal = DefaultUnmarshal
	}()

	runner(t)
}
