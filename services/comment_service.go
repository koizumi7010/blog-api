package services

import (
	"github.com/koizumi7010/blog-api/apperrors"
	"github.com/koizumi7010/blog-api/models"
	"github.com/koizumi7010/blog-api/repositories"
)

func (s *MyAppService) PostCommentService(comment models.Comment) (models.Comment, error) {
	// repositories層のInsertCommentを呼び出す
	newComment, err := repositories.InsertComment(s.db, comment)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "failed to insert comment")
		return models.Comment{}, err
	}

	return newComment, nil
}
