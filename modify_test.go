package proto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifying(t *testing.T) {

	t.Run("Test modifying struct fields", func(t *testing.T) {

		var protoStruct Test
		var nonce int = Template(&protoStruct)
		Modify(&protoStruct)

		assert.Equal(t, fmt.Sprintf("StringPtr_%d_Updated", nonce), *protoStruct.StringPtr)
		assert.Equal(t, fmt.Sprintf("String_%d_Updated", nonce), protoStruct.String)

		assert.Equal(t, nonce + 1, *protoStruct.IntPtr)
		assert.Equal(t, nonce + 1, protoStruct.Int)

		assert.Equal(t, false, *protoStruct.BoolPtr)
		assert.Equal(t, false, protoStruct.Bool)

		assert.Equal(t,
			[]string{
				fmt.Sprintf("StringSlicePtr_%d", nonce),
				fmt.Sprintf("StringSlicePtr_%d_Updated", nonce),
			},
			*protoStruct.StringSlicePtr,
		)
		assert.Equal(t,
			[]string{
				fmt.Sprintf("StringSlice_%d", nonce),
				fmt.Sprintf("StringSlice_%d_Updated", nonce),
			},
			protoStruct.StringSlice,
		)

		assert.Equal(t,
			[]int{
				nonce,
				nonce+1,
			},
			*protoStruct.IntSlicePtr,
		)
		assert.Equal(t,
			[]int{
				nonce,
				nonce+1,
			},
			protoStruct.IntSlice,
		)	

		// var subStringExpected string = fmt.Sprintf("SubStrPtr_%d", nonce)
		// var subTestExpexted SubTest = SubTest{SubInt: nonce, SubStrPtr: &subStringExpected}
		// assert.Empty(t, deep.Equal(protoStruct.SubTest, subTestExpexted))
		// assert.Equal(t, protoStruct.SubInt, nonce)
		// assert.Equal(t, *protoStruct.SubStrPtr, subStringExpected)

		// assert.Empty(t, deep.Equal(protoStruct.StructSlice, []SubTest{subTestExpexted}))

	})
}
