package server

import (
    "net/http"
    "strconv"

    "cryptofyanalytics/internal/models"
    "cryptofyanalytics/internal/storage"

    "github.com/gin-gonic/gin"
)

type Server struct {
    storage storage.Repository
    router  *gin.Engine
}

func NewServer(storage storage.Repository) *Server {
    s := &Server{
        storage: storage,
        router:  gin.Default(),
    }

    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    s.router.GET("/health", s.healthCheck)
    s.router.POST("/prices", s.insertPrice)
    s.router.GET("/prices/:coin_id", s.getPrices)
}

func (s *Server) healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) insertPrice(c *gin.Context) {
    var price models.CryptoPrice
    if err := c.ShouldBindJSON(&price); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := s.storage.InsertPrice(c.Request.Context(), &price); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (s *Server) getPrices(c *gin.Context) {
    coinID := c.Param("coin_id")
    limitStr := c.DefaultQuery("limit", "10")
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
        return
    }

    prices, err := s.storage.GetPrices(c.Request.Context(), coinID, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, prices)
}

func (s *Server) Run(addr string) error {
    return s.router.Run(addr)
}