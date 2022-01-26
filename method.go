package sip

import "strings"

type Method string

// Determine if the given method equals some other given method.
// This is syntactic sugar for case insensitive equality checking.
func (method *Method) Equals(other *Method) bool {
	if method != nil && other != nil {
		return strings.EqualFold(string(*method), string(*other))
	} else {
		return method == other
	}
}

func (method *Method) Is(s string) bool {
	return strings.ToLower(string(*method)) == strings.ToLower(s)
}

//String
func (method Method) String() string {
	return string(method)
}

// It's nicer to avoid using raw strings to represent methods, so the following standard
// method names are defined here as constants for convenience.
const (
	MethodInvite    Method = "INVITE"
	MethodAck       Method = "ACK"
	MethodCancel    Method = "CANCEL"
	MethodBye       Method = "BYE"
	MethodRegister  Method = "REGISTER"
	MethodOptions   Method = "OPTIONS"
	MethodSubscribe Method = "SUBSCRIBE"
	MethodNotify    Method = "NOTIFY"
	MethodRefer     Method = "REFER"
)
