package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"net/http"
	"time"
)

func CircuitBreaker() gin.HandlerFunc {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "api-breaker",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})

	return func(c *gin.Context) {
		result, err := cb.Execute(func() (interface{}, error) {
			c.Next()
			return nil, nil
		})

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Service temporarily unavailable",
			})
			c.Abort()
			return
		}

		if result != nil {
			c.JSON(http.StatusOK, result)
		}
	}
} 