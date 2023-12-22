package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/attribute"
)

func HandleHealthz(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tr := traceProvider.Tracer("videos-ms-main-component")
	id := uuid.New()
	ip := strings.Split(r.RemoteAddr, ":")[0]
	_, span := tr.Start(context.Background(), "healthz")
	span.SetAttributes(attribute.Key("Protocol").String(r.Proto))
	span.SetAttributes(attribute.Key("UUID").String(id.String()))
	span.SetAttributes(attribute.Key("Client IP").String(ip))
	defer span.End()
	Sugar.Infof("client_ip: %v", ip)
	Sugar.Infof("request_id: %v", id)
	fmt.Fprintf(w, "ok!")
}

func HandleGetVideoById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	tr := traceProvider.Tracer("videos-ms-main-component")
	ctx, span := tr.Start(r.Context(), "GET Video By ID")
	id := uuid.New()
	ip := strings.Split(r.RemoteAddr, ":")[0]
	Sugar.Infof("client_ip: %v", ip)
	Sugar.Infof("request_id: %v", id)

	defer span.End()

	span.SetAttributes(attribute.Key("Function").String("GetVideosByIDHandler"))
	span.SetAttributes(attribute.Key("Protocol").String(r.Proto))
	span.SetAttributes(attribute.Key("UUID").String(id.String()))
	span.SetAttributes(attribute.Key("Client IP").String(strings.Split(r.RemoteAddr, ":")[0]))

	if flaky == "true" {
		if rand.Intn(90) < 30 {
			span.RecordError(errors.New("Flaky Error"))
			panic("flaky error occurred ")
		}
	}

	video := video(w, r, p, ctx, id, ip)

	Cors(w)
	fmt.Fprintf(w, "%s", video)
}
