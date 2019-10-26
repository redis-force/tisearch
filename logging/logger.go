package logging

import "fmt"

func Debugf(f string, val ...interface{}) {
	fmt.Printf("DEBUG: "+f+"\n", val...)
}
func Infof(f string, val ...interface{}) {
	fmt.Printf("INFO: "+f+"\n", val...)
}

func Warnf(f string, val ...interface{}) {
	fmt.Printf("WARN: "+f+"\n", val...)
}

func Errorf(f string, val ...interface{}) {
	fmt.Printf("ERROR: "+f+"\n", val...)
}

func Fatalf(f string, val ...interface{}) {
	panic(fmt.Sprintf("FATAL: "+f+"\n", val...))
}

func Debug(val interface{}) {
	fmt.Printf("DEBUG: %v\n", val)
}
func Info(val interface{}) {
	fmt.Printf("INFO: %v\n", val)
}

func Warn(val interface{}) {
	fmt.Printf("WARN: %v\n", val)
}

func Error(val interface{}) {
	fmt.Printf("ERROR: %v\n", val)
}

func Fatal(val interface{}) {
	fmt.Printf("ERROR: %v\n", val)
}
