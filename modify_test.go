package proto

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestModifying(t *testing.T) {

	t.Run("Test modifying struct fields", func(t *testing.T) {

		var protoStruct Test
		var nonce int = Prototype(&protoStruct)

		// Save some values for testing later
		structSliceBefore := protoStruct.StructSlice

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
		

		var subStringExpected string = fmt.Sprintf("SubStrPtr_%d_Updated", nonce)
		var subTestExpexted SubTest = SubTest{SubInt: nonce + 1, SubStrPtr: &subStringExpected}
		assert.Empty(t, deep.Equal(protoStruct.SubTest, subTestExpexted))
		assert.Equal(t, protoStruct.SubInt, nonce + 1)
		assert.Equal(t, *protoStruct.SubStrPtr, subStringExpected)

		assert.NotContains(t, structSliceBefore, protoStruct.StructSlice[len(protoStruct.StructSlice) - 1])

	})

	t.Run("Modify fields with tags", func(t *testing.T) {

		type ModTest struct {
			String string `proto:"stringVal"`
			IntPtr *int    `proto:"1"`

			StringPtr *string `proto:"stringVal2" proto.modify:"modifiedVal"`
			Int	int `proto:"33" proto.modify:"44"`

			StringSlice []string `proto.modify:"val1,val2,val3"`
			IntSlicePtr *[]int `proto.modify:"5,6,7"`
		}

		var modTest ModTest
		Prototype(&modTest)
		Modify(&modTest)

		assert.Equal(t, "stringVal", modTest.String)
		assert.Equal(t, 1, *modTest.IntPtr)

		assert.Equal(t, "modifiedVal", *modTest.StringPtr)
		assert.Equal(t, 44, modTest.Int)

		assert.Empty(t, deep.Equal([]string{"val1","val2","val3"}, modTest.StringSlice))
		assert.Empty(t, deep.Equal([]int{5,6,7}, *modTest.IntSlicePtr))

	})
}
