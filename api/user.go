package api

import (
	"net/http"
	"time"

	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/ahmedabzk/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	HashedPassword string `json:"hashed_password" binding:"required,min=6"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

type userReturnParams struct {
	Username string
	Fullname string
	Email string
	CreatedAt time.Time

}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponds(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if ErrPq, ok := err.(*pq.Error); ok {
			switch ErrPq.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponds(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	rsp := userReturnParams{
		Username: user.Username,
		Fullname: user.FullName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
