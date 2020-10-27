# logger
A simple logger for Go

## Introduction

The golang log package lacks of few important things : level and an asynchronous mode

Logging to different level is important for filter the details of the information that your application gives.
Althought the practise says to do not log in goroutines, probably some application sooner or later will have to.
The golang log package hide a sync.Mutex in the Output function, a bottleneck while used in concurrency.

## Details

This logger package, gives a default syncronous logger (os.Stdout)

The default logger configuration can be changed with 

```go
func Initialize(writer io.Writer, asyncFlag bool, bufferingSize uint)
```

or create a new one with

```go
func NewLogger(writer io.Writer, asyncFlag bool, bufferingSize uint)
```

Use the package level function for the default logger :

```go
func Print(format string, parameters ...interface{})
func Trace(format string, parameters ...interface{})
func Debug(format string, parameters ...interface{})
func Info(format string, parameters ...interface{})
func Warning(format string, parameters ...interface{})
func Error(format string, parameters ...interface{})
func Fatal(format string, parameters ...interface{})
func GetLevel() uint
func SetLevel(newLevel uint)
```

or the equivalent interface for the custom ones :

```go
type Logger interface {
	Print(format string, parameters ...interface{})
	Trace(format string, parameters ...interface{})
	Debug(format string, parameters ...interface{})
	Info(format string, parameters ...interface{})
	Warning(format string, parameters ...interface{})
	Error(format string, parameters ...interface{})
	Fatal(format string, parameters ...interface{})
	GetLevel() uint
	SetLevel(newLevel uint)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Disclaimer

This project has been an exercise for improving my GO skills, wrapping up things I already knew.

The package has never been used, it needs deep testing.

