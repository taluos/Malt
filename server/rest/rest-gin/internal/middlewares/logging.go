package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/taluos/Malt/pkg/log"
)

func LoggingMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        log.Infof("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
        return ""
    })
}