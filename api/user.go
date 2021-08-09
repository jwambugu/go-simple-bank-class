package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/jwambugu/go-simple-bank-class/db/sqlc"
	"github.com/jwambugu/go-simple-bank-class/util"
	"github.com/lib/pq"
	"net/http"
	"time"
)

type (
	createUserRequest struct {
		Username string `json:"username" binding:"required,alphanum"`
		FullName string `json:"full_name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}

	userResponse struct {
		Username          string    `json:"username"`
		FullName          string    `json:"fullName"`
		Email             string    `json:"email"`
		PasswordChangedAt time.Time `json:"passwordChangedAt"`
		CreatedAt         time.Time `json:"createdAt"`
	}

	loginUserRequest struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
	}

	loginUserResponse struct {
		AccessToken string       `json:"access_token"`
		User        userResponse `json:"user"`
	}
)

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	// Create a new user
	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				{
					ctx.JSON(http.StatusForbidden, errorResponse(err))
					return
				}
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Find the user by username
	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Compare the passwords
	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Attempt to generate an access token
	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, response)
}
