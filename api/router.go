package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koizumi7010/blog-api/api/middlewares"
	"github.com/koizumi7010/blog-api/controllers"
	"github.com/koizumi7010/blog-api/services"
)

func NewRouter(db *sql.DB) *mux.Router {
	ser := services.NewMyAppService(db)
	aCon := controllers.NewArticleController(ser)
	cCon := controllers.NewCommentController(ser)

	r := mux.NewRouter()

	r.HandleFunc("/article", aCon.PostArticleHandler).Methods(http.MethodPost)
	r.HandleFunc("/article/list", aCon.ArticleListHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/{id:[0-9]+}", aCon.ArticleDetailHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/nice", aCon.PostNiceHandler).Methods(http.MethodPost)
	r.HandleFunc("/comment", cCon.PostCommentHandler).Methods(http.MethodPost)

	// 全てのハンドラに共通のミドルウェアを適用
	r.Use(middlewares.LoggingMiddleware)
	// r.Use(middlewares.AuthMiddleware) # 認証処理

	return r
}
