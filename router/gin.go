package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/thetkpark/lineman-wongnai-intern/covid"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type GinContext struct {
	*gin.Context
}

func (c *GinContext) JSON(status int, v interface{}) {
	c.Context.JSON(status, v)
}

func (c *GinContext) Error(status int, err error) {
	c.Context.JSON(status, gin.H{
		"error": err.Error(),
	})
}

type GinRouter struct {
	logger *zap.SugaredLogger
	port   string
	*gin.Engine
}

type HandlerFunc func(ctx covid.Context)

func NewGinRouter(port string, logger *zap.SugaredLogger) *GinRouter {
	return &GinRouter{Engine: gin.Default(), port: port, logger: logger}
}

func (r *GinRouter) Get(path string, handler HandlerFunc) {
	r.Engine.GET(path, func(c *gin.Context) {
		handler(&GinContext{c})
	})
}

func (r *GinRouter) ListenAndServe() func() {
	s := &http.Server{
		Addr:    ":" + r.port,
		Handler: r,
	}

	go func() {
		r.logger.Infof("Listening on :%s", r.port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.logger.Fatalf("listen: %s\n", err)
		}
	}()

	return func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()
		stop()
		r.logger.Info("shutting down gracefully")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(timeoutCtx); err != nil {
			r.logger.Infow("unable to shutdown", "error", err)
		}
	}
}
