package middleware

import (
    "net/http"
    "runtime/debug"
    
    "github.com/gin-gonic/gin"
    "github.com/taluos/Malt/pkg/log"
)

func RecoveryMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        log.Errorf("Panic recovered: %v\nStack: %s", recovered, debug.Stack())
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": 500,
            "msg":  "Internal server error",
        })
    })
}