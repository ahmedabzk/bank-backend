package api

import (
	"net/http"

	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID   int64 `json:"to_account_id" binding:"required"`
	Amount        int64 `json:"amount" binding:"required"`
}

func(server *Server) createTransfer(ctx *gin.Context){
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return 
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
		
	}

	transfer, err := server.store.CreateTransfer(ctx, arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return 
	}

	ctx.JSON(http.StatusOK, transfer)
}
