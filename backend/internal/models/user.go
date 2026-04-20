package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int            `db:"id"`
	Email      string         `db:"email"`
	Username   *string        `db:"username"`
	Password   sql.NullString `db:"password"`
	GoogleSub  sql.NullString `db:"google_sub"`
	ImageURL   *string        `db:"image_url"`
	Role       string     `db:"role"`
	IsVerified bool       `db:"is_verified"`
	VerifiedAt *time.Time `db:"verified_at"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}
