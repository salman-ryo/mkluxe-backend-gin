package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaginationMeta struct {
	TotalCount  int64 `json:"total_count"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

func Paginated(c *gin.Context, message string, data interface{}, total int64, page, limit int) {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	meta := PaginationMeta{
		TotalCount:  total,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
