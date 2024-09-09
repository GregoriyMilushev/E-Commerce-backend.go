package models

import (
	"errors"

	"gorm.io/gorm"
)

type Role string

const (
    RoleAdmin  Role = "admin"
    RoleClient Role = "client"
)

type User struct {
    gorm.Model
    Name              string    `gorm:"size:100;not null" json:"name"`
    Email             string    `gorm:"size:100;not null;unique" json:"email"`
    Password          string    `gorm:"size:255;not null" json:"-"`
    Role              Role      `gorm:"size:100;not null" json:"role"`
    GoogleID          *string    `gorm:"size:255;unique" json:"google_id,omitempty"`
    ProfilePictureURL *string    `gorm:"size:255" json:"profile_picture_url"`
    BillingAddress    *string    `gorm:"type:text" json:"billing_address"`
    StripeCustomerID  *string    `gorm:"size:255" json:"stripe_customer_id"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    if !u.IsValidRole() {
        return errors.New("invalid role: must be either 'admin' or 'client'")
    }
    return nil
}

func (u *User) IsValidRole() bool {
    return u.Role == RoleAdmin || u.Role == RoleClient
}
