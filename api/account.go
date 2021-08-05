package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/jwambugu/go-simple-bank-class/db/sqlc"
	"net/http"
)

type (
	createAccountRequest struct {
		Owner    string `json:"owner" binding:"required"`
		Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
	}

	getAccountByIDRequest struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	getAccountsRequest struct {
		PageID   int32 `form:"page_id" binding:"required,min=1"`
		PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
	}
)

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	// Create account
	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAccountByID(ctx *gin.Context) {
	var req getAccountByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get account
	account, err := server.store.GetAccount(ctx, int32(req.ID))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var req getAccountsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	// Get accounts
	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
