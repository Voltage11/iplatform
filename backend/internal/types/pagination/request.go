package pagination

type PaginationRequest struct {
	page  int
	limit int
}

func (p PaginationRequest) GetPage() int  { return p.page }
func (p PaginationRequest) GetLimit() int { return p.limit }

// GetOffset рассчитывает сдвиг для SQL-запроса (OFFSET)
// Для первой страницы возвращает 0.
func (p PaginationRequest) GetOffset() int {
	if p.page < 1 {
		return 0
	}
	return (p.page - 1) * p.GetLimit()
}

// NewPaginationRequest создаёт запрос пагинации с валидными значениями
func NewPaginationRequest(page, limit int) PaginationRequest {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	return PaginationRequest{
		page:  page,
		limit: limit,
	}
}
