package proxy

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulvinamazow/caching-gateway/internal/model"
)

type Cache interface {
	Get(key string) (model.CachedItem, bool, error)
	Set(key string, value model.CachedItem) error
	Clear() error
}

type Handler struct {
	cache  Cache
	origin string
	client *http.Client
}

func NewHandler(cache Cache, origin string) *Handler {
	return &Handler{
		cache:  cache,
		origin: origin,
		client: &http.Client{},
	}
}

func buildCacheKey(method, uri string, body []byte) string {
	hash := sha256.Sum256(body)
	bodyHash := hex.EncodeToString(hash[:])

	return method + ":" + uri + ":" + bodyHash
}

func (handler *Handler) Handle(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	key := buildCacheKey(c.Request.Method, c.Request.URL.RequestURI(), bodyBytes)

	item, found, err := handler.cache.Get(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var skipHeaders = map[string]bool{
		"Access-Control-Allow-Origin":      true,
		"Access-Control-Allow-Methods":     true,
		"Access-Control-Allow-Headers":     true,
		"Access-Control-Allow-Credentials": true,
		"Access-Control-Expose-Headers":    true,
		"Access-Control-Max-Age":           true,
	}

	// CACHE HIT
	if found {
		c.Header("X-Cache", "HIT")

		for k, v := range item.Headers {
			if skipHeaders[k] {
				continue
			}

			for _, vv := range v {
				c.Header(k, vv)
			}
		}

		c.Data(item.StatusCode, item.Headers.Get("Content-Type"), item.Body)
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	originURL := handler.origin + c.Request.URL.RequestURI()

	request, err := http.NewRequest(
		c.Request.Method,
		originURL,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	request = request.WithContext(ctx)

	// headers copy
	for k, v := range c.Request.Header {
		if skipHeaders[k] {
			continue
		}

		for _, vv := range v {
			request.Header.Add(k, vv)
		}
	}

	// HTTP client
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// cache save
	cacheItem := model.CachedItem{
		StatusCode: response.StatusCode,
		Headers:    response.Header,
		Body:       respBody,
	}

	_ = handler.cache.Set(key, cacheItem)

	// MISS header
	c.Header("X-Cache", "MISS")

	// response headers forward
	for k, v := range response.Header {
		for _, vv := range v {
			c.Header(k, vv)
		}
	}

	// response return
	c.Data(response.StatusCode, response.Header.Get("Content-Type"), respBody)
}
