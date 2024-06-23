package pkg

import (
	"encoding/json"
	"errors"
)

func (oe *originalError) MarshalJSON() ([]byte, error) {
	var rawMessages []string
	for _, err := range oe.Errors {
		rawMessages = append(rawMessages, err.Error())
	}
	return json.Marshal(rawMessages)
}
func (oe *originalError) UnmarshalJSON(data []byte) error {
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		var singleErr string
		if err := json.Unmarshal(data, &singleErr); err == nil {
			oe.Errors = []error{errors.New(singleErr)}
			return nil
		}
		return err
	}

	for _, rawMsg := range rawMessages {
		var errStr string
		if err := json.Unmarshal(rawMsg, &errStr); err != nil {
			return err
		}
		oe.Errors = append(oe.Errors, errors.New(errStr))
	}

	return nil
}

func (ee *extenedError) UnmarshalJSON(data []byte) error {
	var simpleError string
	if err := json.Unmarshal(data, &simpleError); err == nil {
		ee.Message = simpleError
		return nil
	}

	var aux extenedError
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

	oError := &originalError{}

	if err := json.Unmarshal(aux.Origin, oError); err == nil {
		ee.OriginalErr = *oError
		return nil
	}
	return nil
}
