package sip

const (
	MaxCseq = 2147483647
)

const (
	StatusTrying                       = 100 //An attempt is made to transfer the call.
	StatusRinging                      = 180 //An attempt is made to ring from the called party.
	StatusCallIsBeingForwarded         = 181 //The call is forwarded.
	StatusQueued                       = 182 //The call is on hold.
	StatusSessionProgress              = 183 //The connection is established.
	StatusEarlyDialogTerminated        = 199 //The dialog was closed during connection setup.
	StatusOK                           = 200 //The request has been processed successfully and the result of the request is transferred in the response.
	StatusAccepted                     = 202 //The request has been accepted, but will be executed at a later time.
	StatusNoNotification               = 204 //The request was executed successfully, but the corresponding response is deliberately not sent.
	StatusMultipleChoices              = 300 //There is no unique destination address for the remote terminal.
	StatusMovedPermanently             = 301 //The called party is permanently reachable somewhere else.
	StatusMovedTemporarily             = 302 //The called party is temporarily reachable somewhere else.
	StatusUseProxy                     = 305 //The specified proxy must be used.
	StatusAlternativeService           = 380 //The call was not successful, but alternative services are available.
	StatusBadRequest                   = 400 //The SIP request is incorrect.
	StatusUnauthorized                 = 401 //The authorization is incorrect.
	StatusPaymentRequired              = 402 //Not yet defined; intended for "not enough credit available".
	StatusForbidden                    = 403 //The request was invalid.
	StatusNotFound                     = 404 //The remote terminal was not found or does not exist.
	StatusMethodNotAllowed             = 405 //The method of the request (for example, SUBSCRIBE or NOTIFY) is not allowed.
	StatusNotAcceptable                = 406 //The call options are not allowed.
	StatusProxyAuthenticationRequired  = 407 //The proxy needs authorization.
	StatusRequestTimeout               = 408 //Timeout - The remote peer does not respond within a reasonable time.
	StatusGone                         = 410 //The desired subscriber can no longer be reached at the specified address.
	StatusConditionalRequestFailed     = 412 //The prerequisites for processing the request could not be met because a request required for this failed.
	StatusRequestEntityTooLarge        = 413 //The message content is too large.
	StatusRequestURITooLong            = 414 //The SIP address (URI) of the request is too long.
	StatusUnsupportedMediaType         = 415 //The codec is not supported.
	StatusUnsupportedURIScheme         = 416 //The SIP address is incorrect.
	StatusUnknownResourcePriority      = 417 //The request should be treated with a certain priority, but the server does not understand the details.
	StatusBadExtension                 = 420 //The server does not understand a protocol extension.
	StatusExtensionRequired            = 421 //The server needs a protocol extension.
	StatusSessionIntervalTooSmall      = 422 //The Session Expires value is too low for the server.
	StatusIntervalTooBrief             = 423 //The value of the desired machining time is too short.
	StatusUseIdentityHeader            = 428 //The identity header is missing.
	StatusProvideReferrerIdentity      = 429 //No valid referred by token is specified.
	StatusFlowFailed                   = 430 //The particular routing failed (proxy internal, endpoints should treat the response like code 400).
	StatusAnonymityDisallowed          = 433 //The server refuses to process anonymous requests.
	StatusBadIdentityInfo              = 436 //The SIP address contained in the identity header is invalid, unavailable, or not supported.
	StatusUnsupportedCertificate       = 437 //The verifier cannot verify the certificate in the identity header.
	StatusInvalidIdentityHeader        = 438 //The certificate in the identity header is invalid.
	StatusFirstHopLacksOutboundSupport = 439 //The registrar supports outbound feature, but the proxy used does not.
	StatusMaxBreadthExceeded           = 440 //It is no longer possible to derive concurrent forks from the query.
	StatusBadInfoPackage               = 469 //Unsuitable Info-Package - Transmission error, resend.
	StatusConsentNeeded                = 470 //The server has no access rights to at least one of the specified SIP addresses.
	StatusTemporarilyUnavailable       = 480 //The called party is currently not reachable.
	StatusCallTransactionDoesNotExist  = 481 //This connection does not exist (anymore).
	StatusLoopDetected                 = 482 //A forwarding loop has been detected.
	StatusTooManyHops                  = 483 //Too many forwarding steps were identified.
	StatusAddressIncomplete            = 484 //The SIP address is incomplete.
	StatusAmbiguous                    = 485 //The SIP address cannot be uniquely resolved.
	StatusBusyHere                     = 486 //The called party is busy.
	StatusRequestTerminated            = 487 //The call attempt was aborted.
	StatusNotAcceptableHere            = 488 //Illegal call attempt.
	StatusBadEvent                     = 489 //The server does not know the specified event.
	StatusRequestPending               = 491 //A request from the same dialog is still being processed.
	StatusUndecipherable               = 493 //The request contains an encrypted MIME body that the recipient cannot decrypt.
	StatusSecurityAgreementRequired    = 494 //The request requires a security agreement, but does not include a security mechanism supported by the server.
	StatusServerInternalError          = 500 //Internal server error.
	StatusNotImplemented               = 501 //The server does not support the SIP request.
	StatusBadGateway                   = 502 //The gateway in the SIP request is corrupted.
	StatusServiceUnavailable           = 503 //The server's SIP service is temporarily unavailable.
	StatusServerTimeout                = 504 //The server cannot reach another server in a reasonable time.
	StatusVersionNotSupported          = 505 //The SIP protocol version is not supported by the server.
	StatusMessageTooLarge              = 513 //The SIP message is too large for UDP; TCP must be used.
	StatusPreconditionFailure          = 580 //The server cannot or does not want to meet the requirements for processing the request.
	StatusBusyEverywhere               = 600 //All terminal devices of the called subscriber are occupied.
	StatusDeclined                     = 603 //The called party has rejected the call attempt.
	StatusDoesNotExistAnywhere         = 604 //The called party no longer exists.
	StatusPartyHangsUp                 = 701 //The called party has hung up.
)

var statusText = map[int]string{
	StatusTrying:                       "Trying",
	StatusRinging:                      "Ringing",
	StatusCallIsBeingForwarded:         "Call Is Being Forwarded",
	StatusQueued:                       "Queued",
	StatusSessionProgress:              "Session Progress",
	StatusEarlyDialogTerminated:        "Early Dialog Terminated",
	StatusOK:                           "OK",
	StatusAccepted:                     "Accepted",
	StatusNoNotification:               "No Notification",
	StatusMultipleChoices:              "Multiple Choices",
	StatusMovedPermanently:             "Moved Permanently",
	StatusMovedTemporarily:             "Moved Temporarily",
	StatusUseProxy:                     "Use Proxy",
	StatusAlternativeService:           "Alternative Service",
	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthenticationRequired:  "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusGone:                         "Gone",
	StatusConditionalRequestFailed:     "Conditional Request Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusUnsupportedURIScheme:         "Unsupported URI Scheme",
	StatusUnknownResourcePriority:      "Unknown Resource-Priority",
	StatusBadExtension:                 "Bad Extension",
	StatusExtensionRequired:            "Extension Required",
	StatusSessionIntervalTooSmall:      "Session Interval Too Small",
	StatusIntervalTooBrief:             "Interval Too Brief",
	StatusUseIdentityHeader:            "Use Identity Header",
	StatusProvideReferrerIdentity:      "Provide Referrer Identity",
	StatusFlowFailed:                   "Flow Failed",
	StatusAnonymityDisallowed:          "Anonymity Disallowed",
	StatusBadIdentityInfo:              "Bad Identity-Info",
	StatusUnsupportedCertificate:       "Unsupported Certificate",
	StatusInvalidIdentityHeader:        "Invalid Identity Header",
	StatusFirstHopLacksOutboundSupport: "First Hop Lacks Outbound Support",
	StatusMaxBreadthExceeded:           "Max-Breadth Exceeded",
	StatusBadInfoPackage:               "Bad Info Package",
	StatusConsentNeeded:                "Consent Needed",
	StatusTemporarilyUnavailable:       "Temporarily Unavailable",
	StatusCallTransactionDoesNotExist:  "Call/Transaction Does Not Exist",
	StatusLoopDetected:                 "Loop Detected",
	StatusTooManyHops:                  "Too Many Hops",
	StatusAddressIncomplete:            "Address Incomplete",
	StatusAmbiguous:                    "Ambiguous",
	StatusBusyHere:                     "Busy Here",
	StatusRequestTerminated:            "Request Terminated",
	StatusNotAcceptableHere:            "Not Acceptable Here",
	StatusBadEvent:                     "Bad Event",
	StatusRequestPending:               "Request Pending",
	StatusUndecipherable:               "Undecipherable",
	StatusSecurityAgreementRequired:    "Security Agreement Required",
	StatusServerInternalError:          "Server Internal Error",
	StatusNotImplemented:               "Not Implemented",
	StatusBadGateway:                   "Bad Gateway",
	StatusServiceUnavailable:           "Service Unavailable",
	StatusServerTimeout:                "Server Time-out",
	StatusVersionNotSupported:          "Version Not Supported",
	StatusMessageTooLarge:              "Message Too Large",
	StatusPreconditionFailure:          "Precondition Failure",
	StatusBusyEverywhere:               "Busy Everywhere",
	StatusDeclined:                     "Declined",
	StatusDoesNotExistAnywhere:         "Does Not Exist Anywhere",
	StatusPartyHangsUp:                 "Party Hangs Up",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}
