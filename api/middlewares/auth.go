package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/koizumi7010/blog-api/apperrors"
	"github.com/koizumi7010/blog-api/common"
	"google.golang.org/api/idtoken"
)

const (
	googleClientID = "974808245159-5t60iq26fa8imcgbc7mlrs7tup30csel.apps.googleusercontent.com"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// HeaderからAuthorizationを取得
		authorization := req.Header.Get("Authorization")

		// Authorizationフィールドを検証
		authHeader := strings.Split(authorization, " ")
		if len(authHeader) != 2 {
			err := apperrors.RequiredAuthorizationHeader.Wrap(errors.New("invalid request header"), "invalid header")
			apperrors.ErrorHandler(w, req, err)
			return
		}

		bearer, idToken := authHeader[0], authHeader[1]
		if bearer != "Bearer" || idToken == "" {
			err := apperrors.RequiredAuthorizationHeader.Wrap(errors.New("invalid request header"), "invalid header")
			apperrors.ErrorHandler(w, req, err)
			return
		}

		// IDトークンを検証
		tokenValidator, err := idtoken.NewValidator(context.Background())
		if err != nil {
			err := apperrors.CannotMakeValidator.Wrap(err, "internal auth error")
			apperrors.ErrorHandler(w, req, err)
			return
		}

		payload, err := tokenValidator.Validate(context.Background(), idToken, googleClientID)
		if err != nil {
			err := apperrors.Unauthorizated.Wrap(err, "invalid token")
			apperrors.ErrorHandler(w, req, err)
			return
		}

		name, ok := payload.Claims["name"]
		if !ok {
			err := apperrors.Unauthorizated.Wrap(err, "invalid token")
			apperrors.ErrorHandler(w, req, err)
			return
		}

		req = common.SetUserName(req, name.(string))

		next.ServeHTTP(w, req)
	})
}
