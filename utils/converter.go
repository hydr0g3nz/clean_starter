package utils

import "github.com/jackc/pgx/v5/pgtype"

func ConvertToBool(pgBool pgtype.Bool) bool {
	if pgBool.Valid {
		return pgBool.Bool
	}
	return false
}
func ConvertToPGBool(value bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  value,
		Valid: true,
	}
}
