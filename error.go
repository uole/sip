package sip

import "fmt"

type SipError struct {
	Code    int
	Message string
}

func (s *SipError) Error() string {
	return fmt.Sprintf("SIP ERROR (%d): %s", s.Code, s.Message)
}

func newSipError(code int, message string) *SipError {
	return &SipError{
		Code:    code,
		Message: message,
	}
}
