// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAdmin(ctx context.Context, arg CreateAdminParams) (Admin, error)
	CreateCourse(ctx context.Context, arg CreateCourseParams) (Course, error)
	CreateExercise(ctx context.Context, arg CreateExerciseParams) (Exercise, error)
	// Maintenance Functionality --
	// Creational Queries
	CreateLanguage(ctx context.Context, arg CreateLanguageParams) (Language, error)
	CreateLesson(ctx context.Context, arg CreateLessonParams) (Lesson, error)
	CreateUserCourse(ctx context.Context, arg CreateUserCourseParams) (UserCourse, error)
	CreateUserProgress(ctx context.Context, arg CreateUserProgressParams) (UserProgress, error)
	DeleteAdmin(ctx context.Context, adminID pgtype.UUID) error
	GetAdminByEmail(ctx context.Context, email string) (GetAdminByEmailRow, error)
	GetAdminById(ctx context.Context, adminID pgtype.UUID) (GetAdminByIdRow, error)
	GetAdminForLogin(ctx context.Context, email string) (GetAdminForLoginRow, error)
	GetAllAdmins(ctx context.Context, arg GetAllAdminsParams) ([]GetAllAdminsRow, error)
	// Retrieval Queries
	GetAllLanguages(ctx context.Context, arg GetAllLanguagesParams) ([]Language, error)
	GetLanguageById(ctx context.Context, languageID pgtype.UUID) (Language, error)
	GetLanguageByName(ctx context.Context, languageName string) (Language, error)
	UpdateAdmin(ctx context.Context, arg UpdateAdminParams) error
}

var _ Querier = (*Queries)(nil)
