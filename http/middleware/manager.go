package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

var pluginManager Manager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetMiddlewareManager()
	})
}

type (
	// ManagerImpl is a manager for the plugins. It stores and maps the plugins to the command.
	ManagerImpl struct {
		mws sync.Map
	}

	Manager interface {
		HasMiddleware(name string) bool
		AddMiddleware(name string, middleware Handler)
		GetMiddleware(name ...string) ([]mux.MiddlewareFunc, error)
	}

	Handler interface {
		Next() func(handler http.Handler) http.Handler
	}
)

// Register a cmd to the manager.
func Register(name string, plugin Handler) {
	GetMiddlewareManager().AddMiddleware(name, plugin)
}

// GetPluginManager gets the cmd manager instance (singleton).
func GetMiddlewareManager() Manager {
	if pluginManager == nil {
		pluginManager = &ManagerImpl{mws: sync.Map{}}
	}
	return pluginManager
}

// HasMiddleware Check if the cmd exists in the manager.
func (pm *ManagerImpl) HasMiddleware(name string) bool {
	_, exists := pm.mws.Load(name)
	return exists
}

// AddMiddleware Add a middleware to the manager.
func (pm *ManagerImpl) AddMiddleware(name string, plugin Handler) {
	if !pm.HasMiddleware(name) {
		pm.mws.Store(name, plugin)
	}
}

// GetMiddleware returns the middleware, if found.
func (pm *ManagerImpl) GetMiddleware(name ...string) ([]mux.MiddlewareFunc, error) {
	var handlers []mux.MiddlewareFunc
	for _, n := range name {
		mPlugin, exists := pm.mws.Load(n)
		if exists {
			middlew := mPlugin.(Handler)
			handlers = append(handlers, middlew.Next())
		}
	}
	return handlers, nil
}
