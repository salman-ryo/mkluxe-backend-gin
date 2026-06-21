package utils

func NormalizePagination(page, limit int, sort string) (int, int, string) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	if sort == "" {
		sort = "-created_at"
	}
	return page, limit, sort
}
