# Go Prototyping

This small utility library provides a couple of small tools which can take a blank struct input and populate the fields of that object with values based on the type of field, and also modify those fields to other values.

The tools are:
- `proto.Prototype` - Generate a struct object with filled in fields based on their type.
- `proto.PrototypeWithVersion` - Like `proto.Prototype` but includes a version constraint for certain fields.
- `proto.Modify`    - Modifies the fields of a struct to a new value.
- `proto.ModifyWithVersion`    - Like `proto.Modify` but includes a version constraint for certain fields.

This was designed primarily as a utility for writing CRUD tests, to generate a generalised struct object to reduce code duplication when writing tests.

**Disclaimer - This was written solely to help with one of my personal projects. I cannot guarantee that it will be regularly maintained.** However please feel free to contribute to expand the currently-limited range of types currently supported.

The populated values are (mostly) based on a random integer (referred to below as the `nonce`), and follow the below format:

| Field Type | Value |
| - | - |
| `string`, `*string` | `"FieldName_<nonce>"` |
| `int`, `*int` | `<nonce>` |
| `[]string`, `*[]string` | `[]string{"FieldName_<nonce>"}` |
| `[]int`, `*[]int` | `[]int{<nonce>,}` |
| `bool`, `*bool` | `true`|

Pointer fields are also generated, and the value of the pointer will correspond to the type as above.

The prototyping will also generate struct fields of the initial struct, and slices of structs.

## Field Tags

### `proto` and `proto.modify`

The `proto` field tag can be added to struct fields with a given value, and the `Prototype` function will (attempt to) use this value instead. If there is a type mismatch, e.g a `string` is provided as the default value for an `int` field, then the function will panic.

At the moment, user-provided values for fields of struct type are not supported.

By default, **fields with the `proto` field tag are ignored during the modify stage**. However, the user can provide the `proto.modify` tag to specify the value to modify to if necessary.

| Field Type | Example Tag |
| - | - |
| `string`, `*string` | `proto:"stringval"` |
| `int`, `*int` | `proto:"123"` |
| `[]string`, `*[]string` | `proto:"val1,val2"` |
| `[]int`, `*[]int` | `proto:"1,2,3"` |
| `bool`, `*bool` | `proto:"false"` |

Note that the above example tags also apply to the `proto.modify` tag.

### `proto.constraint`

The `proto.constraint` tag takes a version constraint in the format expected by the `hashicorp/go-version` library, e.g:
- `> 1.0`
- `<= 2.1`
- `< 0.9, >=1.1`

Fields with this tag are then either prototyped if they satisfy the constraint, or set to the Zero value of the fields type if they do not satisfy the constraint.

# Example Usage

A basic example usage can be found below:
```go
package main

import (
	"fmt"

	proto   "github.com/edklesel/proto"
    version "github.com/hashicorp/go-version"
)

type Test struct {
    String           string
    IntPtr           *int   `proto.modify:"4"`
    Bool             bool
    String2          string `proto:"userValue"`
    DeprecatedString string `proto.constraint:"< 1.2"`
}

func main() {

    var test Test

    // Create a prototype of the Test struct
    nonce := proto.Prototype(&test)

    fmt.Println(test.String == fmt.Sprintf("String_%d", nonce)) // true
    
    fmt.Println(*test.IntPtr == nonce) // true

    fmt.Println(test.Bool == true) // true

    fmt.Println(test.String2 == "userValue") // true

    // Modify the fields
    Modify(&test)

    fmt.Println(test.String == fmt.Sprintf("String_%d_Updated", nonce)) // true

    fmt.Printnln(*test.IntPtr == 4) // true

    fmt.Println(test.Bool == false) // true

    fmt.Println(test.String2 == "userValue") // true

    // If we care about the version, create a version-sensitive prototype
    var test2 Test

    ver, _ := version.NewVersion("1.2")
    nonce = proto.PrototypeWithVersion(&test, ver)
    fmt.Println(test2.DeprecatedString == "") // false

    ver, _ := version.NewVersion("1.3")
    nonce = proto.PrototypeWithVersion(&test, ver)
    fmt.Println(test2.DeprecatedString == "") // true

}
```
