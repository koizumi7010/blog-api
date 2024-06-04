package services

import (
	"database/sql"
	"errors"

	"github.com/koizumi7010/blog-api/apperrors"
	"github.com/koizumi7010/blog-api/models"
	"github.com/koizumi7010/blog-api/repositories"
)

func (s *MyAppService) PostArticleService(article models.Article) (models.Article, error) {
	// repositories層のInsertArticleを呼び出す
	newArticle, err := repositories.InsertArticle(s.db, article)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "failed to insert article")
		return models.Article{}, err
	}

	return newArticle, nil
}

func (s *MyAppService) GetArticleListService(page int) ([]models.Article, error) {
	// repositories層のSelectArticleListを呼び出す
	articles, err := repositories.SelectArticleList(s.db, page)
	if err != nil {
		err = apperrors.GetDataFailed.Wrap(err, "fail to get data")
		return nil, err
	}

	if len(articles) == 0 {
		err = apperrors.NAData.Wrap(ErrNoData, "no data")
		return nil, err
	}

	return articles, nil
}

func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	var article models.Article
	var commentList []models.Comment
	var articleGetErr, commentGetErr error

	type articleResult struct {
		article models.Article
		err     error
	}
	articleCh := make(chan articleResult)
	defer close(articleCh)

	go func(ch chan<- articleResult, db *sql.DB, articleID int) {
		article, err := repositories.SelectArticleDetail(s.db, articleID)
		ch <- articleResult{article: article, err: err}
	}(articleCh, s.db, articleID)

	type commentResult struct {
		commentList []models.Comment
		err         error
	}
	commentCh := make(chan commentResult)
	defer close(commentCh)

	go func(ch chan<- commentResult, db *sql.DB, articleID int) {
		commentList, err := repositories.SelectCommentList(s.db, articleID)
		ch <- commentResult{commentList: commentList, err: err}
	}(commentCh, s.db, articleID)

	for i := 0; i < 2; i++ {
		select {
		case ar := <-articleCh:
			article, articleGetErr = ar.article, ar.err
		case cr := <-commentCh:
			commentList, commentGetErr = cr.commentList, cr.err
		}
	}

	if articleGetErr != nil {
		if errors.Is(articleGetErr, sql.ErrNoRows) {
			err := apperrors.NAData.Wrap(articleGetErr, "no data")
			return models.Article{}, err
		}
		err := apperrors.GetDataFailed.Wrap(articleGetErr, "fail to get data")
		return models.Article{}, err
	}

	if commentGetErr != nil {
		err := apperrors.GetDataFailed.Wrap(commentGetErr, "fail to get data")
		return models.Article{}, err
	}

	article.CommentList = append(article.CommentList, commentList...)

	return article, nil
}

func (s *MyAppService) PostNiceService(article models.Article) (models.Article, error) {
	// repositories層のUpdateNiceNumを呼び出す
	err := repositories.UpdateNiceNum(s.db, article.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperrors.NoTargetData.Wrap(err, "does not exist target article")
			return models.Article{}, err
		}
		err = apperrors.UpdateDataFailed.Wrap(err, "fail to update nice count")
		return models.Article{}, err
	}

	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Contents:  article.Contents,
		UserName:  article.UserName,
		NiceNum:   article.NiceNum + 1,
		CreatedAt: article.CreatedAt,
	}, nil
}
