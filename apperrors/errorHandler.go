package apperrors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/koizumi7010/blog-api/common"
)

func ErrorHandler(w http.ResponseWriter, req *http.Request, err error) {
	// 変換先であるMyAppError型の変数を宣言
	var appErr *MyAppError
	// errors.As関数で引数のerrをMyAppError型のappErrに変換
	if !errors.As(err, &appErr) {
		// 変換できなかった場合は、Unknownエラーを返す
		appErr = &MyAppError{
			ErrCode: Unknown,
			Message: "internal process failed",
			Err:     err,
		}
	}

	traceID := common.GetTraceID(req.Context())
	log.Printf("[%d]error: %s\n", traceID, appErr)

	// ユーザーに返すHTTPステータスコードを収めるための変数
	var statusCode int

	switch appErr.ErrCode {
	case NAData:
		statusCode = http.StatusNotFound
	case NoTargetData, ReqBodyDecodeFailed, BadParam:
		statusCode = http.StatusBadRequest
	case RequiredAuthorizationHeader, CannotMakeValidator, Unauthorizated:
		statusCode = http.StatusUnauthorized
	case NotMatchUser:
		statusCode = http.StatusForbidden
	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(appErr)
}
