package routes
import (
  "net/http"
  "github.com/gin-gonic/gin"
)
func Start(addr string){ r:=gin.Default(); r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"ok","service":"ground-turn"}) })
  r.GET("/api/flight-turnaround", func(c *gin.Context){ c.JSON(http.StatusOK, []gin.H{{"id":1,"name":"航班过站","status":"READY"}}) })
  r.GET("/api/ground-task", func(c *gin.Context){ c.JSON(http.StatusOK, []gin.H{{"id":1,"name":"地勤任务","status":"READY"}}) })
  r.GET("/api/ground-resource", func(c *gin.Context){ c.JSON(http.StatusOK, []gin.H{{"id":1,"name":"保障资源","status":"READY"}}) })
  r.GET("/api/resource-booking", func(c *gin.Context){ c.JSON(http.StatusOK, []gin.H{{"id":1,"name":"资源预约","status":"READY"}}) })
  r.GET("/api/delay-event", func(c *gin.Context){ c.JSON(http.StatusOK, []gin.H{{"id":1,"name":"延误事件","status":"READY"}}) })
  r.Run(addr)
}
