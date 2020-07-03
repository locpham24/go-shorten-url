package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/locpham24/go-shorten-url/db"
	"net/http"
	"net/url"
	"github.com/teris-io/shortid"
	"time"
)
type UrlRequest struct {
	URL string `json:"url"`
}
var redisClient *redis.Client
func InitRouter(redis *redis.Client) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	redisClient = redis
	r.POST("/locpham/generate", getShortenUrl)
	r.GET("/locpham/:key", redirectActualUrl)
	return r
}

func redirectActualUrl(c *gin.Context) {
	uniqueKey := c.Param("key")
	actualUrl, err := redisClient.Get(uniqueKey).Result()
	if err != nil || actualUrl == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "can not get actual url"})
		return
	}
	c.Redirect(http.StatusPermanentRedirect, actualUrl)
}

func getShortenUrl(c *gin.Context) {
	req := UrlRequest{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.URL == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL can not empty"})
		return
	}

	url, err := url.Parse(req.URL)
	if err != nil || url.Host == "" || url.Scheme == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can not parse url"})
		return
	}

	uniqueId, err := shortid.Generate()
	redisDb := db.RedisDb{}
	redisDb.Connect()
	err = redisDb.Client.Set(uniqueId, req.URL, 30*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can not store to redis"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"your_url": req.URL,
		"shorten_url": c.Request.Host + "/locpham/"+uniqueId,
	})
}