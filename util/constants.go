package util

import "errors"

var (
	// HTTP
	HTTP_ANGSURAN              = "HTTP|ANGSURAN"
	SUCCESS_CALCULATE_ANGSURAN = "Success calculate angsuran"
	FAIL_CALCULATE_ANGSURAN    = "Fail calculate angsuran"

	// USECASE ANGSURAN
	USECASE_ANGSURAN = "USECASE|ANGSURAN"

	// USECASE ANGSURAN ERROR
	ErrUsecaseANGSURAN = errors.New("error - usecase ANGSURAN")

	// COMMON ERROR
	ErrInvalidPayload                     = errors.New("error - invalid request payload")
	ErrInternalServer                     = errors.New("error - internal server error")
	ErrDataIsNil                          = errors.New("error - input data is nil")
	ErrInvalidLoanDurationAndInterestRate = errors.New("error - invalid loan duration or interest rate")
)
