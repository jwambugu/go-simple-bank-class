package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/jwambugu/go-simple-bank-class/db/sqlc"
	"net/http"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) isValidAccount(ctx *gin.Context, accountID int64, currency string) bool {
	// Find the account using the account id
	account, err := server.store.GetAccount(ctx, int32(accountID))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return false
	}

	return true
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if the sender account is valid
	if !server.isValidAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	// Check if the receiver account is valid
	if !server.isValidAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// Create a new transfer
	transferTxResult, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transferTxResult)
}
