package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiali/kiali/log"
)

// Get JaegerInfo provides the Jaeger URL and other info, first by checking if a config exists
// then (if not) by inspecting the Kubernetes Jaeger service in Istio installation namespace
func GetJaegerInfo(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Token initialization error: "+err.Error())
		return
	}

	info, code, err := business.Jaeger.GetJaegerInfo()
	if err != nil {
		log.Error(err)
		RespondWithError(w, code, err.Error())
		return
	}
	RespondWithJSON(w, code, info)
}

func GetJaegerServices(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	services, code, err := business.Jaeger.GetJaegerServices()
	if err != nil {
		log.Error(err)
		RespondWithError(w, code, err.Error())
		return
	}
	RespondWithJSON(w, code, services)
}

func TraceServiceDetails(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Trace Service Details initialization error: "+err.Error())
		return
	}
	params := mux.Vars(r)
	namespace := params["namespace"]
	service := params["service"]
	traces, code, err := business.Jaeger.GetJaegerTraces(namespace, service, r.URL.RawQuery)
	if err != nil {
		log.Error(err)
		RespondWithError(w, code, err.Error())
		return
	}
	RespondWithJSON(w, code, traces)
}

func TraceDetails(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Trace Detail initialization error: "+err.Error())
		return
	}
	params := mux.Vars(r)
	traceID := params["traceID"]
	traces, code, err := business.Jaeger.GetJaegerTraceDetail(traceID)
	if err != nil {
		log.Error(err)
		RespondWithError(w, code, err.Error())
		return
	}
	RespondWithJSON(w, code, traces)
}

// ServiceSpans is the API handler to fetch Jaeger spans of a specific service
func ServiceSpans(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	params := mux.Vars(r)
	namespace := params["namespace"]
	service := params["service"]
	queryParams := r.URL.Query()
	startMicros := queryParams.Get("startMicros")
	endMicros := queryParams.Get("endMicros")

	spans, err := business.Jaeger.GetJaegerSpans(namespace, service, startMicros, endMicros)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, spans)
}
