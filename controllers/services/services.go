package services

import "github.com/koizumi7010/blog-api/models"

// /article関連を引き受けるインターフェース
type ArticleServicer interface {
	PostArticleService(article models.Article) (models.Article, error)
	GetArticleListService(page int) ([]models.Article, error)
	GetArticleService(articleID int) (models.Article, error)
	PostNiceService(article models.Article) (models.Article, error)
}

// /commentを引き受けるインターフェース
type CommentServicer interface {
	PostCommentService(comment models.Comment) (models.Comment, error)
}
