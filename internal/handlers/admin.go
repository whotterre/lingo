package handlers

import (
	"errors"
	"fmt"
	db "lingo/internal/db/sqlc"
	"lingo/pkg/auth/tokengen"
	"lingo/utils"
	"log"
	"net/http"
	"strconv"
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

func (h *AdminHandler) CreateNewLesson(c *gin.Context) {
	courseId := c.Param("courseId")
	var req db.CreateLessonParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Convert to PostgreSQL UUID data type from string
	courseUUID, err := utils.StringToPgTypeUUID(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("couldn't convert course ID to Postgres UUID data type %s", err),
		})
		return
	}
	req.CourseID = courseUUID
	// Create new user record for new exercise
	newLesson, err := h.store.CreateLesson(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("Couldn't create new lesson because %s", err),
		})
		return
	}

	// On successful response
	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson created successfully",
		"lesson":  newLesson,
	})

}

// CreateNewExercise creates a new exercise for a lesson in the database
// It requires the lesson ID as a URL parameter.
func (h *AdminHandler) CreateNewExercise(c *gin.Context) {
	lessonId := c.Param("lessonId")
	var req db.CreateExerciseParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Convert to PostgreSQL UUID data type from string
	lessonUUID, err := utils.StringToPgTypeUUID(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("couldn't convert course ID to Postgres UUID data type %s", err),
		})
		return
	}
	req.LessonID = lessonUUID
	// Create new user record for new exercise
	newExercise, err := h.store.CreateExercise(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("Couldn't create new exercise because %s", err),
		})
		return
	}

	// Successful
	c.JSON(http.StatusOK, gin.H{
		"message":    "Exercise created successfully",
		"exerciseId": newExercise.ExerciseID,
		"exercise":   newExercise,
	})

}

/* Retrieval routes */
func (h *AdminHandler) GetAvailableLanguages(c *gin.Context) {
	var req db.GetAllLanguagesParams
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	limit := int32(limitInt)
	req.Limit = limit
	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	req.Limit = limit
	req.Offset = int32(offsetInt)
	languages, err := h.store.GetAllLanguages(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("Couldn't create get languages because %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Available languages",
		"languages": languages,
	})
}

func (h *AdminHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.store.GetAllCourses(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("Couldn't create get courses because %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "Available courses",
		"courses": courses,
	})
}

func (h *AdminHandler) GetAllLessons(c *gin.Context) {
	var req db.GetAllLessonsParams
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	req.Limit = int32(limitInt)
	req.Offset = int32(offsetInt)

	lessons, err := h.store.GetAllLessons(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("Couldn't create get lessons because %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Available lessons",
		"lessons": lessons,
	})
}

func (h *AdminHandler) GetAllExercises(c *gin.Context) {
	var req db.GetAllExercisesParams
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	req.Limit = int32(limitInt)
	req.Offset = int32(offsetInt)

	exercises, err := h.store.GetAllExercises(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("Couldn't create get exercises because %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "Available exercises",
		"exercises": exercises,
	})
}

func (h *AdminHandler) GetLessonsByCourseId(c *gin.Context) {
	var req db.GetLessonsByCourseIdParams
	courseId := c.Param("courseId")

	if courseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// Convert courseId from string to UUID
	courseUUID, err := utils.StringToPgTypeUUID(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid course ID format: %s", err)})
		return
	}

	req.CourseID = courseUUID

	lessons, err := h.store.GetLessonsByCourseId(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lessons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lessons retrieved successfully",
		"lessons": lessons,
	})
}

func (h *AdminHandler) GetExercisesByLessonId(c *gin.Context) {
	var req db.GetExercisesByLessonIdParams
	lessonId := c.Param("lessonId")

	if lessonId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lesson ID is required"})
		return
	}

	// Convert lessonId from string to UUID
	lessonUUID, err := utils.StringToPgTypeUUID(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid lesson ID format: %s", err)})
		return
	}

	req.LessonID = lessonUUID

	exercises, err := h.store.GetExercisesByLessonId(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exercises"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Exercises retrieved successfully",
		"exercises": exercises,
	})
}

func (h *AdminHandler) GetExerciseById(c *gin.Context) {
	exerciseId := c.Param("exerciseId")
	if exerciseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exercise ID is required"})
		return
	}

	// Convert exerciseId from string to UUID
	exerciseUUID, err := utils.StringToPgTypeUUID(exerciseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid exercise ID format: %s", err)})
		return
	}

	exercise, err := h.store.GetExerciseById(c, exerciseUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Exercise retrieved successfully",
		"exercise": exercise,
	})
}

// UpdateAdminDetails updates the details of an admin
func (h *AdminHandler) UpdateAdminDetails(c *gin.Context) {
	var req db.UpdateAdminDetailsParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminId := c.Param("adminId")
	if adminId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin ID is required"})
		return
	}

	// Convert adminId from string to UUID
	adminUUID, err := utils.StringToPgTypeUUID(adminId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid admin ID format: %s", err)})
		return
	}

	req.AdminID = adminUUID

	err = h.store.UpdateAdminDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin details updated successfully",
	})
}

// UpdateAdminPassword updates the password of an admin
func (h *AdminHandler) UpdateAdminPassword(c *gin.Context) {
	var req db.UpdateAdminPasswordParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminId := c.Param("adminId")
	if adminId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin ID is required"})
		return
	}

	// Convert adminId from string to UUID
	adminUUID, err := utils.StringToPgTypeUUID(adminId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid admin ID format: %s", err)})
		return
	}

	req.AdminID = adminUUID

	err = h.store.UpdateAdminPassword(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin password updated successfully",
	})
}

func (h *AdminHandler) UpdateLanguageById(c *gin.Context) {
	var req db.UpdateLanguageDetailsParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	languageId := c.Param("languageId")
	if languageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language ID is required"})
		return
	}

	// Convert languageId from string to UUID
	languageUUID, err := utils.StringToPgTypeUUID(languageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid language ID format: %s", err)})
		return
	}

	req.LanguageID = languageUUID

	err = h.store.UpdateLanguageDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Language updated successfully",
	})
}

func (h *AdminHandler) UpdateCourseById(c *gin.Context) {
	var req db.UpdateCourseDetailsParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseId := c.Param("courseId")
	if courseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// Convert courseId from string to UUID
	courseUUID, err := utils.StringToPgTypeUUID(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid course ID format: %s", err)})
		return
	}

	req.CourseID = courseUUID

	err = h.store.UpdateCourseDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course updated successfully",
	})
}

func (h *AdminHandler) UpdateLessonById(c *gin.Context) {
	var req db.UpdateLessonDetailsParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lessonId := c.Param("lessonId")
	if lessonId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lesson ID is required"})
		return
	}

	// Convert lessonId from string to UUID
	lessonUUID, err := utils.StringToPgTypeUUID(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid lesson ID format: %s", err)})
		return
	}

	req.LessonID = lessonUUID

	err = h.store.UpdateLessonDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lesson"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson updated successfully",
	})
}

func (h *AdminHandler) UpdateExerciseById(c *gin.Context) {
	var req db.UpdateExerciseDetailsParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exerciseId := c.Param("exerciseId")
	if exerciseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exercise ID is required"})
		return
	}

	// Convert exerciseId from string to UUID
	exerciseUUID, err := utils.StringToPgTypeUUID(exerciseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid exercise ID format: %s", err)})
		return
	}

	req.ExerciseID = exerciseUUID

	err = h.store.UpdateExerciseDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercise updated successfully",
	})
}

func (h *AdminHandler) DeleteLanguage(c *gin.Context) {
	languageId := c.Param("languageId")
	if languageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language ID is required"})
		return
	}

	// Convert languageId from string to UUID
	languageUUID, err := utils.StringToPgTypeUUID(languageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid language ID format: %s", err)})
		return
	}

	err = h.store.DeleteLanguage(c, languageUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Language deleted successfully",
	})
}

func (h *AdminHandler) DeleteCourse(c *gin.Context) {
	courseId := c.Param("courseId")
	if courseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	// Convert courseId from string to UUID
	courseUUID, err := utils.StringToPgTypeUUID(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid course ID format: %s", err)})
		return
	}

	err = h.store.DeleteCourse(c, courseUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course deleted successfully",
	})
}

func (h *AdminHandler) DeleteLesson(c *gin.Context) {
	lessonId := c.Param("lessonId")
	if lessonId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lesson ID is required"})
		return
	}

	// Convert lessonId from string to UUID
	lessonUUID, err := utils.StringToPgTypeUUID(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid lesson ID format: %s", err)})
		return
	}

	err = h.store.DeleteLesson(c, lessonUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lesson"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson deleted successfully",
	})
}

func (h *AdminHandler) DeleteExercise(c *gin.Context) {
	exerciseId := c.Param("exerciseId")
	if exerciseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exercise ID is required"})
		return
	}

	// Convert exerciseId from string to UUID
	exerciseUUID, err := utils.StringToPgTypeUUID(exerciseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid exercise ID format: %s", err)})
		return
	}

	err = h.store.DeleteExercise(c, exerciseUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercise deleted successfully",
	})
}
