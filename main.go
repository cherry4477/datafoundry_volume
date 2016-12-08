package main

import (
	"fmt"
	"net/http"
	"os"
	"flag"

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

func startHttpServer() {

	// ...

	router := httprouter.New()

	router.POST("/lapi/v1/namespaces/:namespace/volumes", CreateVolume)
	router.DELETE("/lapi/v1/namespaces/:namespace/volumes/:name", DeleteVolume)

	router.GET("/lapi/v1/volumes", QueryVolumes)
	router.PUT("/lapi/v1/volumes", ManageVolumes)

	router.NotFound = &mux{}
	router.MethodNotAllowed = &mux{}
	glog.Infoln("listening on port 9095")
	glog.Fatalln(http.ListenAndServe(":9095", router))
}

var debug = flag.Bool("debug", false, "is debug mode?")
var cli = flag.Bool("cli", false, "used as cli command?")

func main() {
	openshift.Init(
		os.Getenv("DATAFOUNDRY_HOST_ADDR"),
		os.Getenv("DATAFOUNDRY_ADMIN_USER"),
		os.Getenv("DATAFOUNDRY_ADMIN_PASS"),
	)

	InitGluster(
		os.Getenv("GLUSTER_ENDPOINTS_NAME"),

		os.Getenv("HEKETI_HOST_ADDR"),
		os.Getenv("HEKETI_HOST_PORT"),
		os.Getenv("HEKETI_USER"),
		os.Getenv("HEKETI_KEY"),
	)

	
	flag.Parse()

	if *cli {
		executeCommand()
	} else {
		startHttpServer()
	}
}

