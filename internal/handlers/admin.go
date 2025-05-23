package handlers

import (
	"errors"
	"fmt"
	db "lingo/internal/db/sqlc"
	"lingo/pkg/auth/tokengen"
	"lingo/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AdminHandler struct {
	store *db.SQLStore
	tok tokengen.Maker
}

func NewAdminHandler(store *db.SQLStore, tok tokengen.Maker) *AdminHandler {
	return &AdminHandler{
		store: store,
		tok: tok,
	}
}



func (h *AdminHandler) RegisterAdmin(c *gin.Context) {
	var req db.CreateAdminParams
	// Removed unused variable newAdmin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Transaction to create an admin
	err := h.store.ExecTx(c, func(q db.Querier) error {
		// Check if the admin already exists
		_, err := q.GetAdminByEmail(c, req.Email)
		if err != nil {
			if err != pgx.ErrNoRows {
				return err
			}
			// Hash password
			hashedPassword, err := utils.HashPassword(req.Password)
			if err != nil {
				return err
			}
			req.Password = hashedPassword
			// Create the admin
			admin, err := h.store.CreateAdmin(c, req)
			fmt.Print(admin)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin created successfully",
	})
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AdminHandler) LoginAdmin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := h.store.GetAdminForLogin(c, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process login"})
		return
	}

	if !utils.CompareHashAndPassword(admin.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := h.tok.CreateToken(admin.AdminID, 30*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	refreshToken, err := h.tok.CreateToken(admin.AdminID, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "login successful",
		"token":         token,
		"refresh_token": refreshToken,
	})
}
