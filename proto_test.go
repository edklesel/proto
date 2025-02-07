package proto

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

type SubTest struct {
	SubInt    int
	SubStrPtr *string
}

type Test struct {
	String         string
	StringPtr      *string
	IntPtr         *int
	Int            int
	BoolPtr        *bool
	Bool           bool
	StringSlicePtr *[]string
	StringSlice    []string
	IntSlicePtr    *[]int
	IntSlice       []int
	StructSlice    []SubTest
	SubTest
}

type TestWithTags struct {
	StringNoTag string
	StringTag   string	`proto:"TagValue"`
	IntNoTag    int
	IntTag      int		`proto:"1"`
	BoolNoTag   bool
	BoolTag     bool    `proto:"false"`
	StringSliceNoTag []string
	StringSliceTag []string  `proto:"TagValue1,TagValue2"`
	IntSliceNoTag []int
	IntSliceTag []int  `proto:"1,2,3"`
}

func TestPrototyping(t *testing.T) {

	t.Run("Test basic prototyping", func(t *testing.T) {

		var protoStruct Test
		var nonce int = Template(&protoStruct)

		assert.Equal(t, *protoStruct.StringPtr, fmt.Sprintf("StringPtr_%d", nonce))
		assert.Equal(t, protoStruct.String, fmt.Sprintf("String_%d", nonce))

		assert.Equal(t, *protoStruct.IntPtr, nonce)
		assert.Equal(t, protoStruct.Int, nonce)

		assert.Equal(t, *protoStruct.BoolPtr, true)
		assert.Equal(t, protoStruct.Bool, true)

		assert.Equal(t, *protoStruct.StringSlicePtr, []string{fmt.Sprintf("StringSlicePtr_%d", nonce)})
		assert.Equal(t, protoStruct.StringSlice, []string{fmt.Sprintf("StringSlice_%d", nonce)})

		assert.Equal(t, *protoStruct.IntSlicePtr, []int{nonce})
		assert.Equal(t, protoStruct.IntSlice, []int{nonce})

		var subStringExpected string = fmt.Sprintf("SubStrPtr_%d", nonce)
		var subTestExpexted SubTest = SubTest{SubInt: nonce, SubStrPtr: &subStringExpected}
		assert.Empty(t, deep.Equal(protoStruct.SubTest, subTestExpexted))
		assert.Equal(t, protoStruct.SubInt, nonce)
		assert.Equal(t, *protoStruct.SubStrPtr, subStringExpected)

		assert.Empty(t, deep.Equal(protoStruct.StructSlice, []SubTest{subTestExpexted}))

	})

	t.Run("Test prototyping with tags", func(t *testing.T) {

		var proto TestWithTags
		var nonce int = Template(&proto)

		assert.Equal(t, proto.StringNoTag, fmt.Sprintf("StringNoTag_%d", nonce))
		assert.Equal(t, proto.StringTag, "TagValue")

		assert.Equal(t, proto.IntNoTag, nonce)
		assert.Equal(t, proto.IntTag, 1)

		assert.Equal(t, proto.BoolNoTag, true)
		assert.Equal(t, proto.BoolTag, false)

		assert.Equal(t, proto.StringSliceNoTag, []string{fmt.Sprintf("StringSliceNoTag_%d", nonce)})
		assert.Equal(t, proto.StringSliceTag, []string{"TagValue1","TagValue2"})

		assert.Equal(t, proto.IntSliceNoTag, []int{nonce})
		assert.Equal(t, proto.IntSliceTag, []int{1,2,3})

	})

}
