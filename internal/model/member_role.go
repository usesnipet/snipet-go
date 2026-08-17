package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type MemberRole int

const (
	RoleUser MemberRole = iota + 1
	RoleManager
	RoleAdmin
)

var roleNames = map[MemberRole]string{
	RoleUser:    "user",
	RoleManager: "manager",
	RoleAdmin:   "admin",
}
var roleValues = map[string]MemberRole{
	"user": RoleUser, "manager": RoleManager, "admin": RoleAdmin,
}

func (r MemberRole) String() string              { return roleNames[r] }
func (r MemberRole) AtLeast(min MemberRole) bool { return r >= min }

func (r MemberRole) MarshalJSON() ([]byte, error) { return json.Marshal(r.String()) }
func (r *MemberRole) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	role, ok := roleValues[s]
	if !ok {
		return fmt.Errorf("invalid role %q", s)
	}
	*r = role
	return nil
}

func (r MemberRole) Value() (driver.Value, error) { return r.String(), nil }
func (r *MemberRole) Scan(v any) error {
	s, _ := v.(string)
	role, ok := roleValues[s]
	if !ok {
		return fmt.Errorf("invalid role in db %q", s)
	}
	*r = role
	return nil
}
