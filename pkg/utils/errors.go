package utils

import "transport/lib/errors"

//srms -5000 -> -5500
const (
	// Database errors
	EMSDatabaseError      errors.ErrorType = -5001 //EMS_DATABASE_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSGetDataNotfound    errors.ErrorType = -5002 //EMS_GET_DATA_NOTFOUND - Không tìm thấy dữ liệu trên hệ thống, vui lòng kiểm tra lại!
	EMSGetDataError       errors.ErrorType = -5003 //EMS_ADD_DATA_ERROR - Không thể truy xuất dữ liệu, vui lòng kiểm tra lại!
	EMSAddDataError       errors.ErrorType = -5004 //EMS_ADD_DATA_ERROR - Thêm mới dữ liệu không thành công, vui lòng kiểm tra lại!
	EMSUpdateDataError    errors.ErrorType = -5005 //EMS_UPDATE_DATA_ERROR - Cập nhật dữ liệu không thành công, vui lòng kiểm tra lại!
	EMSUpdateSuccessError errors.ErrorType = -5006 //EMS_UPDATE_SUCCESS_ERROR - Không cho phép cập nhật message ở trạng thái thành công, vui lòng kiểm tra lại!
	EMSDeleteDataError    errors.ErrorType = -5007 //EMS_DELETE_DATA_ERROR - Xóa dữ liệu không thành công, vui lòng kiểm tra lại!
	EMSDataExistedError   errors.ErrorType = -5008 //EMS_DELETE_DATA_ERROR - Dữ liệu đã tồn tại trên hệ thống, vui lòng kiểm tra lại!

	// Logic errors
	EMSDataInvalidError errors.ErrorType = -5020 //EMS_DATA_INVALID_ERROR - Dữ liệu không hợp lệ, vui lòng kiểm tra lại!
	EMSSerializeError   errors.ErrorType = -5021 //EMS_SERIALIZE_DATA_ERROR - Dữ liệu không hợp lệ, vui lòng kiểm tra lại!
	EMSParseQueryError  errors.ErrorType = -5022 //EMS_PARSE_QUERY_ERROR - Dữ liệu Query truyển vào không hợp lệ, vui lòng kiểm tra lại!
	EMSParseTimeError   errors.ErrorType = -5023 //EMS_PARSE_TIME_ERROR - Xảy ra lỗi khi parse thời gian, vui lòng kiểm tra lại!

	// Business errors
	EMSCallAPIError      errors.ErrorType = -5050 //EMS_CALL_API_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSWaitPreviousError errors.ErrorType = -5051 //EMS_WAIT_PREVIOUS_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!

	// AMQP errors
	EMSAMQPOpenConnectionError   errors.ErrorType = -5070 //EMS_AMQP_OPEN_CONNECTION_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPCloseConnectionError  errors.ErrorType = -5071 //EMS_AMQP_CLOSE_CONNECTION_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPConnectionClosedError errors.ErrorType = -5072 //EMS_AMQP_CONNECTION_CLOSED_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPOpenChannelError      errors.ErrorType = -5073 //EMS_AMQP_OPEN_CHANNEL_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPCloseChannelError     errors.ErrorType = -5074 //EMS_AMQP_CLOSE_CHANNEL_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPChannelClosedError    errors.ErrorType = -5075 //EMS_AMQP_CHANNEL_CLOSED_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPExchangeDeclareError  errors.ErrorType = -5076 //EMS_AMQP_EXCHANGE_DECLARE_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPQueueDeclareError     errors.ErrorType = -5077 //EMS_AMQP_QUEUE_DECLARE_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPQueueBindError        errors.ErrorType = -5078 //EMS_AMQP_QUEUE_BIND_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!

	// Consumer errors
	EMSAMQPConsumerSubscribeError errors.ErrorType = -5100 //EMS_AMQP_CONSUMER_SUBSCRIBE_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPConsumerQosError       errors.ErrorType = -5101 //EMS_AMQP_CONSUMER_QOS_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPConsumerConsumeError   errors.ErrorType = -5102 //EMS_AMQP_CONSUMER_CONSUME_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!

	// Publisher errors
	EMSAMQPublisherConfirmError errors.ErrorType = -5150 //EMS_AMQP_PUBLISHER_CONFIRM_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
	EMSAMQPublisherPublishError errors.ErrorType = -5151 //EMS_AMQP_PUBLISHER_PUBLISH_ERROR - Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!
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
	EMSDatabaseError: {
		ErrCodeStr: "EMS_DATABASE_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSGetDataNotfound: {
		ErrCodeStr: "EMS_GET_DATA_NOTFOUND",
		Message:    "Không tìm thấy dữ liệu trên hệ thống, vui lòng kiểm tra lại!",
	},
	EMSGetDataError: {
		ErrCodeStr: "EMS_ADD_DATA_ERROR",
		Message:    "Không thể truy xuất dữ liệu, vui lòng kiểm tra lại!",
	},
	EMSAddDataError: {
		ErrCodeStr: "EMS_ADD_DATA_ERROR",
		Message:    "Thêm mới dữ liệu không thành công, vui lòng kiểm tra lại!",
	},
	EMSUpdateDataError: {
		ErrCodeStr: "EMS_UPDATE_DATA_ERROR",
		Message:    "Cập nhật dữ liệu không thành công, vui lòng kiểm tra lại!",
	},
	EMSUpdateSuccessError: {
		ErrCodeStr: "EMS_UPDATE_SUCCESS_ERROR",
		Message:    "Không cho phép cập nhật message ở trạng thái thành công, vui lòng kiểm tra lại!",
	},
	EMSDeleteDataError: {
		ErrCodeStr: "EMS_DELETE_DATA_ERROR",
		Message:    "Xóa dữ liệu không thành công, vui lòng kiểm tra lại!",
	},
	EMSDataExistedError: {
		ErrCodeStr: "EMS_DELETE_DATA_ERROR",
		Message:    "Dữ liệu đã tồn tại trên hệ thống, vui lòng kiểm tra lại!",
	},
	EMSDataInvalidError: {
		ErrCodeStr: "EMS_DATA_INVALID_ERROR",
		Message:    "Dữ liệu không hợp lệ, vui lòng kiểm tra lại!",
	},
	EMSSerializeError: {
		ErrCodeStr: "EMS_SERIALIZE_DATA_ERROR",
		Message:    "Dữ liệu không hợp lệ, vui lòng kiểm tra lại!",
	},
	EMSParseQueryError: {
		ErrCodeStr: "EMS_PARSE_QUERY_ERROR",
		Message:    "Dữ liệu Query truyển vào không hợp lệ, vui lòng kiểm tra lại!",
	},
	EMSParseTimeError: {
		ErrCodeStr: "EMS_PARSE_TIME_ERROR",
		Message:    "Xảy ra lỗi khi parse thời gian, vui lòng kiểm tra lại!",
	},
	EMSCallAPIError: {
		ErrCodeStr: "EMS_CALL_API_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPOpenConnectionError: {
		ErrCodeStr: "EMS_AMQP_OPEN_CONNECTION_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPCloseConnectionError: {
		ErrCodeStr: "EMS_AMQP_CLOSE_CONNECTION_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPConnectionClosedError: {
		ErrCodeStr: "EMS_AMQP_CONNECTION_CLOSED_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPOpenChannelError: {
		ErrCodeStr: "EMS_AMQP_OPEN_CHANNEL_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPCloseChannelError: {
		ErrCodeStr: "EMS_AMQP_CLOSE_CHANNEL_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPChannelClosedError: {
		ErrCodeStr: "EMS_AMQP_CHANNEL_CLOSED_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPExchangeDeclareError: {
		ErrCodeStr: "EMS_AMQP_EXCHANGE_DECLARE_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPQueueDeclareError: {
		ErrCodeStr: "EMS_AMQP_QUEUE_DECLARE_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPQueueBindError: {
		ErrCodeStr: "EMS_AMQP_QUEUE_BIND_ERROR",
		Message:    "Đã xảy ra lỗi khi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPConsumerSubscribeError: {
		ErrCodeStr: "EMS_AMQP_CONSUMER_SUBSCRIBE_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPConsumerQosError: {
		ErrCodeStr: "EMS_AMQP_CONSUMER_QOS_ERROR",
		Message:    "Đã xảy ra lỗi, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPConsumerConsumeError: {
		ErrCodeStr: "EMS_AMQP_CONSUMER_CONSUME_ERROR",
		Message:    "Đã xảy ra lỗi khi consume message, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPublisherConfirmError: {
		ErrCodeStr: "EMS_AMQP_PUBLISHER_CONFIRM_ERROR",
		Message:    "Đã xảy ra lỗi khi publish message, vui lòng liên hệ với quản trị viên!",
	},
	EMSAMQPublisherPublishError: {
		ErrCodeStr: "EMS_AMQP_PUBLISHER_PUBLISH_ERROR",
		Message:    "Đã xảy ra lỗi khi publish message, vui lòng liên hệ với quản trị viên!",
	},
}

var ErrorWaitPrevious = EMSWaitPreviousError.New()
