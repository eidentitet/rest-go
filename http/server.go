package http

import (
	configuration "github.com/eidentitet/rest-go/config"
	chttp "github.com/eidentitet/rest-go/http/meta"
	"github.com/eidentitet/rest-go/http/middleware"
	"github.com/eidentitet/rest-go/tls"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ApiHandler struct {
	Routes  []chttp.Route
	Mux     *mux.Router
	Version string
	//Config
}

//type Config struct {
//	OpenID OpenID `fig:"openid" validate:"required"`
//	Server Server `fig:"server" validate:"required"`
//}

func (api ApiHandler) Start() {
	appConfig := configuration.GetAppConfig()

	r := mux.NewRouter()
	api.Mux = r

	// Add your routes as needed
	api.buildPostRoutes()

	// Start HTTP Server
	switch appConfig.Server.TLS {
	case true:
		server := tls.GetHTTPSServer(appConfig.Server.Name, appConfig.Server.Addr)
		server.Handler = r
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Errorln(err)
		}
		break
	case false:
		srv := &http.Server{
			Addr: appConfig.Server.Addr,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: appConfig.Server.GracefulTimeout,
			ReadTimeout:  appConfig.Server.GracefulTimeout,
			IdleTimeout:  time.Second * 60,
			Handler:      r,
		}

		// Run our server in a goroutine so that it doesn't block.
		if err := srv.ListenAndServe(); err != nil {
			log.Errorln(err)
		}
		break
	}
}

// Set up routes
func (api ApiHandler) buildPostRoutes() {
	pluginManager := middleware.GetMiddlewareManager()
	for _, route := range api.Routes {
		apiVx := api.Mux.PathPrefix("/api/" + api.Version).Subrouter()
		middlewares, _ := pluginManager.GetMiddleware(route.Middlewares...)
		apiVx.Use(middlewares...)
		apiVx.HandleFunc(route.Path, route.Handler).Methods(route.Method).Name(route.Name)
		log.Infoln(route)
	}
}
