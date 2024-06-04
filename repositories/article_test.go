package repositories_test

import (
	"testing"

	"github.com/koizumi7010/blog-api/models"
	"github.com/koizumi7010/blog-api/repositories"
	"github.com/koizumi7010/blog-api/repositories/testdata"

	_ "github.com/go-sql-driver/mysql"
)

func TestSelectArticleList(t *testing.T) {
	expectedNum := len(testdata.ArticleTestData)
	articles, err := repositories.SelectArticleList(testDB, 1)
	if err != nil {
		t.Error(err)
	}
	if len(articles) != expectedNum {
		t.Errorf("get %d articles, but want %d articles\n", len(articles), expectedNum)
	}
}

func TestSelectArticleDetail(t *testing.T) {
	tests := []struct {
		testTitle string
		expected  models.Article
	}{
		{
			testTitle: "subtest1",
			expected:  testdata.ArticleTestData[0],
		},
		{
			testTitle: "subtest2",
			expected:  testdata.ArticleTestData[1],
		},
	}

	for _, test := range tests {
		t.Run(test.testTitle, func(t *testing.T) {
			got, err := repositories.SelectArticleDetail(testDB, test.expected.ID)
			if err != nil {
				t.Fatal(err)
			}

			if got.ID != test.expected.ID {
				t.Errorf("ID: get %d, but want %d\n", got.ID, test.expected.ID)
			}
			if got.Title != test.expected.Title {
				t.Errorf("Title: get %s, but want %s\n", got.Title, test.expected.Title)
			}
			if got.Contents != test.expected.Contents {
				t.Errorf("Contents: get %s, but want %s\n", got.Contents, test.expected.Contents)
			}
			if got.UserName != test.expected.UserName {
				t.Errorf("UserName: get %s, but want %s\n", got.UserName, test.expected.UserName)
			}
			if got.NiceNum != test.expected.NiceNum {
				t.Errorf("NiceNum: get %d, but want %d\n", got.NiceNum, test.expected.NiceNum)
			}
		})
	}
}

func TestInsertArticle(t *testing.T) {
	article := models.Article{
		Title:    "insertTest",
		Contents: "test",
		UserName: "saki",
	}

	expectedArticleNum := 3
	newArticle, err := repositories.InsertArticle(testDB, article)
	if err != nil {
		t.Error(err)
	}
	if newArticle.ID != expectedArticleNum {
		t.Errorf("ID: get %d, but want %d\n", newArticle.ID, expectedArticleNum)
	}

	// Clean up
	t.Cleanup(func() {
		const sqlStr = `
			delete from articles
			where title = ? and contents = ? and username = ?;
		`
		testDB.Exec(sqlStr, article.Title, article.Contents, article.UserName)
	})
}

func TestUpdateNiceNum(t *testing.T) {
	articleID := 1
	before, err := repositories.SelectArticleDetail(testDB, articleID)
	if err != nil {
		t.Error(err)
	}
	err = repositories.UpdateNiceNum(testDB, articleID)
	if err != nil {
		t.Error(err)
	}

	after, err := repositories.SelectArticleDetail(testDB, articleID)
	if err != nil {
		t.Error(err)
	}
	if before.NiceNum+1 != after.NiceNum {
		t.Error("failed to update nice_num")
	}
}
