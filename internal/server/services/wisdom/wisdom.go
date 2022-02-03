package wisdom

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"

	"wisdom/internal/pkg/pow"
	"wisdom/internal/server/config"
	"wisdom/internal/server/logger"
	"wisdom/internal/server/signature"
)

//Wisdom implementation of service for gets wisdom
type Wisdom interface {
	Start() error
	Stop(ctx context.Context) error
}

type wisdom struct {
	config config.Loader
	server *http.Server
	logger *logger.Logger
}

//NewWisdomService gets service which implements http server
func NewWisdomService(c config.Loader, l *logger.Logger) Wisdom {
	service := &wisdom{
		config: c,
		logger: l,
	}

	mux := http.NewServeMux()
	mux.Handle("/wisdom", service.signatureMiddleware(service.powMiddleware(service.wisdomHandler())))

	service.server = &http.Server{
		Handler: mux,
	}

	return service
}

// check implementation
var _ Wisdom = (*wisdom)(nil)

//Start http server
func (wdm *wisdom) Start() error {
	address := wdm.config.Get().App.HTTPWisdomAddr
	if address == "" {
		address = defaultWisdomAddr
	}

	ln, err := net.Listen("tcp", address)

	if err != nil {
		return errors.Wrap(err, "can't create listener")
	}

	if err := wdm.server.Serve(ln); err != nil {
		return errors.Wrap(err, "can't start server")
	}
	return nil
}

//Stop safely shut down for server
func (wdm *wisdom) Stop(ctx context.Context) error {
	return wdm.server.Shutdown(ctx)
}

func (wdm *wisdom) wisdomHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := wdm.getRandomWisdom()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			_, _ = fmt.Fprintf(w, "error encode response: %+v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func (wdm *wisdom) signatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signatureService := signature.NewSignatureService(wdm.config)
		if r.Method == http.MethodGet {
			sign := signatureService.Generate()
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(sign); err != nil {
				wdm.logger.Error().Msgf("error encode response: %+v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		sign := r.Header.Get("Salt")
		if sign == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "signature must be send")
			return
		}
		signIsOk, err := signatureService.Validate(sign)
		if err != nil {
			wdm.logger.Error().Msgf("signature validation error: %+v", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
		if !signIsOk {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

func (wdm *wisdom) powMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		block := &pow.Block{}
		if err := json.NewDecoder(r.Body).Decode(block); err != nil {
			wdm.logger.Error().Msgf("block couldn't be parsed: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := pow.NewPoofOfWork(block)
		if p.Validate() == false {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
