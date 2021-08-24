// Copyright (c) 2021 deadc0de6

package remote

import (
	"github.com/fatih/color"
)

func outputOk(pre string, str string) string {
	col := color.New(color.FgGreen)
	return pre + col.Sprintln(str)
}

func outputErr(pre string, str string) string {
	col := color.New(color.FgRed)
	return pre + col.Sprintln(str)
}

func outputTitle(str string) string {
	col := color.New(color.FgBlue).Add(color.Bold)
	return col.Sprintln(str)
}
