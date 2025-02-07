# Go Prototyping

This small utility library provides `proto.Template`, a tool which can take a blank struct input and populate the fields of that object with values based on the type of field.

This was designed primarily as a utility for writing CRUD tests, to generate a generalised struct object to reduce code duplication when writing tests.

**Disclaimer - This was written solely to help with one of my personal projects. I cannot guarantee that it will be regularly maintained.** However please feel free to contribute to expand the currently-limited range of types currently supported.

The populated values are (mostly) based on a random integer (referred to below as the `nonce`), and follow the below format:

| Field Type | Value |
| - | - |
| `string` | `"FieldName_<nonce>"` |
| `int` | `<nonce>` |
| `[]string` | `[]string{"FieldName_<nonce>"}` |
| `[]int` | `[]int{<nonce>,}` |
| `bool` | `true`|

Pointer fields are also templated, and the value of the pointer will correspond to the type as above.

The templating will also template struct fields of the initial struct, and slices of structs.

You can provide the `proto` field tag to struct fields, and the `Template` function will (attempt to) use this value instead. If there is a type mismatch, e.g a `string` is provided as the default value for an `int` field, then the function will panic.

An example usage can be found below:
```go
package main

import (
	"fmt"

	proto "github.com/edklesel/proto"
)

type Test struct {
    String  string
    IntPtr  *int
    Bool    bool
    String2 string `proto:"userValue"`
}

func main() {

    var test Test
    nonce := proto.Template(&test)

    fmt.Println(test.String == fmt.Sprintf("String_%d", nonce)) // true
    
    fmt.Println(*test.IntPtr == nonce) // true

    fmt.Println(test.Bool == true) // true

    fmt.Println(test.String2 == "userValue") // true

}
```