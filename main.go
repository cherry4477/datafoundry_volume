package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/golang/glog"

	"github.com/asiainfoLDP/datafoundry_volume/openshift"
)

var Debug = false

type mux struct{}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)
	RespError(w, fmt.Errorf("not found"), http.StatusNotFound)
}

func main() {

	go openshift.Init(
		os.Getenv("DATAFOUNDRY_HOST_ADDR"),
		os.Getenv("DATAFOUNDRY_ADMIN_USER"),
		os.Getenv("DATAFOUNDRY_ADMIN_PASS"),
	)

	// ...

	router := httprouter.New()

	router.POST("/lapi/v1/namespaces/:namespace/volumes", CreateVolume)
	router.DELETE("/lapi/v1/namespaces/:namespace/volumes/:name", DeleteVolume)

	router.NotFound = &mux{}
	router.MethodNotAllowed = &mux{}
	glog.Infoln("listening on port 9095")
	glog.Fatalln(http.ListenAndServe(":9095", router))
}

