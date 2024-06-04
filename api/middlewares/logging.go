package middlewares

import (
	"log"
	"net/http"

	"github.com/koizumi7010/blog-api/common"
)

// 委譲によってHeader(),Write(),WriteHeader()のメソッドを持つ = ServeHTTP()の第一引数に渡せる
type resLoggingWriter struct {
	http.ResponseWriter
	code int
}

func NewResLoggingWriter(w http.ResponseWriter) *resLoggingWriter {
	return &resLoggingWriter{ResponseWriter: w, code: http.StatusOK}
}

// HandlerがHTTPレスポンスコードを書き込む時に使うメソッド
// WriteHeaderメソッドをオーバーライドして、レスポンスコードを記録する
func (rsw *resLoggingWriter) WriteHeader(code int) {
	rsw.code = code
	rsw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		traceID := newTraceID()

		// リクエスト情報をログ出力
		log.Printf("[%d]%s %s\n", traceID, req.RequestURI, req.Method)

		ctx := common.SetTraceID(req.Context(), traceID)
		req = req.WithContext(ctx)
		rlw := NewResLoggingWriter(w)

		next.ServeHTTP(rlw, req)

		log.Printf("[%d]res: %d", traceID, rlw.code)
	})
}
