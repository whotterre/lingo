package handlers

import (
	"errors"
	"fmt"
	db "lingo/internal/db/sqlc"
	"lingo/pkg/auth/tokengen"
	"lingo/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AdminHandler struct {
	store *db.SQLStore
	tok   tokengen.Maker
}

func NewAdminHandler(store *db.SQLStore, tok tokengen.Maker) *AdminHandler {
	return &AdminHandler{
		store: store,
		tok:   tok,
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

// CreateNewLanguage creates a new language in the database
func (h *AdminHandler) CreateNewLanguage(c *gin.Context) {
	var req db.CreateLanguageParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the language already exists
	_, err := h.store.GetLanguageByName(c, req.LanguageName)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "language already exists"})
		return
	}
	if err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check language"})
		return
	}

	language, err := h.store.CreateLanguage(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Language created successfully",
		"language": language,
	})
}

// CreateNewCourse creates a new course in the database
// It requires the language ID as a URL parameter.
func (h *AdminHandler) CreateNewCourse(c *gin.Context) {
	var req db.CreateCourseParams
	langId := c.Param("langId")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert lang id from string to uuid
	langIdUUID, err := utils.StringToPgTypeUUID(langId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language ID format"})
		return
	}

	// Check if the language exists
	_, err = h.store.GetLanguageById(c, langIdUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check language"})
		return
	}

	// Assign the language ID to the request
	req.LanguageID = langIdUUID

	course, err := h.store.CreateCourse(c, req)
	if err != nil {
		log.Print("Error creating course:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course created successfully",
		"course":  course,
	})
}
