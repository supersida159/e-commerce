package common

import (
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	// Bad Request Errors (400)
	CodeInvalidMessageKey ErrorCode = -100 + iota
	CodeUsernameTooShort
	CodePasswordMustContainNumber
	CodeInvalidEmail
	CodeUserNotExist

	CodeFailedToSavePhoto
	CodeInvalidDateFormat
	CodeJSONBindingError
	// This one for validator
	CodeMissingRequiredField
	CodeFieldBelowMinimum
	CodeFieldAboveMaximum
	CodeInvalidURL

	// Already Exist Errors (405)
	CodePhoneNumberExists ErrorCode = -110 + iota
	CodeUserNameExists
	CodeProductSKUExists
	CodeEmailExists

	// Not Found Errors (404)
	CodeEntityNotExist ErrorCode = -111 + iota
	CodeResourceNotFound
	CodePageNotFound
	CodeQueryNotFound
	CodeUpdateFailed
	CodeProductNotFound
	CodeOrderNotFound

	// Unauthorized Errors (401)
	CodeUnauthorizedAccess ErrorCode = -201 + iota
	CodeInvalidCredentials
	CodeInvalidToken
	CodeExpiredToken
	CodeInvalidRefreshToken
	CodeExpiredRefreshToken
	CodeWrongPasswordOrUsername
	CodeTokenLeaked
	CodeTokenBlocked
	CodeTokenMustBeUpdated

	// Forbidden Errors (403)
	CodeForbidden ErrorCode = -202 + iota
	CodeAccessDenied

	// Conflict Errors (409)
	CodeResourceConflict ErrorCode = -302 + iota
	CodeDuplicateEntry
	CodeOutOfStock

	// Unprocessable Entity Errors (422)
	CodeUnprocessableEntity ErrorCode = -303 + iota
	CodeInvalidInputData
	CodeInvalidPaymentInfo

	// Internal Server Errors (500)
	CodeInternalServerError ErrorCode = -401 + iota
	CodeDatabaseConnectionError
	CodeFailedToConvertToDTO
	CodeJobRetryError
	CodePaymentGatewayError
	CodeErrorDB

	// Not Implemented Errors (501)
	CodeNotImplemented ErrorCode = -402 + iota
	CodeFeatureNotAvailable

	// Service Unavailable Errors (503)
	CodeServiceUnavailable ErrorCode = -403 + iota
	CodeSystemMaintenance

	// Custom Errors (400)
	CodeInvalidRequestParameter ErrorCode = -600 + iota
	CodeInvalidQuantity
	CodeInvalidPrice
	CodeInvalidDiscountAmount
	CodeInvalidCouponCode
	CodeExpiredCoupon
	CodeInvalidShippingAddress
)

type AppError struct {
	StatusCode   int       `json:"status_code"`
	RootErr      error     `json:"-"`
	MessageEn    string    `json:"message_en"`
	MessageVi    string    `json:"message_vi"`
	Key          string    `json:"error_key"`
	InternalCode ErrorCode `json:"code"`
}

// Error implements error.
func (a *AppError) Error() string {
	panic("unimplemented")
}

func NewErrorResponse(rootErr error, messageEn, messageVi, key string, internalStatusCode ErrorCode, StatusCode int) *AppError {
	return &AppError{
		RootErr:      rootErr,
		MessageEn:    messageEn,
		MessageVi:    messageVi,
		Key:          key,
		StatusCode:   StatusCode,
		InternalCode: internalStatusCode,
	}
}

// Bad Request Errors (400)
func ErrInvalidMessageKey(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid message key",
		"Message Key không hợp lệ",
		"INVALID_MESSAGE_KEY",
		CodeInvalidMessageKey,
		http.StatusBadRequest,
	)
}

func ErrUsernameTooShort(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Username must be at least 8 characters",
		"Tên người dùng phải có ít nhất 8 ký tự",
		"USERNAME_TOO_SHORT",
		CodeUsernameTooShort,
		http.StatusBadRequest,
	)
}

func ErrPasswordMustContainNumber(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Password must contain at least one number",
		"Mật khẩu phải chứa ít nhất một số",
		"PASSWORD_MUST_CONTAIN_NUMBER",
		CodePasswordMustContainNumber,
		http.StatusBadRequest,
	)
}

func ErrInvalidEmail(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Email is invalid",
		"Email không hợp lệ",
		"INVALID_EMAIL",
		CodeInvalidEmail,
		http.StatusBadRequest,
	)
}

func ErrUserNotExist(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"User does not exist",
		"Người dùng không tồn tại",
		"USER_NOT_EXIST",
		CodeUserNotExist,
		http.StatusBadRequest,
	)
}

func ErrMissingRequiredField(field string, err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		fmt.Sprintf("Missing required field %s", field),
		fmt.Sprintf("Thiếu trường bắt buộc %s", field),
		"MISSING_REQUIRED_FIELD",
		CodeMissingRequiredField,
		http.StatusBadRequest,
	)
}

func ErrFailedToSavePhoto(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Failed to save photo",
		"Upload ảnh thất bại",
		"FAILED_TO_SAVE_PHOTO",
		CodeFailedToSavePhoto,
		http.StatusBadRequest,
	)
}

func ErrJSONBlindding(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid json format",
		"Định dạng json không hợp lệ",
		"INVALID_JSON_FORMAT",
		CodeJSONBindingError,
		http.StatusBadRequest,
	)
}

// Already Exist Errors (405)
func ErrPhoneNumberExists(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Phone number already exists",
		"Số điện thoại đã tồn tại",
		"PHONE_NUMBER_EXISTS",
		CodePhoneNumberExists,
		http.StatusMethodNotAllowed,
	)
}

func ErrUserNameExists(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"User name already exists",
		"Tên người dùng đã tồn tại",
		"USER_NAME_EXISTS",
		CodeUserNameExists,
		http.StatusMethodNotAllowed,
	)
}

func ErrProductSKUExists(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Product SKU already exists",
		"Mã SKU sản phẩm đã tồn tại",
		"PRODUCT_SKU_EXISTS",
		CodeProductSKUExists,
		http.StatusMethodNotAllowed,
	)
}

func ErrEmailExists(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Email already exists",
		"Email đã tồn tại",
		"EMAIL_EXISTS",
		CodeEmailExists,
		http.StatusMethodNotAllowed,
	)
}

// Not Found Errors (404)
func ErrEntityNotExist(entity string, err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		fmt.Sprintf("Entity %s does not exist", entity),
		fmt.Sprintf("Dữ liệu %s không tồn tại", entity),
		"ENTITY_NOT_EXIST",
		CodeEntityNotExist,
		http.StatusNotFound,
	)
}

func ErrResourceNotFound(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Resource not found",
		"Tài nguyên không tìm thấy",
		"RESOURCE_NOT_FOUND",
		CodeResourceNotFound,
		http.StatusNotFound,
	)
}

func ErrProductNotFound(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Product not found",
		"Không tìm thấy sản phẩm",
		"PRODUCT_NOT_FOUND",
		CodeProductNotFound,
		http.StatusNotFound,
	)
}

func ErrOrderNotFound(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Order not found",
		"Không tìm thấy đơn hàng",
		"ORDER_NOT_FOUND",
		CodeOrderNotFound,
		http.StatusNotFound,
	)
}

// Unauthorized Errors (401)
func ErrUnauthorizedAccess(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Unauthorized access",
		"Truy cập không được phép",
		"UNAUTHORIZED_ACCESS",
		CodeUnauthorizedAccess,
		http.StatusUnauthorized,
	)
}

func ErrInvalidCredentials(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid credentials",
		"Thông tin đăng nhập không hợp lệ",
		"INVALID_CREDENTIALS",
		CodeInvalidCredentials,
		http.StatusUnauthorized,
	)
}

// Forbidden Errors (403)
func ErrForbidden(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Forbidden: You do not have permission",
		"Cấm: Bạn không có quyền truy cập",
		"FORBIDDEN",
		CodeForbidden,
		http.StatusForbidden,
	)
}

func ErrAccessDenied(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Access denied",
		"Truy cập bị từ chối",
		"ACCESS_DENIED",
		CodeAccessDenied,
		http.StatusForbidden,
	)
}

// Conflict Errors (409)
func ErrResourceConflict(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Conflict: resource already exists",
		"Xung đột: Tài nguyên đã tồn tại",
		"RESOURCE_CONFLICT",
		CodeResourceConflict,
		http.StatusConflict,
	)
}

func ErrDuplicateEntry(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Conflict: duplicate entry",
		"Xung đột: Mục trùng lặp",
		"DUPLICATE_ENTRY",
		CodeDuplicateEntry,
		http.StatusConflict,
	)
}

func ErrOutOfStock(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Product is out of stock",
		"Sản phẩm đã hết hàng",
		"OUT_OF_STOCK",
		CodeOutOfStock,
		http.StatusConflict,
	)
}

// Unprocessable Entity Errors (422)
func ErrUnprocessableEntity(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Unprocessable Entity",
		"Thực thể không thể xử lý",
		"UNPROCESSABLE_ENTITY",
		CodeUnprocessableEntity,
		http.StatusUnprocessableEntity,
	)
}

func ErrInvalidInputData(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid input data",
		"Dữ liệu đầu vào không hợp lệ",
		"INVALID_INPUT_DATA",
		CodeInvalidInputData,
		http.StatusUnprocessableEntity,
	)
}

func ErrInvalidPaymentInfo(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid payment information",
		"Thông tin thanh toán không hợp lệ",
		"INVALID_PAYMENT_INFO",
		CodeInvalidPaymentInfo,
		http.StatusUnprocessableEntity,
	)
}

// Internal Server Errors (500)
func ErrInternalServerError(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Internal server error",
		"Lỗi máy chủ nội bộ",
		"INTERNAL_SERVER_ERROR",
		CodeInternalServerError,
		http.StatusInternalServerError,
	)
}

func ErrDatabaseConnectionError(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Database connection error",
		"Lỗi kết nối cơ sở dữ liệu",
		"DATABASE_CONNECTION_ERROR",
		CodeDatabaseConnectionError,
		http.StatusInternalServerError,
	)
}

func ErrFailedToConvertToDTO(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Failed to convert to ResponseDTO",
		"Không thể chuyển đổi sang ResponseDTO",
		"FAILED_TO_CONVERT_TO_DTO",
		CodeFailedToConvertToDTO,
		http.StatusInternalServerError,
	)
}

func ErrPaymentGatewayError(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Payment gateway error",
		"Lỗi cổng thanh toán",
		"PAYMENT_GATEWAY_ERROR",
		CodePaymentGatewayError,
		http.StatusInternalServerError,
	)
}

// Not Implemented Errors (501)
func ErrNotImplemented(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Not implemented",
		"Chức năng chưa được triển khai",
		"NOT_IMPLEMENTED",
		CodeNotImplemented,
		http.StatusNotImplemented,
	)
}

func ErrFeatureNotAvailable(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Feature not available",
		"Tính năng chưa có sẵn",
		"FEATURE_NOT_AVAILABLE",
		CodeFeatureNotAvailable,
		http.StatusNotImplemented,
	)
}

// Service Unavailable Errors (503)
func ErrServiceUnavailable(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Service unavailable",
		"Dịch vụ không khả dụng",
		"SERVICE_UNAVAILABLE",
		CodeServiceUnavailable,
		http.StatusServiceUnavailable,
	)
}

func ErrSystemMaintenance(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"System maintenance",
		"Hệ thống đang bảo trì",
		"SYSTEM_MAINTENANCE",
		CodeSystemMaintenance,
		http.StatusServiceUnavailable,
	)
}

// Custom Errors (400)
func ErrInvalidRequestParameter(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid request parameter",
		"Tham số yêu cầu không hợp lệ",
		"INVALID_REQUEST_PARAMETER",
		CodeInvalidRequestParameter,
		http.StatusBadRequest,
	)
}

func ErrInvalidQuantity(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid quantity",
		"Số lượng không hợp lệ",
		"INVALID_QUANTITY",
		CodeInvalidQuantity,
		http.StatusBadRequest,
	)
}

func ErrInvalidPrice(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid price",
		"Giá không hợp lệ",
		"INVALID_PRICE",
		CodeInvalidPrice,
		http.StatusBadRequest,
	)
}

func ErrInvalidDiscountAmount(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid discount amount",
		"Số tiền giảm giá không hợp lệ",
		"INVALID_DISCOUNT_AMOUNT",
		CodeInvalidDiscountAmount,
		http.StatusBadRequest,
	)
}

func ErrInvalidCouponCode(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid coupon code",
		"Mã giảm giá không hợp lệ",
		"INVALID_COUPON_CODE",
		CodeInvalidCouponCode,
		http.StatusBadRequest,
	)
}

func ErrExpiredCoupon(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Coupon has expired",
		"Mã giảm giá đã hết hạn",
		"EXPIRED_COUPON",
		CodeExpiredCoupon,
		http.StatusBadRequest,
	)
}

func ErrInvalidShippingAddress(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid shipping address",
		"Địa chỉ giao hàng không hợp lệ",
		"INVALID_SHIPPING_ADDRESS",
		CodeInvalidShippingAddress,
		http.StatusBadRequest,
	)
}

func ErrInvalidToken(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid token",
		"Token không hợp lệ",
		"INVALID_TOKEN",
		CodeInvalidToken,
		http.StatusUnauthorized,
	)
}

func ErrTokenExpired(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Token expired",
		"Token hết hạn",
		"TOKEN_EXPIRED",
		CodeExpiredToken,
		http.StatusUnauthorized,
	)
}

func ErrRefreshTokenExpired(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Refresh token expired",
		"Refresh token hết hạn",
		"REFRESH_TOKEN_EXPIRED",
		CodeExpiredRefreshToken,
		http.StatusUnauthorized,
	)
}

func ErrWrongPasswordOrUsername(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Wrong password or username",
		"Mật khẩu hoặc tên người dùng không đúng",
		"WRONG_PASSWORD_OR_USERNAME",
		CodeWrongPasswordOrUsername,
		http.StatusUnauthorized,
	)
}

func ErrFieldBelowMinimum(fieldName string, min interface{}, err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		fmt.Sprintf("%s must be at least %v", fieldName, min),
		fmt.Sprintf("%s phải ít nhất là %v", fieldName, min),
		"FIELD_BELOW_MINIMUM",
		CodeFieldBelowMinimum,
		http.StatusBadRequest,
	)
}

func ErrFieldAboveMaximum(fieldName string, max interface{}, err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		fmt.Sprintf("%s must not exceed %v", fieldName, max),
		fmt.Sprintf("%s không được vượt quá %v", fieldName, max),
		"FIELD_ABOVE_MAXIMUM",
		CodeFieldAboveMaximum,
		http.StatusBadRequest,
	)
}

func ErrInvalidURL(fieldName string, err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		fmt.Sprintf("%s must be a valid URL", fieldName),
		fmt.Sprintf("%s phải là một URL hợp lệ", fieldName),
		"INVALID_URL",
		CodeInvalidURL,
		http.StatusBadRequest,
	)
}

func ErrInvalidDateTime(fieldName string, format string, err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		fmt.Sprintf("%s must be a valid date and time in the format %s", fieldName, format),
		fmt.Sprintf("%s phải là ngày giờ hợp lệ theo định dạng %s", fieldName, format),
		"INVALID_DATETIME",
		CodeInvalidDateFormat,
		http.StatusBadRequest,
	)
}

func ErrTokenLeaked(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Token has been compromised and is no longer valid",
		"Token đã bị xâm phạm và không còn hợp lệ",
		"TOKEN_LEAKED",
		CodeTokenLeaked,
		http.StatusUnauthorized,
	)
}

func ErrTokenBlocked(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Token has been blocked and is no longer valid",
		"Token đã bị chặn và không còn hợp lệ",
		"TOKEN_BLOCKED",
		CodeTokenBlocked,
		http.StatusUnauthorized,
	)
}

func ErrTokenMustBeUpdated(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Token must be updated",
		"Token cần được cập nhật",
		"TOKEN_MUST_BE_UPDATED",
		CodeTokenMustBeUpdated,
		http.StatusUnauthorized,
	)
}

func ErrJobRetry(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Job retry error occurred",
		"Đã xảy ra lỗi khi thử lại công việc",
		"JOB_RETRY_ERROR",
		CodeJobRetryError,
		http.StatusInternalServerError,
	)
}

func ErrDB(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Database error occurred",
		"Đã xảy ra lỗi cơ sở dữ liệu",
		"DATABASE_ERROR",
		CodeErrorDB,
		http.StatusInternalServerError,
	)
}
