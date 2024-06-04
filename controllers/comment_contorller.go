package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/koizumi7010/blog-api/apperrors"
	"github.com/koizumi7010/blog-api/controllers/services"
	"github.com/koizumi7010/blog-api/models"
)

type CommentController struct {
	service services.CommentServicer
}

func NewCommentController(s services.CommentServicer) *CommentController {
	return &CommentController{service: s}
}

// POST /comment
func (c *CommentController) PostCommentHandler(w http.ResponseWriter, req *http.Request) {
	var reqComment models.Comment
	if err := json.NewDecoder(req.Body).Decode(&reqComment); err != nil {
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "bad body request")
		apperrors.ErrorHandler(w, req, err)
		return
	}

	comment, err := c.service.PostCommentService(reqComment)
	if err != nil {
		apperrors.ErrorHandler(w, req, err)
		return
	}

	json.NewEncoder(w).Encode(comment)
}
