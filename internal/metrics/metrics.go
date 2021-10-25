package metrics

import (
	"expvar"
	_ "expvar"
	"fmt"

	"github.com/gin-gonic/gin"
)

var reqCount *expvar.Int
var docsGenerated *expvar.Int
var errCount *expvar.Int

type collector struct {
	handler func(c *gin.Context)
}

type Collector interface {
	GetHandler() func(c *gin.Context)
	IncReq()
	IncDocs()
	IncErrs()
}

func (c *collector) GetHandler() func(c *gin.Context) {
	return c.handler
}

func (c *collector) IncReq() {
	reqCount.Add(1)
}

func (c *collector) IncDocs() {
	docsGenerated.Add(1)
}

func (c *collector) IncErrs() {
	errCount.Add(1)
}

func metricsHandler(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "application/json; charset = utf-8")

	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(w, "% q:% q", key, str)
		} else {
			fmt.Fprintf(w, "% q:% v", key, value)
		}
	}

	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func NewMetricsCollector() Collector {
	reqCount = expvar.NewInt("req_count")
	docsGenerated = expvar.NewInt("docs_generated")
	errCount = expvar.NewInt("err_count")
	return &collector{metricsHandler}
}
