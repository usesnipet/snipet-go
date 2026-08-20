package model

import "time"

type Tenant struct {
	ID string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name string  `gorm:"type:varchar(255);not null" json:"name"`
	Slug string  `gorm:"type:varchar(255);not null;unique" json:"slug"`
	Icon *string `gorm:"type:varchar(255)" json:"icon"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	Members     []Member           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"members"`
	Invitations []TenantInvitation `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"invitations"`
	Agents      []Agent            `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"agents"`
	APIKeys     []APIKey           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"api_keys"`
	Clients     []App              `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"clients"`
	Knowledges  []Knowledge        `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"knowledges"`
	LLMs        []LLM              `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"llms"`
}
