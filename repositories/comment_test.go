package repositories_test

import (
	"testing"

	"github.com/koizumi7010/blog-api/models"
	"github.com/koizumi7010/blog-api/repositories"

	_ "github.com/go-sql-driver/mysql"
)

func TestInsetComment(t *testing.T) {
	comment := models.Comment{
		ArticleID: 1,
		Message:   "test",
	}

	expectedCommentID := 3
	newComment, err := repositories.InsertComment(testDB, comment)
	if err != nil {
		t.Error(err)
	}
	if newComment.CommentID != expectedCommentID {
		t.Errorf("CommentID: get %d, but want %d\n", newComment.CommentID, expectedCommentID)
	}

	// Clean up
	t.Cleanup(func() {
		const sqlStr = `
			delete from comments
			where article_id = ? and message = ?;
		`
		testDB.Exec(sqlStr, comment.ArticleID, comment.Message)
	})
}

func TestSelectCommentList(t *testing.T) {
	articleID := 1
	got, err := repositories.SelectCommentList(testDB, articleID)
	if err != nil {
		t.Error(err)
	}

	for _, comment := range got {
		if comment.ArticleID != articleID {
			t.Errorf("want comment of articleID %d but got ID %d\n", articleID, comment.ArticleID)
		}
	}
}
