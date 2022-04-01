package errors

import "fmt"

/*
Task: Errors needed for multiwriter

You may find this blog post useful:
http://blog.golang.org/error-handling-and-go

Similar to a the Stringer interface, the error interface also defines a
method that returns a string.

type error interface {
    Error() string
}

Thus also the error type can describe itself as a string. The fmt package (and
many others) use this Error() method to print errors.

Implement the Error() method for the Errors type defined above.

The following conditions should be covered:

1. When there are no errors in the slice, it should return:

"(0 errors)"

2. When there is one error in the slice, it should return:

The error string return by the corresponding Error() method.

3. When there are two errors in the slice, it should return:

The first error + " (and 1 other error)"

4. When there are X>1 errors in the slice, it should return:

The first error + " (and X other errors)"
*/
func (m Errors) Error() string {
	var er Errors
	for i := 0; i < len(m); i++ {
		if m[i] != nil {
			er = append(er, m[i])
		}
	}

	errLen := len(er)

	switch errLen {
	case 1:
		return er[0].Error()
	case 2:
		return er[0].Error() + " (and 1 other error)"
	case 0:
		return "(0 errors)"
	default:
		str := fmt.Sprintf(" (and %d other errors)", errLen-1)
		return er[0].Error() + str
	}
}
