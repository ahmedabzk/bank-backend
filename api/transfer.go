package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency){
		return 
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency){
		return 
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponds(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponds(err))
		return false
	}

	if account.Currency != currency{
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, ErrorResponds(err))
		return false
	}
	return true
}
