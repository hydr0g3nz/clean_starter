package entity

// User represents the users table
type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
}

func (u User) ToNotEmptyValueMap() map[string]interface{} {
	result := make(map[string]interface{})
	if u.FirstName != "" {
		result["first_name"] = u.FirstName
	}
	if u.LastName != "" {
		result["last_name"] = u.LastName
	}
	if u.Email != "" {
		result["email"] = u.Email
	}
	if u.Phone != "" {
		result["phone"] = u.Phone
	}
	return result
}

type UserFilter struct {
	FirstName *string
	LastName  *string
	Email     *string
	Phone     *string
}
