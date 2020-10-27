package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// getGID get goroutine ID
func getGID() uint64 {
	// get Gorouting ID from go enviroment
	// Scott Mansfield
	// Goroutine IDs
	// https://blog.sgmansfield.com/2015/12/goroutine-ids/
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// loggingLineComposer ...
func loggingLineComposer(format string, level uint, parameters ...interface{}) []byte {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	msec := now.Nanosecond() / 1e6
	timeStamp := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d.%03d", year, month, day, hour, min, sec, msec)

	file := "???"
	line := 0

	var ok bool
	_, file, line, ok = runtime.Caller(2)

	if ok {
		// substring (needs to convert in a byte slice)
		// from chr past last '/' till end of filename
		file = string([]byte(file)[strings.LastIndexByte(file, '/')+1:])
	}
	lvlName := levelNames[level]
	gID := getGID()

	realParameters := append([]interface{}{timeStamp, file, line, lvlName, gID}, parameters...)

	return []byte(
		fmt.Sprintf(
			"%s | %s:%d | %s | %d | "+format+"\n",
			realParameters...,
		),
	)
}

const (
	// TraceLvl the lowest logging level, good for extra info
	TraceLvl uint = iota
	// DebugLvl development logging level
	DebugLvl
	// InfoLvl application logging level
	InfoLvl
	// WarningLvl oddities logging level
	WarningLvl
	// ErrorLvl errors logging level
	ErrorLvl
	// FatalLvl panic logging level
	FatalLvl
)

// levelNames log level names
var levelNames [6]string

// logger
type logger struct {
	synch     io.Writer
	async     chan<- []byte
	mutex     *sync.Mutex
	asyncFlag bool
	level     *uint
}

// Logger ___
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

// GetLevel return the current logging level
func (l logger) GetLevel() uint {
	return *l.level
}

// SetLevel change the current logging level
func (l logger) SetLevel(newLevel uint) {
	if newLevel > FatalLvl {
		newLevel = FatalLvl
	}
	*l.level = newLevel
}

func (l *logger) logLine(line []byte) {
	if l.asyncFlag == true {
		l.async <- line
		return
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.synch.Write(line)
}

// Print logs ignoring level restrictions (TraceLvl)
func (l logger) Print(format string, parameters ...interface{}) {
	l.logLine(loggingLineComposer(format, TraceLvl, parameters...))
}

// Trace logs extra informations (TraceLvl)
func (l logger) Trace(format string, parameters ...interface{}) {
	if *l.level > TraceLvl {
		return
	}
	l.logLine(loggingLineComposer(format, TraceLvl, parameters...))
}

// Debug logs development informations (DebugLvl)
func (l logger) Debug(format string, parameters ...interface{}) {
	if *l.level > DebugLvl {
		return
	}
	l.logLine(loggingLineComposer(format, DebugLvl, parameters...))
}

// Info logs application informations (InfoLvl)
func (l logger) Info(format string, parameters ...interface{}) {
	if *l.level > InfoLvl {
		return
	}
	l.logLine(loggingLineComposer(format, InfoLvl, parameters...))
}

// Warning logs oddities informations (WarningLvl)
func (l logger) Warning(format string, parameters ...interface{}) {
	if *l.level > WarningLvl {
		return
	}
	l.logLine(loggingLineComposer(format, WarningLvl, parameters...))
}

// Error logs problems informations (ErrorLvl)
func (l logger) Error(format string, parameters ...interface{}) {
	if *l.level > ErrorLvl {
		return
	}
	l.logLine(loggingLineComposer(format, ErrorLvl, parameters...))
}

// Fatal logs critical informations (FatalLvl)
func (l logger) Fatal(format string, parameters ...interface{}) {
	line := loggingLineComposer(format, FatalLvl, parameters...)

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.synch.Write(line)

	os.Exit(-1)
}

func (l logger) SetWriter(w io.Writer) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.synch = w
}

// asyncLogging helper function for run async logger in goroutine
func asyncLogging(w io.Writer, ch <-chan []byte) {
	for {
		line, ok := <-ch

		if !ok {
			break
		}
		_, err := w.Write(line)

		if err != nil {
			break
		}
	}
}

// NewLogger initialize a new logger
func NewLogger(writer io.Writer, asyncFlag bool, bufferingSize uint) Logger {
	channel := make(chan []byte, bufferingSize)

	go asyncLogging(writer, channel)

	level := InfoLvl

	return logger{
		synch:     writer,
		async:     channel,
		mutex:     &sync.Mutex{},
		asyncFlag: asyncFlag,
		level:     &level,
	}
}

// Initialize initialize the current logger
func Initialize(writer io.Writer, asyncFlag bool, bufferingSize uint) {
	if defaultLogger.async != nil {
		close(defaultLogger.async)
	}
	channel := make(chan []byte, bufferingSize)

	go asyncLogging(writer, channel)

	level := InfoLvl

	defaultLogger = logger{
		synch:     writer,
		async:     channel,
		mutex:     &sync.Mutex{},
		asyncFlag: asyncFlag,
		level:     &level,
	}
}

// defaultLogger base logger, initially tied with os.Stdoout
var defaultLogger logger

func init() {
	levelNames[TraceLvl] = "TRACE"
	levelNames[DebugLvl] = "DEBUG"
	levelNames[InfoLvl] = "INFO"
	levelNames[WarningLvl] = "WARNING"
	levelNames[ErrorLvl] = "ERROR"
	levelNames[FatalLvl] = "FATAL"

	Initialize(os.Stdout, false, 0)
}

// GetLevel return the current logging level
func GetLevel() uint {
	return defaultLogger.GetLevel()
}

// SetLevel change the current logging level
func SetLevel(newLevel uint) {
	defaultLogger.SetLevel(newLevel)
}

// SetLevelName change the label of level during logging
func SetLevelName(level uint, name string) {
	if level > FatalLvl {
		return
	}
	levelNames[level] = name
}

// Print logs ignoring level restrictions (TraceLvl)
func Print(format string, parameters ...interface{}) {
	defaultLogger.logLine(loggingLineComposer(format, TraceLvl, parameters...))
}

// Trace logs extra informations (TraceLvl)
func Trace(format string, parameters ...interface{}) {
	if *defaultLogger.level > TraceLvl {
		return
	}
	defaultLogger.logLine(loggingLineComposer(format, TraceLvl, parameters...))
}

// Debug logs development informations (DebugLvl)
func Debug(format string, parameters ...interface{}) {
	if *defaultLogger.level > DebugLvl {
		return
	}
	defaultLogger.logLine(loggingLineComposer(format, DebugLvl, parameters...))
}

// Info logs application informations (InfoLvl)
func Info(format string, parameters ...interface{}) {
	if *defaultLogger.level > InfoLvl {
		return
	}
	defaultLogger.logLine(loggingLineComposer(format, InfoLvl, parameters...))
}

// Warning logs oddities informations (WarningLvl)
func Warning(format string, parameters ...interface{}) {
	if *defaultLogger.level > WarningLvl {
		return
	}
	defaultLogger.logLine(loggingLineComposer(format, WarningLvl, parameters...))
}

// Error logs problems informations (ErrorLvl)
func Error(format string, parameters ...interface{}) {
	if *defaultLogger.level > ErrorLvl {
		return
	}
	defaultLogger.logLine(loggingLineComposer(format, ErrorLvl, parameters...))
}

// Fatal logs critical informations (FatalLvl)
func Fatal(format string, parameters ...interface{}) {
	line := loggingLineComposer(format, FatalLvl, parameters...)

	defaultLogger.mutex.Lock()
	defer defaultLogger.mutex.Unlock()
	defaultLogger.synch.Write(line)

	os.Exit(-1)
}
