package vo

import (
	"strings"

	errs "github.com/hydr0g3nz/clean_stater/internal/domain/error"
)

// UserRole represents the user role
type UserRole string

const (
	RoleCandidate UserRole = "candidate"
	RoleCompanyHR UserRole = "company_hr"
	RoleAdmin     UserRole = "admin"
)

func (role UserRole) String() string {
	switch role {
	case RoleCandidate:
		return "candidate"
	case RoleCompanyHR:
		return "company_hr"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// ParseUserRole parses a user role string and returns the corresponding UserRole
func ParseUserRole(role string) (UserRole, error) {
	switch strings.ToLower(role) {
	case string(RoleCandidate):
		return RoleCandidate, nil
	case string(RoleCompanyHR):
		return RoleCompanyHR, nil
	case string(RoleAdmin):
		return RoleAdmin, nil
	default:
		return "", errs.ErrInvalidUserRole
	}
}

// IsValid returns true if the user role is valid, false otherwise
func (r UserRole) IsValid() bool {
	switch r {
	case RoleCandidate, RoleCompanyHR, RoleAdmin:
		return true
	default:
		return false
	}
}
