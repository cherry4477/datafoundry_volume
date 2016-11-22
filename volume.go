package main

import (
	//"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"encoding/base32"
	"os"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"

	"github.com/asiainfoLDP/datafoundry_volume/openshift"
	"github.com/asiainfoLDP/datahub_commons/common"
	kapiresource "k8s.io/kubernetes/pkg/api/resource"
	kapi "k8s.io/kubernetes/pkg/api/v1"

	heketi "github.com/heketi/heketi/client/api/go-client"
	"github.com/heketi/heketi/pkg/glusterfs/api"
)

const (
	MinVolumnSize = 10
	MaxVolumnSize = 200

	Gi = int64(1) << 30
)

var invalidVolumnSize = fmt.Errorf(
	"volumn size must be integer multiple of 10G and in range [%d, %d].",
	MinVolumnSize, MaxVolumnSize)

func heketiClient() *heketi.Client {
	return heketi.NewClient(
		fmt.Sprintf("http://%s:%s",
			os.Getenv("HEKETI_HOST_ADDR"),
			os.Getenv("HEKETI_HOST_PORT"),
		),
		os.Getenv("HEKETI_USER"),
		os.Getenv("HEKETI_KEY"),
	)
}

func glusterEndpointsName() string {
	return os.Getenv("GLUSTER_ENDPOINTS_NAME")
}

//==============================================================
//
//==============================================================

//func NewElevenLengthID() string {
//	t := time.Now().UnixNano()
//	bs := make([]byte, 8)
//	for i := uint(0); i < 8; i ++ {
//		bs[i] = byte((t >> i) & 0xff)
//	}
//	return string(base64.RawURLEncoding.EncodeToString(bs))
//}

var base32Encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")
func NewThirteenLengthID() string {
	t := time.Now().UnixNano()
	bs := make([]byte, 8)
	for i := uint(0); i < 8; i ++ {
		bs[i] = byte((t >> i) & 0xff)
	}
	
	dest := make([]byte, 16)
	base32Encoding.Encode(dest, bs)
	return string(dest[:13])
}

//func PvcName2PvName(namespace, volName string) string {
//	return fmt.Sprintf("%s-%s", namespace, volName) // don't change
//}

// volSource: "gluster", ...
// volSize: GB
func BuildRandomPvName(volSource string, volSize int) string {
	return fmt.Sprintf("%s-%dg-%s", strings.ToLower(volSource), volSize, NewThirteenLengthID())
}

func VolumeId2VolumeName(volId string) string {
	return "vol_" + volId
}

func VolumeName2VolumeId(volName string) string {
	const prefix = "vol_"
	if strings.HasPrefix(volName, prefix) {
		return volName[len(prefix):]
	}

	return volName
}

//func PvName2VolumeName(namespace, pvName string) string {
//	prefix := PvcName2PvName(pvName, "")
//	if strings.HasPrefix(pvName, prefix) {
//		return pvName[len(prefix):]
//	}
//
//	return pvName
//}

//==============================================================
//
//==============================================================

//==============================================================
//
//==============================================================

func CreateVolume(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	glog.Infoln("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	var username string
	var err error

	if username, err = getDFUserame(r.Header.Get("Authorization")); err != nil {
		RespError(w, err, http.StatusUnauthorized)
		return
	}

	// params

	namespace, e := MustStringParamInPath(params, "namespace", StringParamType_UrlWord)
	if e != nil {
		RespError(w, e, http.StatusBadRequest)
		return
	}

	m, err := common.ParseRequestJsonAsMap(r)
	if err != nil {
		glog.Error(err)
		RespError(w, err, http.StatusBadRequest)
		return
	}

	size, e := MustIntParamInMap(m, "size")
	if e != nil {
		RespError(w, e, http.StatusBadRequest)
		return
	}
	if size < MinVolumnSize || size > MaxVolumnSize || (size%10) != 0 {
		RespError(w, invalidVolumnSize, http.StatusBadRequest)
		return
	}

	pvcname, e := MustStringParamInMap(m, "name", StringParamType_UrlWord)
	if e != nil {
		RespError(w, e, http.StatusBadRequest)
		return
	}
	valid, msg := NameIsDNSLabel(pvcname, false)
	if !valid {
		RespError(w, errors.New(msg), http.StatusBadRequest)
		return
	}

	// todo: check permission

	// _ = username
	// _ = namespace

	// ...

	resourceList := make(kapi.ResourceList)
	resourceList[kapi.ResourceStorage] = *kapiresource.NewQuantity(int64(size*Gi), kapiresource.BinarySI)

	// create volumn
	hkiClient := heketiClient()

	_, _ = hkiClient.ClusterList() // test
	clusterlist, err := hkiClient.ClusterList()
	if err != nil {
		glog.Error(err)
		RespError(w, err, http.StatusInternalServerError)
		return
	}

	go func() {

		req := &api.VolumeCreateRequest{}
		req.Size = int(size)
		//req.Name = pvcname + "-" + namespace + "-" + username + "-jd" // ! don't set name, otherwise, can't get volume id from pv

		req.Clusters = clusterlist.Clusters //[]string{"68aa170df797272ac2ac90fac1f7460b"} //hacked by san
		req.Durability.Type = api.DurabilityReplicate
		req.Durability.Replicate.Replica = 3
		req.Durability.Disperse.Data = 4
		req.Durability.Disperse.Redundancy = 2

		// if snapshotFactor > 1.0 {
		//	req.Snapshot.Factor = float32(snapshotFactor)
		//	req.Snapshot.Enable = true
		// }

		var succeeded = false

		glog.Warningf("creating volume by %s@%s", username, namespace)
		volume, err := hkiClient.VolumeCreate(req)
		if err != nil {
			glog.Error("create volume by heketi error:", err)
			//RespError(w, err, http.StatusBadRequest)
			return
		}

		defer func() {
			if succeeded {
				//glog.Infoln("success")
				return
			}

			err := hkiClient.VolumeDelete(volume.Id)
			if err != nil {
				glog.Warningf("delete volume (%s, %s) on failed to CreateVolume", pvcname, volume.Id)
			}
		}()

		// create pv

		inputPV := &kapi.PersistentVolume{}
		{
			inputPV.Kind = "PersistentVolume"
			inputPV.APIVersion = "v1"
			inputPV.Annotations = make(map[string]string)
			inputPV.Annotations["datafoundry.io/gluster-volume"] = volume.Id
			inputPV.Annotations["datafoundry.io/requester"] = username + "@" + namespace
			//inputPV.Name = PvcName2PvName(namespace, pvcname) // this is not a good idea, kubernetes may not math pv and pvc with the same name
			inputPV.Name = BuildRandomPvName("gluster", int(size))
			inputPV.Spec.Capacity = resourceList
			inputPV.Spec.PersistentVolumeSource = kapi.PersistentVolumeSource{
				Glusterfs: &kapi.GlusterfsVolumeSource{
					EndpointsName: glusterEndpointsName(),
					Path:          VolumeId2VolumeName(volume.Id),
				},
			}
			inputPV.Spec.AccessModes = []kapi.PersistentVolumeAccessMode{
				kapi.ReadWriteMany,
			}
			inputPV.Spec.PersistentVolumeReclaimPolicy = kapi.PersistentVolumeReclaimRecycle
		}

		outputPV := &kapi.PersistentVolume{}
		osrPV := openshift.NewOpenshiftREST(nil)
		osrPV.KPost("/persistentvolumes", inputPV, outputPV)
		if osrPV.Err != nil {
			glog.Warningf("create pv error CreateVolume: pvname=%s, error: %s", inputPV.Name, osrPV.Err)

			//RespError(w, osrPV.Err, http.StatusBadRequest)
			return
		}
		defer func() {
			if succeeded {
				//glog.Infoln("success")
				return
			}

			osrPV := openshift.NewOpenshiftREST(nil)
			osrPV.KDelete("/persistentvolumes/"+inputPV.Name, nil)
			if osrPV.Err != nil {
				glog.Warningf("delete pv error on failed to CreateVolume: pvname=%s, error: %s", inputPV.Name, osrPV.Err)
			}
		}()

		succeeded = true
		glog.Infof("create volume(%s) by %s@%s successfuly.", volume.Id, username, namespace)

		// update pvc

		currentPVC := &kapi.PersistentVolumeClaim{}
		osrPVC := openshift.NewOpenshiftREST(openshift.NewOpenshiftClient(r.Header.Get("Authorization")))
		osrPVC.KGet("/namespaces/"+namespace+"/persistentvolumeclaims/"+pvcname, currentPVC)
		if osrPVC.Err != nil {
			glog.Warningf("get pvc error on failed to CreateVolume: pvname=%s, error: %s", pvcname, osrPVC.Err)
			return
		}

		if currentPVC.Annotations == nil {
			currentPVC.Annotations = make(map[string]string)
		}
		currentPVC.Annotations["datafoundry.io/gluster-volume"] = volume.Id
		currentPVC.Annotations["datafoundry.io/requester"] = username + "@" + namespace

		osrPVC.KPut("/namespaces/"+namespace+"/persistentvolumeclaims/"+pvcname, &currentPVC, nil)
		if osrPVC.Err != nil {
			glog.Warningf("update pvc error on CreateVolume: pvcname=%s, error: %s", pvcname, osrPVC.Err)
			return
		}
	}()

	// create pvc

	inputPVC := &kapi.PersistentVolumeClaim{}
	{
		inputPVC.Kind = "PersistentVolumeClaim"
		inputPVC.APIVersion = "v1"
		inputPVC.Name = pvcname
		inputPVC.Spec.AccessModes = []kapi.PersistentVolumeAccessMode{
			kapi.ReadWriteMany,
		}
		inputPVC.Spec.Resources = kapi.ResourceRequirements{
			Requests: resourceList,
		}
	}

	outputPVC := &kapi.PersistentVolumeClaim{}
	osrPVC := openshift.NewOpenshiftREST(openshift.NewOpenshiftClient(r.Header.Get("Authorization")))
	osrPVC.KPost("/namespaces/"+namespace+"/persistentvolumeclaims", &inputPVC, &outputPVC)
	if osrPVC.Err != nil {
		glog.Warningf("create pvc error on CreateVolume: pvcname=%s, error: %s", pvcname, osrPVC.Err)

		//RespError(w, osrPVC.Err, http.StatusBadRequest)
		return
	}

	//defer func() {
	//	if succeeded {
	//		return
	//	}
	//
	//	osrPVC := openshift.NewOpenshiftREST(openshift.NewOpenshiftClient(r.Header.Get("Authorization")))
	//	osrPVC.KDelete("/namespaces/"+namespace+"/persistentvolumeclaims/"+inputPVC.Name, nil)
	//	if osrPVC.Err != nil {
	//		glog.Warningf("delete pvc error on failed to CreateVolume: pvcname=%s, error: %s", pvcname, osrPVC.Err)
	//	}
	//}()

	// ...

	RespOK(w, outputPVC)
}

func DeleteVolume(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	glog.Infoln("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	var username string
	var err error

	if username, err = getDFUserame(r.Header.Get("Authorization")); err != nil {
		RespError(w, err, http.StatusUnauthorized)
		return
	}

	// ...

	namespace, e := MustStringParamInPath(params, "namespace", StringParamType_UrlWord)
	if e != nil {
		RespError(w, e, http.StatusBadRequest)
		return
	}

	pvcname, e := MustStringParamInPath(params, "name", StringParamType_UrlWord)
	if e != nil {
		RespError(w, e, http.StatusBadRequest)
		return
	}
	valid, msg := NameIsDNSLabel(pvcname, false)
	if !valid {
		RespError(w, errors.New(msg), http.StatusBadRequest)
		return
	}

	// todo: check permission

	_ = username
	_ = namespace

	// get pv (will delete it at the end, for it stores the volumn id info)


	// get pvc (to get the pv name)
	outputPVC := &kapi.PersistentVolumeClaim{}
	osrGetPVC := openshift.NewOpenshiftREST(openshift.NewOpenshiftClient(r.Header.Get("Authorization")))
	osrGetPVC.KGet("/namespaces/" + namespace+"/persistentvolumeclaims/"+pvcname, &outputPVC)
	if osrGetPVC.Err != nil {
		glog.Infof("get pvc error: pvcname=%s, error: %s", pvcname, osrGetPVC.Err)
		RespError(w, osrGetPVC.Err, http.StatusBadRequest)
		return
	}

	// bug: PvcName2PvName may return a wrong pv
	//pvName := PvcName2PvName(namespace, pvcname)
	// use this insead
	pvName := outputPVC.Spec.VolumeName

	// delete pvc

	// func() {
	osrDeletePVC := openshift.NewOpenshiftREST(openshift.NewOpenshiftClient(r.Header.Get("Authorization")))
	osrDeletePVC.KDelete("/namespaces/"+namespace+"/persistentvolumeclaims/"+pvcname, nil)
	if osrDeletePVC.Err != nil {
		glog.Infof("delete pvc error: pvcname=%s, error: %s", pvcname, osrDeletePVC.Err)
		RespError(w, osrDeletePVC.Err, http.StatusBadRequest)
		return
	}
	// }()

	// delete volume

	if pvName != "" {
		// todo: do it on a task server
		go func() {
			//get pv
			pv := &kapi.PersistentVolume{}
			osrGetPV := openshift.NewOpenshiftREST(nil)
			osrGetPV.KGet("/persistentvolumes/"+pvName, pv)
			if osrGetPV.Err != nil {
				glog.Warningf("get pv %s info error:%v", pvName, osrGetPV.Err)
				//RespError(w, osrGetPV.Err, http.StatusBadRequest)
				return
			}
			// delete pv

			osrDeletePV := openshift.NewOpenshiftREST(nil)
			osrDeletePV.KDelete("/persistentvolumes/"+pv.Name, nil)
			if osrDeletePV.Err != nil {
				// todo: retry once?

				glog.Warningf("delete pv error: pvname=%s, error: %s", pv.Name, osrDeletePV.Err)

				//RespError(w, osrDeletePV.Err, http.StatusBadRequest)
				return
			}

			hkiClient := heketiClient()

			glusterfs := pv.Spec.PersistentVolumeSource.Glusterfs
			if glusterfs != nil {
				volId := VolumeName2VolumeId(glusterfs.Path) //pv.Annotations["datafoundry.io/gluster-volume"] //
				err := hkiClient.VolumeDelete(volId)
				if err != nil {
					glog.Infof("delete volume error: pvcname=%s, volid=%s, error: %s", pvcname, volId, err)

					// todo: log it
				} else {
					glog.Info("delete volume success, volumeid:", volId)
				}
			} else {
				glog.Infof("pv.Spec.PersistentVolumeSource.Glusterfs == nil. pvcname=%s", pvcname)
			}

		}()
	}

	// ...

	RespOK(w, nil)
}
