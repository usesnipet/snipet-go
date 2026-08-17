package model

import "time"

type Tenant struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name string  `gorm:"type:varchar(255);not null" json:"name"`
	Slug string  `gorm:"type:varchar(255);not null;unique" json:"slug"`
	Icon *string `gorm:"type:varchar(255)" json:"icon"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	Members     []Member           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	Invitations []TenantInvitation `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	Agents      []Agent            `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	APIKeys     []APIKey           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	Clients     []Client           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	Knowledges  []Knowledge        `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
	LLMs        []LLM              `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"-"`
}
