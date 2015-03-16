package beanstream

import (
	"fmt"
)

type BeanstreamApiException struct {
	Status    int
	Code      int
	Category  int
	Message   string
	Reference string
	Details   []ErrorDetail
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *BeanstreamApiException) Error() string {
	return fmt.Sprintf("%v %v: code(%v) category(%v) message: %v   ref: %v  details: %v", e.Status, e.ErrorType(), e.Code, e.Category, e.Message, e.Reference, e.Details)
}

func (e *BeanstreamApiException) ErrorType() string {
	switch e.Status {
	case -1:
		return "UnexpectedException"
	case 302:
		return "BusinessRuleException"
	case 400:
		return "InvalidRequestException"
	case 401:
		return "UnauthorizedException"
	case 402:
		return "BusinessRuleException"
	case 403:
		return "ForbiddenException"
	case 404:
		return "NotFoundException"
	case 405:
		return "InvalidRequestException"
	case 415:
		return "InvalidRequestException"
	default:
		return "InternalServerException"
	}
}
