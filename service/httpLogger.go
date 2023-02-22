package service

// HTTP request logging via zap package.

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HTTPLogCfg is config setting for Ginzap
type HTTPLogCfg struct {
	TimeFormat string
	UTC        bool
	SkipPaths  []string
}

// Logger returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Parameters:
//  1. A logger object (zap)
//  2. A time package format string (e.g. time.RFC3339).
//  3. A boolean stating whether to use UTC time zone or local.
func Logger(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return LoggerWithConfig(logger, &HTTPLogCfg{TimeFormat: timeFormat, UTC: utc})
}

// LoggerWithConfig returns a gin.HandlerFunc using configs
func LoggerWithConfig(logger *zap.Logger, conf *HTTPLogCfg) gin.HandlerFunc {
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}

			if len(c.Errors) > 0 {
				for _, e := range c.Errors.Errors() {
					logger.Error(e)
				}
			} else {
				fields := []zapcore.Field{
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.Duration("latency", latency),
				}
				if conf.TimeFormat != "" {
					fields = append(fields, zap.String("time", end.Format(conf.TimeFormat)))
				}

				switch {
				case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
					{
						logger.Warn(path, fields...)
					}
				case c.Writer.Status() >= http.StatusInternalServerError:
					{
						logger.Error(path, fields...)
					}
				default:
					logger.Info(path, fields...)
				}
			}
		}
	}
}
