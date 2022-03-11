package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shuifa/go-bank/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorization) == 0 {
			err := errors.New("authorization not provider")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse(err))
			return
		}

		fields := strings.Fields(authorization)
		if len(fields) != 2 {
			err := errors.New("invalid authorization fields len")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse(err))
			return
		}

		if fields[0] != AuthorizationTypeBearer {
			err := errors.New("invalid authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse(err))
			return
		}

		payload, err := maker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
