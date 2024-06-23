package pkg

import (
	"encoding/json"
	"errors"
	"time"
)

func (oe OriginErr) MarshalJSON() ([]byte, error) {
	var rawMessages []string
	for _, err := range oe {
		rawMessages = append(rawMessages, err.Error())
	}
	return json.Marshal(rawMessages)
}
func (oe OriginErr) UnmarshalJSON(data []byte) error {
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		var singleErr string
		if err := json.Unmarshal(data, &singleErr); err == nil {
			oe = OriginErr{errors.New(singleErr)}
			return nil
		}
		return err
	}

	for _, rawMsg := range rawMessages {
		var errStr string
		if err := json.Unmarshal(rawMsg, &errStr); err != nil {
			return err
		}
		oe = append(oe, errors.New(errStr))
	}

	return nil
}

func (ee *extenedError) UnmarshalJSON(data []byte) error {
	var simpleError string
	if err := json.Unmarshal(data, &simpleError); err == nil {
		ee.Message = simpleError
		return nil
	}
	var aux struct {
		Class      *string         `json:"class"`
		Code       *string         `json:"code"`
		Message    string          `json:"message"`
		Origin     json.RawMessage `json:"original_error"`
		SourceFile string          `json:"file"`
		Line       int             `json:"line"`
		FuncName   string          `json:"func_name"`
		StackTrace []string        `json:"stack_trace"`
		Timestamp  time.Time       `json:"timestamp"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ee.Class = aux.Class
	ee.Code = aux.Code
	ee.Message = aux.Message
	ee.SourceFile = aux.SourceFile
	ee.Line = aux.Line
	ee.FuncName = aux.FuncName
	ee.StackTrace = aux.StackTrace
	ee.Timestamp = aux.Timestamp

	oError := &OriginErr{}

	if err := json.Unmarshal(aux.Origin, oError); err == nil {
		ee.Origin = *oError
		return nil
	}
	return nil
}
