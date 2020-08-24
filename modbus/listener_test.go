package modbus

import (
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAndTransform(t *testing.T) {
	items := []struct {
		name string
		item Item
		res  []byte
		err  error
	}{
		{
			name: "two_bytes_coil",
			item: Item{
				Function: Coil,
				Quantity: 9,
				Value:    []int16{1, 0, 1, 0, 0, 0, 1, 1, 1},
			},
			res: []byte{197, 1},
			err: nil,
		},
		{
			name: "one_bytes_coil",
			item: Item{
				Function: Coil,
				Quantity: 8,
				Value:    []int16{1, 0, 0, 0, 1, 0, 0, 0},
			},
			res: []byte{17},
			err: nil,
		},
		{
			name: "register",
			item: Item{
				Function: HoldingRegister,
				Quantity: 3,
				Value:    []int16{16, 256, 34},
			},
			res: []byte{0, 16, 1, 0, 0, 34},
			err: nil,
		},
		{
			name: "register",
			item: Item{
				Function: HoldingRegister,
				Quantity: 3,
				Value:    []int16{16},
			},
			res: nil,
			err: errors.Errorf("quantity not equal to value length"),
		},
	}
	for _, tt := range items {
		t.Run(tt.name, func(t *testing.T) {
			res, err := validateAndTransform(tt.item)
			if err != nil {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, tt.err)
			}
			assert.Equal(t, res, tt.res)
		})
	}
}
