package viewers

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestPrettyJSONError(t *testing.T) {
	_, err := prettyJSON("{invalid}")
	assert.Error(t, err)
}
