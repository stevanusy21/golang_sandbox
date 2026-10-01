package domain

import "errors"

var (
	ErrRecordNotFound          = errors.New("Record not found")
	ErrFailedToProcessPayment  = errors.New("Failed to process payment")
	ErrFailedToSavePayment     = errors.New("Failed to save payment")
	ErrFailedToUpdatePayment   = errors.New("Failed to update payment")
	ErrFailedToPublishEvent    = errors.New("Failed to publish event")
	ErrFailedToHandleCallback  = errors.New("Failed to handle callback")
	ErrQueryFailed             = errors.New("Query failed")
)
