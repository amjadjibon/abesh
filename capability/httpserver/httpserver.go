package httpserver

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"github.com/mkawserm/abesh/utility"
)

var ErrPathNotDefined = errors.New("path not defined")
var ErrMethodNotDefined = errors.New("method not defined")

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "abesh_httpserver_response_status",
		Help: "Status of HTTP Response",
	},
	[]string{"path", "status"},
)

var panicCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "abesh_httpserver_panic_counter",
		Help: "Abesh HTTP Server Panic Counter",
	},
	[]string{"contractid"},
)

type EventResponse struct {
	Error error
	Event *model.Event
}

type HTTPServer struct {
	mHost     string
	mPort     string
	mCertFile string
	mKeyFile  string

	mStaticDir  string
	mStaticPath string
	mHealthPath string

	mDefault404HandlerEnabled bool
	mValues                   model.ConfigMap
	mHttpServer               *http.Server
	mHttpServerMux            *http.ServeMux
	mEventTransmitter         iface.IEventTransmitter

	mRequestTimeout     time.Duration
	mDefaultContentType string

	mEmbeddedStaticFSMap map[string]embed.FS

	d401m string
	d403m string
	d404m string
	d405m string
	d408m string
	d409m string
	d499m string
	d500m string

	mIsMetricsEnabled bool
	mMetricPath       string
}

func (h *HTTPServer) Name() string {
	return "abesh_httpserver"
}

func (h *HTTPServer) Version() string {
	return constant.Version
}

func (h *HTTPServer) Category() string {
	return string(constant.CategoryTrigger)
}

func (h *HTTPServer) ContractId() string {
	return "abesh:httpserver"
}

func (h *HTTPServer) GetConfigMap() model.ConfigMap {
	return h.mValues
}

func (h *HTTPServer) buildDefaultMessage(code uint32) string {
	return fmt.Sprintf(`
		{
			"code": "SE_%d",
			"lang": "en",
			"message": "%d ERROR",
			"data": {}
		}
	`, code, code)
}

func (h *HTTPServer) SetConfigMap(values model.ConfigMap) error {
	h.mValues = values

	h.mHost = h.mValues.String("host", "0.0.0.0")
	h.mPort = h.mValues.String("port", "8080")

	h.mCertFile = values.String("cert_file", "")
	h.mKeyFile = values.String("key_file", "")

	h.mStaticDir = values.String("static_dir", "")
	h.mStaticPath = values.String("static_path", "/static/")
	h.mHealthPath = values.String("health_path", "")

	h.mRequestTimeout = h.mValues.Duration("default_request_timeout", time.Second)
	h.mDefault404HandlerEnabled = h.mValues.Bool("default_404_handler_enabled", true)
	h.mDefaultContentType = values.String("default_content_type", "application/json")

	h.d401m = h.buildDefaultMessage(401)
	h.d403m = h.buildDefaultMessage(403)
	h.d404m = h.buildDefaultMessage(404)
	h.d405m = h.buildDefaultMessage(405)
	h.d408m = h.buildDefaultMessage(408)
	h.d409m = h.buildDefaultMessage(409)
	h.d499m = h.buildDefaultMessage(499)
	h.d500m = h.buildDefaultMessage(500)

	h.mIsMetricsEnabled = h.mValues.Bool("metrics_enabled", false)
	h.mMetricPath = values.String("metric_path", "/metrics")

	return nil
}

func (h *HTTPServer) getMessage(key, defaultValue, lang string) string {
	data := h.mValues.String(fmt.Sprintf("%s_%s", key, lang), "")

	if len(data) == 0 {
		data = h.mValues.String(key, defaultValue)
	}

	return data
}

func (h *HTTPServer) getLanguage(r *http.Request) string {
	l := r.Header.Get("Accept-Language")
	if len(l) == 0 {
		l = "en"
	}

	return l
}

func (h *HTTPServer) SetEventTransmitter(eventTransmitter iface.IEventTransmitter) error {
	h.mEventTransmitter = eventTransmitter
	return nil
}

func (h *HTTPServer) GetEventTransmitter() iface.IEventTransmitter {
	return h.mEventTransmitter
}

func (h *HTTPServer) AddEmbeddedStaticFS(pattern string, fs embed.FS) {
	// NOTE: must be called after setup otherwise panic will occur
	h.mEmbeddedStaticFSMap[pattern] = fs
}

func (h *HTTPServer) New() iface.ICapability {
	return &HTTPServer{}
}

func (h *HTTPServer) AddHandlerFunc(pattern string, handler http.HandlerFunc) {
	h.mHttpServerMux.HandleFunc(pattern, handler)
}

func (h *HTTPServer) AddHandler(pattern string, handler http.Handler) {
	h.mHttpServerMux.Handle(pattern, handler)
}

func (h *HTTPServer) Setup() error {
	h.mHttpServer = new(http.Server)
	h.mHttpServerMux = new(http.ServeMux)
	h.mEmbeddedStaticFSMap = make(map[string]embed.FS)

	// setup server details
	h.mHttpServer.Handler = h.mHttpServerMux
	h.mHttpServer.Addr = h.mHost + ":" + h.mPort

	if h.mDefault404HandlerEnabled {
		h.mHttpServerMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			h.debugMessage(request)

			timerStart := time.Now()

			defer func() {
				logger.L(h.ContractId()).Debug("request completed")
				elapsed := time.Since(timerStart)
				logger.L(h.ContractId()).Debug("request execution time", zap.Duration("seconds", elapsed))
			}()

			h.s404m(request, writer, nil)
			return
		})
	}

	// register data path
	if len(h.mStaticDir) != 0 {
		fi, e := os.Stat(h.mStaticDir)

		if e != nil {
			logger.L(h.ContractId()).Error(e.Error())
		} else {
			if fi.IsDir() {
				logger.L(h.ContractId()).Debug("data path", zap.String("static_path", h.mStaticPath))
				h.mHttpServerMux.Handle(h.mStaticPath, http.StripPrefix(h.mStaticPath, http.FileServer(http.Dir(h.mStaticDir))))
			} else {
				logger.L(h.ContractId()).Error("provided static_dir in the manifest conf is not directory")
			}
		}
	}

	// register health path
	if len(h.mHealthPath) != 0 {
		h.mHttpServerMux.HandleFunc(h.mHealthPath, func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusOK)
			logger.L(h.ContractId()).Info("HEALTH OK")
		})
	}

	logger.L(h.ContractId()).Info("http server setup complete",
		zap.String("host", h.mHost),
		zap.String("port", h.mPort))

	if h.mIsMetricsEnabled {
		h.AddHandler(h.mMetricPath, promhttp.Handler())
		logger.L(h.ContractId()).Info("metrics enabled", zap.String("metric_path", h.mMetricPath))
	}

	return nil
}

func (h *HTTPServer) Start(_ context.Context) error {
	logger.L(h.ContractId()).Debug("registering embedded data fs")
	for p, d := range h.mEmbeddedStaticFSMap {
		h.mHttpServerMux.Handle(p, http.FileServer(http.FS(d)))
	}

	logger.L(h.ContractId()).Info("http server started at " + h.mHttpServer.Addr)

	if len(h.mCertFile) != 0 && len(h.mKeyFile) != 0 {
		if err := h.mHttpServer.ListenAndServeTLS(h.mCertFile, h.mKeyFile); err != http.ErrServerClosed {
			return err
		}
	} else {
		if err := h.mHttpServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	if h.mHttpServer != nil {
		return h.mHttpServer.Shutdown(ctx)
	}

	return nil
}

func (h *HTTPServer) TransmitInputEvent(contractId string, inputEvent *model.Event) {
	if h.GetEventTransmitter() != nil {
		go func() {
			err := h.GetEventTransmitter().TransmitInputEvent(contractId, inputEvent)
			if err != nil {
				logger.L(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}

		}()
	}
}

func (h *HTTPServer) TransmitOutputEvent(contractId string, outputEvent *model.Event) {
	if h.GetEventTransmitter() != nil {
		go func() {
			err := h.GetEventTransmitter().TransmitOutputEvent(contractId, outputEvent)
			if err != nil {
				logger.L(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}
		}()
	}
}

func (h *HTTPServer) writeMessage(statusCode int, defaultMessage string, request *http.Request, writer http.ResponseWriter, errLocal error) {
	if errLocal != nil {
		logger.L(h.ContractId()).Error(errLocal.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}

	writer.Header().Add("Content-Type", h.mDefaultContentType)
	writer.WriteHeader(statusCode)
	if _, err := writer.Write([]byte(h.getMessage(fmt.Sprintf("s%dm", statusCode), defaultMessage, h.getLanguage(request)))); err != nil {
		logger.L(h.ContractId()).Error(err.Error(),
			zap.String("version", h.Version()),
			zap.String("name", h.Name()),
			zap.String("contract_id", h.ContractId()))
	}
}

func (h *HTTPServer) s401m(request *http.Request, writer http.ResponseWriter, errLocal error) {
	responseStatus.WithLabelValues(request.URL.Path, "401").Inc()
	h.writeMessage(401, h.d401m, request, writer, errLocal)
}

func (h *HTTPServer) s403m(request *http.Request, writer http.ResponseWriter, errLocal error) {
	responseStatus.WithLabelValues(request.URL.Path, "403").Inc()
	h.writeMessage(403, h.d403m, request, writer, errLocal)
}

func (h *HTTPServer) s404m(request *http.Request, writer http.ResponseWriter, errLocal error) {
	responseStatus.WithLabelValues(request.URL.Path, "404").Inc()
	h.writeMessage(404, h.d404m, request, writer, errLocal)
}

func (h *HTTPServer) s405m(request *http.Request, writer http.ResponseWriter, errLocal error) {
	responseStatus.WithLabelValues(request.URL.Path, "405").Inc()
	h.writeMessage(405, h.d405m, request, writer, errLocal)
}

func (h *HTTPServer) s408m(request *http.Request, writer http.ResponseWriter, errLocal error) {
	responseStatus.WithLabelValues(request.URL.Path, "408").Inc()
	h.writeMessage(408, h.d408m, request, writer, errLocal)
}

func (h *HTTPServer) s499m(request *http.Request, writer http.ResponseWriter, errLocal error) {
	responseStatus.WithLabelValues(request.URL.Path, "499").Inc()
	h.writeMessage(499, h.d499m, request, writer, errLocal)
}

func (h *HTTPServer) s500m(request *http.Request, writer http.ResponseWriter, errLocal error) {
	responseStatus.WithLabelValues(request.URL.Path, "500").Inc()
	h.writeMessage(500, h.d500m, request, writer, errLocal)
}

func (h *HTTPServer) debugMessage(request *http.Request) {
	logger.L(h.ContractId()).Debug("request local timeout in seconds", zap.Duration("timeout", h.mRequestTimeout))
	logger.L(h.ContractId()).Debug("request started")
	logger.L(h.ContractId()).Debug("request data",
		zap.String("path", request.URL.Path),
		zap.String("method", request.Method),
		zap.String("path_with_query", request.RequestURI))
}

func (h *HTTPServer) AddService(
	authorizer iface.IAuthorizer,
	authorizerExpression string,
	triggerValues model.ConfigMap,
	service iface.IService) error {

	logger.L(h.ContractId()).Debug("service add",
		zap.Any("authorizer", authorizer),
		zap.Any("expression", authorizerExpression),
		zap.Any("triggerValues", triggerValues))

	var method string
	var path string
	var methodList []string

	if method = triggerValues.String("method", ""); len(method) == 0 {
		return ErrMethodNotDefined
	}

	method = strings.ToUpper(strings.TrimSpace(method))
	methodList = strings.Split(method, ",")
	if len(methodList) > 0 {
		sort.Strings(methodList)
	}

	if path = triggerValues.String("path", ""); len(path) == 0 {
		return ErrPathNotDefined
	}

	path = strings.TrimSpace(path)

	requestHandler := func(writer http.ResponseWriter, request *http.Request) {
		var err error
		timerStart := time.Now()

		defer func() {
			logger.L(h.ContractId()).Debug("request completed")
			elapsed := time.Since(timerStart)
			logger.L(h.ContractId()).Debug("request execution time", zap.Duration("seconds", elapsed))
		}()

		defer func() {
			if r := recover(); r != nil {
				panicMsg := fmt.Sprintf("%v", r)
				logger.L(h.ContractId()).Info("recovering from panic")

				// add as much information as possible
				logger.L(h.ContractId()).Error("panic data",
					zap.String("host_name", request.URL.Hostname()),
					zap.String("host", request.URL.Host),
					zap.String("path", request.URL.Path),
					zap.String("method", request.Method),
					zap.String("uri", request.RequestURI),
					zap.String("panic_msg", panicMsg))

				go func() {
					panicCounter.WithLabelValues(h.ContractId()).Inc()
				}()

				h.s500m(request, writer, nil)
				return
			}
		}()

		h.debugMessage(request)

		if !utility.IsIn(methodList, request.Method) {
			h.s405m(request, writer, nil)
			return
		}

		var data []byte

		headers := make(map[string]string)

		metadata := &model.Metadata{}
		metadata.Method = request.Method
		metadata.Path = request.URL.EscapedPath()
		metadata.Headers = make(map[string]string)
		metadata.Query = make(map[string]string)
		metadata.ContractIdList = append(metadata.ContractIdList, h.ContractId())

		for k, v := range request.Header {
			if len(v) > 0 {
				metadata.Headers[k] = v[0]
				headers[strings.ToLower(strings.TrimSpace(k))] = v[0]
			}
		}

		for k, v := range request.URL.Query() {
			if len(v) > 0 {
				metadata.Query[k] = v[0]
			}
		}

		if authorizer != nil {
			if !authorizer.IsAuthorized(authorizerExpression, metadata) {
				h.s403m(request, writer, nil)
				return
			}
		}

		if data, err = ioutil.ReadAll(request.Body); err != nil {
			h.s500m(request, writer, err)
			return
		}

		inputEvent := &model.Event{
			Metadata: metadata,
			TypeUrl:  utility.GetValue(headers, "content-type", "application/text"),
			Value:    data,
		}

		// transmit input event
		h.TransmitInputEvent(service.ContractId(), inputEvent)

		nCtx, cancel := context.WithTimeout(request.Context(), h.mRequestTimeout)
		defer cancel()

		ch := make(chan EventResponse, 1)

		func() {
			if request.Context().Err() != nil {
				ch <- EventResponse{
					Event: nil,
					Error: request.Context().Err(),
				}
			} else {
				go func() {
					event, errInner := service.Serve(nCtx, inputEvent)
					ch <- EventResponse{Event: event, Error: errInner}
				}()
			}
		}()

		select {
		case <-nCtx.Done():
			h.s408m(request, writer, nil)
			return
		case r := <-ch:
			if r.Error == context.DeadlineExceeded {
				h.s408m(request, writer, r.Error)
				return
			}

			if r.Error == context.Canceled {
				h.s499m(request, writer, r.Error)
				return
			}

			if r.Error != nil {
				h.s500m(request, writer, r.Error)
				return
			}

			// NOTE: PROMETHEUS RESPONSE STATISTICS
			go func() {
				responseStatus.WithLabelValues(request.URL.Path,
					fmt.Sprintf("%d", r.Event.Metadata.StatusCode)).Inc()
			}()

			// transmit output event
			h.TransmitOutputEvent(service.ContractId(), r.Event)

			// NOTE: handle success from service
			for k, v := range r.Event.Metadata.Headers {
				writer.Header().Add(k, v)
			}

			writer.WriteHeader(int(r.Event.Metadata.StatusCode))

			if _, err = writer.Write(r.Event.Value); err != nil {
				logger.L(h.ContractId()).Error(err.Error(),
					zap.String("version", h.Version()),
					zap.String("name", h.Name()),
					zap.String("contract_id", h.ContractId()))
			}
		}
	}

	h.mHttpServerMux.HandleFunc(path, requestHandler)

	return nil
}

func init() {
	prometheus.MustRegister(panicCounter, responseStatus)
	registry.GlobalRegistry().AddCapability(&HTTPServer{})
}
