package pkg

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type ErrorClass struct {
	codePrefix string
	code       int
	name       string
	codeRange  *[]int
}
type ExtendedError interface {
	error
	IsClassified() bool
	IsCoded() bool
	GetClass() string
	GetOrigin() error
	GetPrevious() error
	GetSourceFile() string
	GetLine() int
	GetFuncName() string
	GetStackTrace() []string
	GetTime() string
}
type OriginErr []error

func (oe OriginErr) Error() string {
	buff := bytes.NewBufferString("")
	for i := 0; i < len(oe); i++ {

		buff.WriteString((oe)[i].Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

type extendedError struct {
	Code       *string   `json:"code"`
	Class      *string   `json:"class"`
	Message    string    `json:"message"`
	Origin     OriginErr `json:"original_error"`
	SourceFile string    `json:"file"`
	Line       int       `json:"line"`
	FuncName   string    `json:"func_name"`
	StackTrace []string  `json:"-"`
	Timestamp  time.Time `json:"timestamp"`
	Previous   error     `json:"-"`
}

func (e *extendedError) IsClassified() bool {
	return e.Class != nil && *e.Class != ""
}
func (e *extendedError) IsCoded() bool {
	return e.Code != nil && *e.Code != ""
}
func (e *extendedError) GetClass() string {
	return *e.Class
}

func (e *extendedError) GetOrigin() error {
	return e.Origin
}
func (e *extendedError) GetPrevious() error {
	return e.Previous
}
func (e *extendedError) GetSourceFile() string {
	return e.SourceFile
}
func (e *extendedError) GetLine() int {
	return e.Line
}
func (e *extendedError) GetFuncName() string {
	return e.FuncName
}
func (e *extendedError) GetStackTrace() []string {
	return e.StackTrace
}
func (e *extendedError) GetTime() string {
	return e.Timestamp.Format(time.RFC3339)
}
func (e *extendedError) Error() string {
	return e.Message
}

type ErrorBuilder struct {
	extendedError *extendedError
}

func NewBuilder() *ErrorBuilder {
	return &ErrorBuilder{extendedError: &extendedError{}}
}

func (e *ErrorBuilder) SetClass(class string) *ErrorBuilder {
	e.extendedError.Class = &class
	return e
}
func (e *ErrorBuilder) SetCode(code string) *ErrorBuilder {
	e.extendedError.Code = &code
	return e
}
func (e *ErrorBuilder) SetMessage(message string) *ErrorBuilder {
	e.extendedError.Message = message
	return e
}
func (e *ErrorBuilder) SetOrigin(origin error) *ErrorBuilder {
	e.extendedError.Origin = OriginErr{origin}

	return e
}
func (e *ErrorBuilder) SetPrevious(previous error) *ErrorBuilder {
	e.extendedError.Previous = previous
	return e
}
func (e *ErrorBuilder) SetSourceFile(sourceFile string) *ErrorBuilder {
	e.extendedError.SourceFile = sourceFile
	return e
}
func (e *ErrorBuilder) SetLine(line int) *ErrorBuilder {
	e.extendedError.Line = line
	return e
}
func (e *ErrorBuilder) SetFuncName(funcName string) *ErrorBuilder {
	e.extendedError.FuncName = funcName
	return e
}
func (e *ErrorBuilder) SetStackTrace(stackTrace []string) *ErrorBuilder {
	e.extendedError.StackTrace = stackTrace
	return e
}
func (e *ErrorBuilder) SetTimestamp(timestamp time.Time) *ErrorBuilder {
	e.extendedError.Timestamp = timestamp
	return e
}

func (e *ErrorBuilder) Build() ExtendedError {
	pc, file, line, ok := Caller(2)
	e.SetTimestamp(time.Now())
	if e.extendedError.IsCoded() {
		e.SetMessage(fmt.Sprintf("%s: %s", *e.extendedError.Code, e.extendedError.Message))
	} else if defaultPolicy.CodesRequired {
		panic("error code not set")
	}

	if defaultPolicy.ClassesRequired && !e.extendedError.IsClassified() {
		panic("error class not set")
	}
	if ok {
		e.SetStackTrace(GetStackTrace(2))
		e.SetSourceFile(file)
		e.SetLine(line)
		e.SetFuncName(FuncForPC(pc).Name())
	}
	return e.extendedError
}
func Class(name string, codePrefix string, codeRange *[]int) *ErrorClass {
	return &ErrorClass{
		name:       name,
		codePrefix: codePrefix,
		codeRange:  codeRange}
}
func (c *ErrorClass) nextCode() string {
	if c.codeRange != nil {
		if c.code == 0 {
			c.code = (*c.codeRange)[0]
		}
		if c.code+1 >= (*c.codeRange)[1] {
			panic("error code out of range")
		}
		c.code = c.code + 1
		return fmt.Sprintf("%s%d", c.codePrefix, c.code)
	}
	return ""
}

func (c *ErrorClass) New(msg string) error {
	return NewBuilder().SetClass(c.name).SetCode(c.nextCode()).SetMessage(msg).Build()
}

func (c *ErrorClass) Newf(format string, args ...interface{}) error {
	return NewBuilder().SetClass(c.name).SetCode(c.nextCode()).SetMessage(fmt.Sprintf(format, args...)).Build()
}

//goland:noinspection GoUnusedFunction
func New(msg string) error {
	return NewBuilder().SetMessage(msg).Build()
}

//goland:noinspection GoUnusedFunction
func From(err error) error {
	return NewBuilder().SetOrigin(err).SetMessage(err.Error()).Build()
}

//goland:noinspection GoUnusedFunction
func Newf(format string, args ...interface{}) error {
	return NewBuilder().SetMessage(fmt.Sprintf(format, args...)).Build()
}

func (c *ErrorClass) Wrap(err error, msg string) error {
	return NewBuilder().
		SetClass(c.name).
		SetCode(c.nextCode()).
		SetMessage(msg).
		SetOrigin(err).
		Build()
}

func (c *ErrorClass) Wrapf(err error, format string, args ...interface{}) error {
	return NewBuilder().
		SetClass(c.name).
		SetCode(c.nextCode()).
		SetMessage(fmt.Sprintf(format, args...)).
		SetOrigin(err).
		Build()
}
