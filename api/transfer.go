package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/shuifa/go-bank/db/sqlc"
	"github.com/shuifa/go-bank/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrResponse(err))
		return
	}

	formAccount, valid := server.validCurrency(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	payload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	if payload.Username != formAccount.Owner {
		ctx.JSON(http.StatusUnauthorized, ErrResponse(errors.New("from account does not belong to authenticated user")))
		return
	}

	if _, exist := server.validCurrency(ctx, req.ToAccountID, req.Currency); !exist {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
	return
}

func (server *Server) validCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, ErrResponse(err))
				return db.Account{}, false
			}
			ctx.JSON(http.StatusInternalServerError, ErrResponse(err))
			return db.Account{}, false
		}
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency miss match %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, ErrResponse(err))
		return db.Account{}, false
	}

	return account, true
}
