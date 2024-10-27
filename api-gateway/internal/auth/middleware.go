package auth

import (
	"context"
	"errors"
	pb "github.com/cloud9cloud9/go-grpc-todo/api-gateway/internal/auth/pb"
	"github.com/cloud9cloud9/go-grpc-todo/api-gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authHeader = "Authorization"
	bearer     = "Bearer "
	Key        = "userId"
)

var (
	unauthorizedBody  = "unauthorized"
	errUserIdNotFound = errors.New("user id not found")
	errInvalidUserID  = errors.New("user id is of invalid type")
)

type Middleware struct {
	svc *ServiceClient
}

func InitMiddleware(svc *ServiceClient) *Middleware {
	return &Middleware{
		svc: svc,
	}
}

func (m *Middleware) UserIdentity(ctx *gin.Context) {
	header := ctx.Request.Header.Get(authHeader)

	if header == "" {
		utils.NewErrorResponse(ctx, http.StatusUnauthorized, unauthorizedBody)
		return
	}

	token := strings.Split(header, bearer)

	if len(token) != 2 {
		utils.NewErrorResponse(ctx, http.StatusUnauthorized, unauthorizedBody)
		return
	}

	res, err := m.svc.Client.Validate(context.Background(), &pb.ValidateRequest{
		Token: token[1],
	})

	if err != nil || res.Status != http.StatusOK {
		utils.NewErrorResponse(ctx, http.StatusUnauthorized, unauthorizedBody)
		return
	}

	ctx.Set(Key, res.UserId)

	ctx.Next()
}

func GetUserId(c *gin.Context) (int64, error) {
	id, ok := c.Get(Key)
	if !ok {
		return 0, errUserIdNotFound
	}

	idInt, ok := id.(int64)
	if !ok {
		return 0, errInvalidUserID
	}

	return idInt, nil
}
