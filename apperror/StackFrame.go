package apperror

import (
	"fmt"
	"github.com/streamrail/concurrent-map"
	"io/ioutil"
	"runtime"
	"strings"
)

// A StackFrame contains all necessary information about to generate a line
// in a callstack.
type StackFrame struct {
	// The path to the file containing this ProgramCounter
	File string
	// The LineNumber in that file
	LineNumber int
	// The Name of the function that contains this ProgramCounter
	Name string
	// The Package that contains this function
	Package string
	// The underlying ProgramCounter
	ProgramCounter uintptr
	// whether pc is within function of sigpanic
	InSigPanic bool
}

// NewStackFrame popoulates a stack frame object from the program counter.
func NewStackFrame(pc uintptr, faultInstruction bool) (frame *StackFrame) {
	frame = &StackFrame{ProgramCounter: pc}
	if frame.ProgramCounter == 0 {
		return
	}

	frameFunc := runtime.FuncForPC(frame.ProgramCounter)

	if frameFunc == nil {
		return
	}

	frame.InSigPanic = (frameFunc.Name() == "runtime.sigpanic")
	frame.Package, frame.Name = packageAndName(frameFunc)

	// pc - 1 because the program counters we use are usually return addresses,
	// and we want to show the line that corresponds to the function call
	xpc := pc
	if !faultInstruction {
		xpc--
	}
	frame.File, frame.LineNumber = frameFunc.FileLine(xpc)
	return
}

func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastSlash := strings.LastIndex(name, "/"); lastSlash >= 0 {
		pkg += name[:lastSlash] + "/"
		name = name[lastSlash+1:]
	}

	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

// String returns the stackframe formatted in the same way as go does
// in runtime/debug.Stack()
func (frame *StackFrame) String() string {
	str := fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.ProgramCounter)

	source := frame.SourceLine()

	return str + "\t" + frame.Name + ": " + source + "\n"
}

var sourceMap cmap.ConcurrentMap = cmap.New()

// SourceLine gets the line of code (from File and Line) of the original source if possible.
func (frame *StackFrame) SourceLine() string {
	var lines []string

	if value, exists := sourceMap.Get(frame.File); exists {
		lines = value.([]string)
	} else {
		data, _ := ioutil.ReadFile(frame.File)
		lines = strings.Split(string(data), "\n")

		sourceMap.Set(frame.File, lines)
	}

	if frame.LineNumber <= 0 || frame.LineNumber >= len(lines) {
		return "???"
	}

	// -1 because line-numbers are 1 based, but our array is 0 based
	return strings.Trim(lines[frame.LineNumber-1], " \t")
}
