package sql

func sqlAndWhere(where string) string {
	if where == "" {
		return ""
	}
	return " AND"
}
