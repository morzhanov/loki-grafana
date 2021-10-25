package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/morzhanov/go-elk-example/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/go-elk-example/internal/doc"
	"github.com/morzhanov/go-elk-example/internal/es"
	"go.uber.org/zap"
)

type rest struct {
	esearch es.ElasticSearch
	log     *zap.Logger
	mc      metrics.Collector
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type REST interface {
	Listen()
}

func (r *rest) handleHttpErr(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
	r.log.Error("error in the handler", zap.Error(err))
}

func (r *rest) handleFind(c *gin.Context) {
	fName := c.Param("field")
	val := c.Param("value")
	res, err := r.esearch.Find(fName, val)
	if err != nil {
		r.handleHttpErr(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (r *rest) handleUpdate(c *gin.Context) {
	id := c.Param("id")
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		r.handleHttpErr(c, err)
		return
	}
	var d doc.Document
	if err := json.Unmarshal(jsonData, &d); err != nil {
		r.handleHttpErr(c, err)
		return
	}

	if err := r.esearch.Update(id, &d); err != nil {
		r.handleHttpErr(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (r *rest) handleDelete(c *gin.Context) {
	id := c.Param("id")
	if err := r.esearch.Delete(id); err != nil {
		r.handleHttpErr(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (r *rest) metricsMiddleware(c *gin.Context) {
	c.Writer = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Next()
	statusCode := c.Writer.Status()
	r.log.Info("handled REST request", zap.Int("status_code", statusCode))
	r.mc.IncReq()
	if statusCode >= 300 {
		r.mc.IncErrs()
	}
}

func (r *rest) Listen() {
	router := gin.Default()
	router.GET("/debug/vars", r.mc.GetHandler())
	router.GET("/:field/:value", r.handleFind)
	router.PUT("/:id", r.handleUpdate)
	router.DELETE("/:id", r.handleDelete)
	router.Use(r.metricsMiddleware)
	if err := router.Run(); err != nil {
		r.log.Fatal("error during rest controller execution", zap.Error(err))
	}
}

func NewREST(esearch es.ElasticSearch, l *zap.Logger, mc metrics.Collector) REST {
	return &rest{esearch, l, mc}
}
