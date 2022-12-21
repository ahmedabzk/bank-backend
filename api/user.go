package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/ahmedabzk/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"hashed_password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	Username  string    `json:"username"`
	Fullname  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		Fullname:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
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

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponds(err))
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponds(err))
		return
	}

	accessToken,_, err := server.tokenMaker.CreateToken(user.Username, server.config.ACCESS_TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
