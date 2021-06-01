package baseutils

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type ErrorMonitor struct {
	Name             string //自定义命名 用于直观原因:passport/odin ==自定义维度,用英文标识
	ExceptionMethod  string //异常方法
	ExceptionClass   string //异常类
	ExceptionMessage string //异常简要描述
	ExceptionLevel   string //等级 error/warn/info
}

var (
	totalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "exception",
			Subsystem: "request",
			Name:      "counter",
			Help:      "Total number of Exception/Error processed by the workers",
		},
		// We will want to monitor the worker ID that processed the
		// job, and the type of job that was processed
		[]string{"name", "exception_method", "exception_class", "exception_message", "level"},
	)
	// http_server_request_duration_seconds，Histogram类型指标，bucket代表duration的分布区间
	WebRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_http_server_request_duration_seconds",
			Help:    "go http_server request duration distribution",
			Buckets: []float64{0.1, 0.2, 0.5, 0.7, 0.9, 1},
		},
		[]string{"method", "uri"},
	)
)

func InitMonitor() {
	prometheus.MustRegister(totalCounterVec)
	prometheus.MustRegister(WebRequestDuration)
}

func init() {
	prometheus.MustRegister(totalCounterVec)
	prometheus.MustRegister(WebRequestDuration)
}

//添加报警信息
//自定义命名 用于直观原因:passport/odin ==自定义维度,用英文标识
//异常方法
//异常类
//异常简要描述
//等级 error/warn/info
func ExceptionMonitorAdd(name, exceptionMethod, exceptionClass, exceptionMessage, exceptionLevel string) {
	if len(name) == 0 {
		name = "empty"
	}
	if len(exceptionMethod) == 0 {
		exceptionMethod = "empty"
	}

	if len(exceptionClass) == 0 {
		exceptionClass = "empty"
	}

	if len(exceptionMessage) == 0 {
		exceptionMessage = "empty"
	}

	if len(exceptionLevel) != 0 && exceptionLevel == "error" {
		exceptionLevel = "error"
	} else if len(exceptionLevel) != 0 && exceptionLevel == "warning" {
		exceptionLevel = "warning"
	} else {
		exceptionLevel = "info"
	}
	data := ErrorMonitor{
		Name:             name,
		ExceptionMethod:  exceptionMethod,
		ExceptionLevel:   exceptionLevel,
		ExceptionClass:   exceptionClass,
		ExceptionMessage: exceptionMessage,
	}
	ExceptionMonitor(data)
}

//用户统计报警信息
func ExceptionMonitor(data ErrorMonitor) {
	totalCounterVec.WithLabelValues(data.Name, data.ExceptionMethod, data.ExceptionClass, data.ExceptionMessage, data.ExceptionLevel).Inc()
}

//在每个输出的方法中调用该方法
func RequestFinally(request *http.Request, startTime time.Time) {
	duration := time.Since(startTime)
	WebRequestDuration.With(prometheus.Labels{"method": request.Method, "uri": request.URL.Path}).Observe(duration.Seconds())
}

//RequestFinally 与 RequestMethodFinally 二选一
func RequestMethodFinally(method, uri string, startTime time.Time) {
	duration := time.Since(startTime)
	WebRequestDuration.With(prometheus.Labels{"method": method, "uri": uri}).Observe(duration.Seconds())
}
