package output

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
)

// Stdout struct
type Stdout struct {
	output map[string]string
	mut    *sync.Mutex
}

func outputOk(pre string, str string) string {
	col := color.New(color.FgGreen)
	return pre + col.Sprintln(str)
}

func outputErr(pre string, str string) string {
	col := color.New(color.FgRed)
	return pre + col.Sprintln(str)
}

func checkPre(ok bool) string {
	pre := "ok"
	col := color.New(color.FgGreen)
	if !ok {
		pre = "ERROR"
		col = color.New(color.FgRed)
	}
	return fmt.Sprintf("[%s]", col.Sprintf("%s", pre))
}

func outputTitle(str string) string {
	col := color.New(color.FgBlue).Add(color.Bold)
	return col.Sprintln(str)
}

// get the entry or add a new one
func (o *Stdout) getOrAdd(key string) string {
	v, ok := o.output[key]
	if !ok {
		v = outputTitle(fmt.Sprintf("checking %s", key))
	}
	return v
}

// StackErr add a new error
func (o *Stdout) StackErr(key string, pre string, content string) {
	o.mut.Lock()
	v := o.getOrAdd(key)
	defer o.mut.Unlock()

	// append error
	v += "  "
	v += checkPre(false)
	v += outputErr(fmt.Sprintf(" %s: ", pre), content)
	o.output[key] = v
}

// StackOk add a new success
func (o *Stdout) StackOk(key string, pre string, content string) {
	o.mut.Lock()
	v := o.getOrAdd(key)
	defer o.mut.Unlock()

	// append success
	v += "  "
	v += checkPre(true)
	v += outputOk(fmt.Sprintf(" %s: ", pre), content)
	o.output[key] = v
}

// Flush flush output
func (o *Stdout) Flush(key string) {
	o.mut.Lock()
	defer o.mut.Unlock()

	v, ok := o.output[key]
	if ok {
		fmt.Print(v)
	}
}

// NewStdout new instance
func NewStdout(_ map[string]string) (*Stdout, error) {
	o := &Stdout{
		output: make(map[string]string),
		mut:    &sync.Mutex{},
	}
	return o, nil
}
