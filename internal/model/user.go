package model

import "time"

type Challenge string

const (
	ChallengeActiveAccount  Challenge = "active_account"
	ChallengeChangePassword Challenge = "change_password"
)

type User struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name         string      `gorm:"type:varchar(255);not null" json:"name"`
	Email        string      `gorm:"type:varchar(255);not null;unique" json:"email"`
	PasswordHash *string     `gorm:"type:varchar(255)" json:"-"`
	Picture      *string     `gorm:"type:varchar(255)" json:"picture"`
	IsAdmin      bool        `gorm:"type:boolean;not null;default:false" json:"-"`
	Challenges   []Challenge `gorm:"type:jsonb;not null;serializer:json" json:"-"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	Accounts []Account `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Tokens   []Token   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Members  []Member  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
