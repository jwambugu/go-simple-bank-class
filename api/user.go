package api

import (
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

	createUserResponse struct {
		Username          string    `json:"username"`
		FullName          string    `json:"fullName"`
		Email             string    `json:"email"`
		PasswordChangedAt time.Time `json:"passwordChangedAt"`
		CreatedAt         time.Time `json:"createdAt"`
	}
)

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

	userResponse := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, userResponse)
}
