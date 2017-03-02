package apperror

import (
	"fmt"
	"net/http"
)

// The error codes.
type ErrorCode string

// NOTE: only external facing error codes are listed here.
const (
	InvalidHeader                 = "InvalidHeader"
	InvalidParameter              = "InvalidParameter"
	DuplicateParameter            = "DuplicateParameter"
	InvalidRequestBody            = "InvalidRequestBody"
	ResourceNotFound              = "ResourceNotFound"
	InternalServerError           = "InternalServerError"
	InvalidUpdate                 = "InvalidUpdate"
	InvalidOperation              = "InvalidOperation"
	NotAuthenticatedError         = "NotAuthenticated"
	AuthorizationError            = "AuthorizationError"
	ClientRequired                = "ClientRequired"
	FreelancerRequired            = "FreelancerRequired"
	PreConditionError             = "PreConditionError"
	FreelancerAccountAlreadyExist = "FreelancerAccountAlreadyExist"
	ItemSuspended                 = "ItemSuspended"
	ItemNotSuspended              = "ItemNotSuspended"
	OfferAlreadyAccepted          = "OfferAlreadyAccepted"
	InvalidOfferStatus            = "InvalidOfferStatus"
	ContractAlreadyExist          = "ContractAlreadyExist"
	TooManyQueryOption            = "TooManyQueryOption"
	ReportImproperItem            = "ReportImproperItem"
	ReportAlreadyRead             = "ReportAlreadyRead"
	ReportAlreadyExist            = "ReportAlreadyExist"
	ProposalAlreadyExist          = "ProposalAlreadyExist"
	ProposalAlreadyInvited        = "ProposalAlreadyInvited"
	OfferAlreadyExist             = "OfferAlreadyExist"
	HeadCountLimitExceeded        = "HeadCountLimitExceeded"
	AssociatedProposalNotFound    = "AssociatedProposalNotFound"
	JobAlreadyArchived            = "JobAlreadyArchived"
	JobNotActive                  = "JobNotActive"
	JobNotArchived                = "JobNotArchived"
	JobTagAlreadyExist            = "JobTagAlreadyExist"
	JobTagAlreadyDeleted          = "JobTagAlreadyDeleted"
	ConcurrentUpdate              = "ConcurrentUpdate"
	JobExpired                    = "JobExpired"
	NotEnoughBalance              = "NotEnoughBalance"
	ContributorAlreadyExist       = "ContributorAlreadyExist"
	GenericWrap                   = "GenericWrap"
	ConflictOperation             = "ConflictOperation"
)

var errorCodesToStatusCode = map[ErrorCode]int{
	InvalidHeader:                 http.StatusBadRequest,
	InvalidParameter:              http.StatusBadRequest,
	DuplicateParameter:            http.StatusBadRequest,
	InvalidRequestBody:            http.StatusBadRequest,
	ResourceNotFound:              http.StatusNotFound,
	InternalServerError:           http.StatusInternalServerError,
	InvalidOperation:              http.StatusBadRequest,
	FreelancerAccountAlreadyExist: http.StatusBadRequest,
	ItemNotSuspended:              http.StatusBadRequest,
	OfferAlreadyAccepted:          http.StatusBadRequest,
	ContractAlreadyExist:          http.StatusBadRequest,
	TooManyQueryOption:            http.StatusBadRequest,
	InvalidUpdate:                 http.StatusBadRequest,
	ReportImproperItem:            http.StatusBadRequest,
	ReportAlreadyRead:             http.StatusBadRequest,
	ReportAlreadyExist:            http.StatusBadRequest,
	AuthorizationError:            http.StatusForbidden,
	ClientRequired:                http.StatusForbidden,
	FreelancerRequired:            http.StatusForbidden,
	NotAuthenticatedError:         http.StatusUnauthorized,
	PreConditionError:             http.StatusPreconditionFailed,
	InvalidOfferStatus:            http.StatusBadRequest,
	AssociatedProposalNotFound:    http.StatusBadRequest,
	OfferAlreadyExist:             http.StatusBadRequest,
	HeadCountLimitExceeded:        http.StatusBadRequest,
	ProposalAlreadyInvited:        http.StatusBadRequest,
	ProposalAlreadyExist:          http.StatusBadRequest,
	JobAlreadyArchived:            http.StatusBadRequest,
	JobNotActive:                  http.StatusBadRequest,
	JobNotArchived:                http.StatusBadRequest,
	JobExpired:                    http.StatusBadRequest,
	NotEnoughBalance:              http.StatusBadRequest,
	JobTagAlreadyExist:            http.StatusBadRequest,
	ContributorAlreadyExist:       http.StatusBadRequest,
	ConflictOperation:             http.StatusConflict,
}

func InvalidHeaderError(name string) *AppError {
	appError := &AppError{
		ErrorCode: InvalidHeader,
		Message:   "invalid header",
		Target:    name,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewInvalidParameterError(name string) *AppError {
	appError := &AppError{
		ErrorCode: InvalidParameter,
		Message:   "The given parameter is missing or invalid",
		Target:    name,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewInvalidParameterWithMsgError(name string, msg string) *AppError {
	appError := &AppError{
		ErrorCode: InvalidParameter,
		Message:   msg,
		Target:    name,
	}

	appError.CaptureStackTrace()
	return appError
}

func ParameterTooLongError(name string) *AppError {
	appError := &AppError{
		ErrorCode: InvalidParameter,
		Message:   "too long",
		Target:    name,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewParameterRequiredError(name string) *AppError {
	appError := &AppError{
		ErrorCode: InvalidParameter,
		Message:   "The given parameter is required",
		Target:    name,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewInvalidRequestBodyError(cause error) *AppError {
	appError := &AppError{
		ErrorCode: InvalidRequestBody,
		Message:   "Request body is invalid",
		cause:     cause,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewResourceNotFoundError(target string) *AppError {
	appError := &AppError{
		ErrorCode: ResourceNotFound,
		Message:   "Resource not found",
		Target:    target,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewInternalServerError(msg string, cause error) *AppError {
	appError := &AppError{
		ErrorCode: InternalServerError,
		Message:   fmt.Sprintf("Internal server error: %s", msg),
		cause:     cause,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewDuplicateParameterError(name string, val interface{}) *AppError {
	appError := &AppError{
		ErrorCode: DuplicateParameter,
		Message:   fmt.Sprintf("the parameter is duplicated, name=%s, val=%v", name, val),
		Target:    name,
	}

	appError.CaptureStackTrace()
	return appError
}

func NewTooManyQueryOptionError(msg string) *AppError {
	appError := &AppError{
		ErrorCode: TooManyQueryOption,
		Message:   fmt.Sprintf("Too many query option parameters: %s", msg),
	}

	appError.CaptureStackTrace()
	return appError
}

func NewInvalidOperationError(reason string) *AppError {
	appError := &AppError{
		ErrorCode: InvalidOperation,
		Message:   fmt.Sprintf("invalid operation: %s", reason),
	}

	appError.CaptureStackTrace()
	return appError
}

func NewAuthorizationError() *AppError {
	appError := &AppError{
		ErrorCode: AuthorizationError,
		Message:   "authorization error",
	}

	appError.CaptureStackTrace()
	return appError
}

func NewNotAuthenticatedError() *AppError {
	appError := &AppError{
		ErrorCode: NotAuthenticatedError,
		Message:   "not authenticated error",
	}

	appError.CaptureStackTrace()
	return appError
}

func NewFreelancerAccountAlreadyExistError() *AppError {
	appError := &AppError{
		ErrorCode: FreelancerAccountAlreadyExist,
		Message:   "Freelancer account already exists.",
	}

	appError.CaptureStackTrace()
	return appError
}

func ItemSuspendedError() *AppError {
	appError := &AppError{
		ErrorCode: ItemSuspended,
		Message:   "Account is suspended",
	}
	appError.CaptureStackTrace()
	return appError
}
func ItemNotSuspendedError() *AppError {
	appError := &AppError{
		ErrorCode: ItemNotSuspended,
		Message:   "Account is not suspended",
	}
	appError.CaptureStackTrace()
	return appError
}

func OfferAlreadyAcceptedError() *AppError {
	appError := &AppError{
		ErrorCode: OfferAlreadyAccepted,
		Message:   "This offer is already accepted by the freelancer.",
	}
	appError.CaptureStackTrace()
	return appError
}

func ContractAlreadyExistError() *AppError {
	appError := &AppError{
		ErrorCode: ContractAlreadyExist,
		Message:   "Contract already exists.",
	}
	appError.CaptureStackTrace()
	return appError
}

func ReportImproperItemError(msg string) *AppError {
	appError := &AppError{
		ErrorCode: ReportImproperItem,
		Message:   fmt.Sprintf("Report improper item: %s", msg),
	}
	appError.CaptureStackTrace()
	return appError
}

func ReportAlreadyReadError() *AppError {
	appError := &AppError{
		ErrorCode: ReportAlreadyRead,
		Message:   "report has already been read",
	}
	appError.CaptureStackTrace()
	return appError
}

func ReportAlreadyExistError() *AppError {
	appError := &AppError{
		ErrorCode: ReportAlreadyExist,
		Message:   "Report already exist",
	}
	appError.CaptureStackTrace()
	return appError
}

func ClientRequiredError() *AppError {
	appError := &AppError{
		ErrorCode: ClientRequired,
		Message:   "client required",
	}
	appError.CaptureStackTrace()
	return appError
}

func FreelancerRequiredError() *AppError {
	appError := &AppError{
		ErrorCode: FreelancerRequired,
		Message:   "freelancer required",
	}
	appError.CaptureStackTrace()
	return appError
}

func ProposalAlreadyInvitedError() *AppError {
	appError := &AppError{
		ErrorCode: ProposalAlreadyInvited,
		Message:   "proposal is already invited",
	}
	appError.CaptureStackTrace()
	return appError
}

func ProposalAlreadyExistError() *AppError {
	appError := &AppError{
		ErrorCode: ProposalAlreadyExist,
		Message:   "proposal already exist",
	}
	appError.CaptureStackTrace()
	return appError
}

func OfferAlreadyExistError() *AppError {
	appError := &AppError{
		ErrorCode: OfferAlreadyExist,
		Message:   "offer is already exist",
	}
	appError.CaptureStackTrace()
	return appError
}

func InvalidOfferStatusError() *AppError {
	appError := &AppError{
		ErrorCode: InvalidOfferStatus,
		Message:   "invalid offer status",
	}
	appError.CaptureStackTrace()
	return appError
}

func HeadCountLimitExceededError() *AppError {
	appError := &AppError{
		ErrorCode: HeadCountLimitExceeded,
		Message:   "head count limit exceeded",
	}

	appError.CaptureStackTrace()
	return appError
}

func AssociatedProposalNotFoundError() *AppError {
	appError := &AppError{
		ErrorCode: AssociatedProposalNotFound,
		Message:   "associated proposal not found",
	}

	appError.CaptureStackTrace()
	return appError
}

func JobAlreadyArchivedError() *AppError {
	appError := &AppError{
		ErrorCode: JobAlreadyArchived,
		Message:   "job is already archived",
	}

	appError.CaptureStackTrace()
	return appError
}

func JobNotActiveError() *AppError {
	appError := &AppError{
		ErrorCode: JobNotActive,
		Message:   "job is inactive, can not archive it",
	}

	appError.CaptureStackTrace()
	return appError
}

func JobNotArchivedError() *AppError {
	appError := &AppError{
		ErrorCode: JobNotArchived,
		Message:   "job is not archived, can not unarchive it",
	}

	appError.CaptureStackTrace()
	return appError
}

func JobTagAlreadyExistError() *AppError {
	appError := &AppError{
		ErrorCode: JobTagAlreadyExist,
		Message:   "job tag already exists",
	}

	appError.CaptureStackTrace()
	return appError
}

func JobTagAlreadyDeletedError() *AppError {
	appError := &AppError{
		ErrorCode: JobTagAlreadyDeleted,
		Message:   "job tag already deleted",
	}

	appError.CaptureStackTrace()
	return appError
}

func ConcurrentUpdateError() *AppError {
	appError := &AppError{
		ErrorCode: ConcurrentUpdate,
		Message:   "concurrent update is not allowed",
	}

	appError.CaptureStackTrace()
	return appError
}

func JobExpiredError() *AppError {
	appError := &AppError{
		ErrorCode: JobExpired,
		Message:   "Job is expired",
	}

	appError.CaptureStackTrace()
	return appError
}

func NotEnoughBalanceError() *AppError {
	appError := &AppError{
		ErrorCode: NotEnoughBalance,
		Message:   "The credits are not enough to pay",
	}
	appError.CaptureStackTrace()
	return appError
}

func ContributorAlreadyExistError(val string) *AppError {
	appError := &AppError{
		ErrorCode: ContributorAlreadyExist,
		Message:   "Contributor already exists",
		Target:    val,
	}
	appError.CaptureStackTrace()
	return appError
}

func ConflictOperationError() *AppError {
	appError := &AppError{
		ErrorCode: ConflictOperation,
		Message:   "conflict operation",
	}
	appError.CaptureStackTrace()
	return appError
}
