package crypto

import (
	"testing"

	"github.com/golly-go/golly/utils"
	"github.com/magiconair/properties/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	examples :=
		[]interface{}{
			"somedata",
			1,
			100000000.0,
		}

	for _, example := range examples {
		b := utils.GetBytes(example)
		encrypted := Encrypt(b, "testing123")

		x := Decrypt(encrypted, "testing123")

		assert.Equal(t, b, x)

	}
}
