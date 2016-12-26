package main

import (
	"github.com/asiainfoLDP/datafoundry_volume/openshift"
	userapi "github.com/openshift/origin/pkg/user/api/v1"
	//projectapi "github.com/openshift/origin/pkg/project/api/v1"
	kapi "k8s.io/kubernetes/pkg/api/v1"
	//kapiresource "k8s.io/kubernetes/pkg/api/resource"
	"github.com/golang/glog"
)

//================================================================
//
//================================================================

func authDF(userToken string) (*userapi.User, error) {
	if Debug {
		return &userapi.User{
			ObjectMeta: kapi.ObjectMeta{
				Name: "local",
			},
		}, nil
	}

	u := &userapi.User{}
	osRest := openshift.NewOpenshiftREST(openshift.NewOpenshiftClient(userToken))

	uri := "/users/~"
	osRest.OGet(uri, u)
	if osRest.Err != nil {
		glog.Infof("authDF, uri(%s) error: %s", uri, osRest.Err)
		return nil, osRest.Err
	}

	return u, nil
}

func dfUser(user *userapi.User) string {
	return user.Name
}

func getDFUserame(token string) (string, error) {
	user, err := authDF(token)
	if err != nil {
		return "", err
	}
	return dfUser(user), nil
}
