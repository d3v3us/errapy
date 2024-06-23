package pkg

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type ErrorExtended interface {
	error
	IsClassifed() bool
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

type extenedError struct {
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
type class struct {
	codePrefix string
	code       int
	name       string
	codeRange  *[]int
}

func (e *extenedError) IsClassifed() bool {
	return e.Class != nil && *e.Class != ""
}
func (e *extenedError) IsCoded() bool {
	return e.Code != nil && *e.Code != ""
}
func (e *extenedError) GetClass() string {
	return *e.Class
}

func (e *extenedError) GetOrigin() error {
	return e.Origin
}
func (e *extenedError) GetPrevious() error {
	return e.Previous
}
func (e *extenedError) GetSourceFile() string {
	return e.SourceFile
}
func (e *extenedError) GetLine() int {
	return e.Line
}
func (e *extenedError) GetFuncName() string {
	return e.FuncName
}
func (e *extenedError) GetStackTrace() []string {
	return e.StackTrace
}
func (e *extenedError) GetTime() string {
	return e.Timestamp.Format(time.RFC3339)
}
func (e *extenedError) Error() string {
	return e.Message
}

type errorExtendedBuilder struct {
	*extenedError
}

func NewBuilder() *errorExtendedBuilder {
	return &errorExtendedBuilder{extenedError: &extenedError{}}
}

func (e *errorExtendedBuilder) SetClass(class string) *errorExtendedBuilder {
	e.Class = &class
	return e
}
func (e *errorExtendedBuilder) SetCode(code string) *errorExtendedBuilder {
	e.Code = &code
	return e
}
func (e *errorExtendedBuilder) SetMessage(message string) *errorExtendedBuilder {
	e.Message = message
	return e
}
func (e *errorExtendedBuilder) SetOrigin(origin error) *errorExtendedBuilder {
	e.Origin = OriginErr{origin}

	return e
}
func (e *errorExtendedBuilder) SetPrevious(previous error) *errorExtendedBuilder {
	e.Previous = previous
	return e
}
func (e *errorExtendedBuilder) SetSourceFile(sourceFile string) *errorExtendedBuilder {
	e.SourceFile = sourceFile
	return e
}
func (e *errorExtendedBuilder) SetLine(line int) *errorExtendedBuilder {
	e.Line = line
	return e
}
func (e *errorExtendedBuilder) SetFuncName(funcName string) *errorExtendedBuilder {
	e.FuncName = funcName
	return e
}
func (e *errorExtendedBuilder) SetStackTrace(stackTrace []string) *errorExtendedBuilder {
	e.StackTrace = stackTrace
	return e
}
func (e *errorExtendedBuilder) SetTimestamp(timestamp time.Time) *errorExtendedBuilder {
	e.Timestamp = timestamp
	return e
}

func (e *errorExtendedBuilder) Build() ErrorExtended {
	pc, file, line, ok := Caller(2)
	e.SetTimestamp(time.Now())
	if ok {
		e.SetStackTrace(GetStackTrace(2))
		e.SetSourceFile(file)
		e.SetLine(line)
		fn := FuncForPC(pc)
		e.SetFuncName(fn.Name())
	}
	return e.extenedError
}
func Class(name string, codeRange *[]int) *class {
	return &class{
		name:      name,
		codeRange: codeRange}
}
func (c *class) nextCode() string {
	if c.codeRange != nil {
		if c.code+1 >= (*c.codeRange)[1] {
			panic("error code out of range")
		}
		c.code = c.code + 1
		return fmt.Sprintf("%s%d", c.codePrefix, c.code)
	}
	return ""
}

func (c *class) New(msg string) error {
	return NewBuilder().SetClass(c.name).SetCode(c.nextCode()).SetMessage(msg).Build()
}

func (c *class) Newf(format string, args ...interface{}) error {
	return NewBuilder().SetClass(c.name).SetCode(c.nextCode()).SetMessage(fmt.Sprintf(format, args...)).Build()
}

func New(msg string) error {
	return NewBuilder().SetMessage(msg).Build()
}
func From(err error) error {
	return NewBuilder().SetOrigin(err).SetMessage(err.Error()).Build()
}
func Newf(format string, args ...interface{}) error {
	return NewBuilder().SetMessage(fmt.Sprintf(format, args...)).Build()
}

func (c *class) Wrap(err error, msg string) error {
	return NewBuilder().
		SetClass(c.name).
		SetCode(c.nextCode()).
		SetMessage(msg).
		SetOrigin(err).
		Build()
}

func (c *class) Wrapf(err error, format string, args ...interface{}) error {
	return NewBuilder().
		SetClass(c.name).
		SetCode(c.nextCode()).
		SetMessage(fmt.Sprintf(format, args...)).
		SetOrigin(err).
		Build()
}
