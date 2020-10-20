package utils

import "transport/lib/errors"

//srms -5000 -> -5500
const (
	// Database errors
	BOTDatabaseError      errors.ErrorType = -5001 //BOT_DATABASE_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	BOTGetDataNotfound    errors.ErrorType = -5002 //BOT_GET_DATA_NOTFOUND - Không tìm thấy dữ liệu trên hệ thống, vui lòng kiểm tra lại!
	BOTGetDataError       errors.ErrorType = -5003 //BOT_ADD_DATA_ERROR - Không thể truy xuất dữ liệu, vui lòng kiểm tra lại!
	BOTAddDataError       errors.ErrorType = -5004 //BOT_ADD_DATA_ERROR - Thêm mới dữ liệu không thành công, vui lòng kiểm tra lại!
	BOTUpdateDataError    errors.ErrorType = -5005 //BOT_UPDATE_DATA_ERROR - Cập nhật dữ liệu không thành công, vui lòng kiểm tra lại!
	BOTUpdateSuccessError errors.ErrorType = -5006 //BOT_UPDATE_SUCCESS_ERROR - Không cho phép cập nhật message ở trạng thái thành công, vui lòng kiểm tra lại!
	BOTDeleteDataError    errors.ErrorType = -5007 //BOT_DELETE_DATA_ERROR - Xóa dữ liệu không thành công, vui lòng kiểm tra lại!
	BOTDataExistedError   errors.ErrorType = -5008 //BOT_DELETE_DATA_ERROR - Dữ liệu đã tồn tại trên hệ thống, vui lòng kiểm tra lại!

	// Logic errors
	BOTDataInvalidError errors.ErrorType = -5020 //BOT_DATA_INVALID_ERROR - Dữ liệu không hợp lệ, vui lòng kiểm tra lại!
	BOTSerializeError   errors.ErrorType = -5021 //BOT_SERIALIZE_DATA_ERROR - Dữ liệu không hợp lệ, vui lòng kiểm tra lại!
	BOTParseQueryError  errors.ErrorType = -5022 //BOT_PARSE_QUERY_ERROR - Dữ liệu Query truyển vào không hợp lệ, vui lòng kiểm tra lại!
	BOTParseTimeError   errors.ErrorType = -5023 //BOT_PARSE_TIME_ERROR - Xảy ra lỗi khi parse thời gian, vui lòng kiểm tra lại!

	// Business errors
	BOTCallAPIError errors.ErrorType = -5050 //BOT_CALL_API_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
)

var DefaultErrorMap = errors.ErrorMap{
	errors.Success: {
		ErrCodeStr: "SUCCESS",
		Message:    "Thành công",
	},
	errors.NotFound: {
		ErrCodeStr: "NOT_FOUND",
		Message:    "Không tìm thấy tài nguyên",
	},
	errors.BadRequest: {
		ErrCodeStr: "BAD_REQUEST",
		Message:    "Dữ liệu gửi lên không đúng, vui lòng kiểm tra lại"},
	errors.Unauthorized: {
		ErrCodeStr: "UNAUTHORIZED",
		Message:    "Không có quyền truy cập, vui lòng liên hệ quản trị viên để biết thêm thông tin!",
	},
	BOTDatabaseError: {
		ErrCodeStr: "BOT_DATABASE_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	BOTGetDataNotfound: {
		ErrCodeStr: "BOT_GET_DATA_NOTFOUND",
		Message:    "Không tìm thấy dữ liệu trên hệ thống, vui lòng kiểm tra lại!",
	},
	BOTGetDataError: {
		ErrCodeStr: "BOT_ADD_DATA_ERROR",
		Message:    "Không thể truy xuất dữ liệu, vui lòng kiểm tra lại!",
	},
	BOTAddDataError: {
		ErrCodeStr: "BOT_ADD_DATA_ERROR",
		Message:    "Thêm mới dữ liệu không thành công, vui lòng kiểm tra lại!",
	},
	BOTUpdateDataError: {
		ErrCodeStr: "BOT_UPDATE_DATA_ERROR",
		Message:    "Cập nhật dữ liệu không thành công, vui lòng kiểm tra lại!",
	},
	BOTUpdateSuccessError: {
		ErrCodeStr: "BOT_UPDATE_SUCCESS_ERROR",
		Message:    "Không cho phép cập nhật message ở trạng thái thành công, vui lòng kiểm tra lại!",
	},
	BOTDeleteDataError: {
		ErrCodeStr: "BOT_DELETE_DATA_ERROR",
		Message:    "Xóa dữ liệu không thành công, vui lòng kiểm tra lại!",
	},
	BOTDataExistedError: {
		ErrCodeStr: "BOT_DELETE_DATA_ERROR",
		Message:    "Dữ liệu đã tồn tại trên hệ thống, vui lòng kiểm tra lại!",
	},
	BOTDataInvalidError: {
		ErrCodeStr: "BOT_DATA_INVALID_ERROR",
		Message:    "Dữ liệu không hợp lệ, vui lòng kiểm tra lại!",
	},
	BOTSerializeError: {
		ErrCodeStr: "BOT_SERIALIZE_DATA_ERROR",
		Message:    "Dữ liệu không hợp lệ, vui lòng kiểm tra lại!",
	},
	BOTParseQueryError: {
		ErrCodeStr: "BOT_PARSE_QUERY_ERROR",
		Message:    "Dữ liệu Query truyển vào không hợp lệ, vui lòng kiểm tra lại!",
	},
	BOTParseTimeError: {
		ErrCodeStr: "BOT_PARSE_TIME_ERROR",
		Message:    "Xảy ra lỗi khi parse thời gian, vui lòng kiểm tra lại!",
	},
	BOTCallAPIError: {
		ErrCodeStr: "BOT_CALL_API_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
}
