package output

import "fmt"

// Output struct
type Output interface {
	StackErr(string, string, string)
	StackOk(string, string, string)
	Flush(string)
}

// GetOutput returns an output instance
func GetOutput(name string, options map[string]string) (Output, error) {
	switch name {
	case "stdout":
		return NewStdout(options)
	}
	return nil, fmt.Errorf("no such output: %s", name)
}
