package logging

import "fmt"

func Debugf(f string, val ...interface{}) {
	fmt.Printf(f+"\n", val...)
}
func Infof(f string, val ...interface{}) {
	fmt.Printf(f+"\n", val...)
}

func Warnf(f string, val ...interface{}) {
	fmt.Printf(f+"\n", val...)
}

func Errorf(f string, val ...interface{}) {
	fmt.Printf(f+"\n", val...)
}
