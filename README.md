
## Overview

管理员创建/删除pvc

## APIs

### POST /lapi/v1/namespaces/:namespace/volumes

管理员创建pvc

Path Parameters:
```
namespace: 
```

Body Parameters (json):
```
name: pvc名称。
size: pvc大小。
```

### DELETE /lapi/v1/namespaces/:namespace/volumes/:name

管理员删除pvc

Path Parameters:
```
namespace: 
name: pvc名称。
```

## 部署

```
oc new-app --name datafoundryservicevolume https://github.com/asiainfoLDP/datafoundry_volume.git#develop \
    -e  GLUSTER_ENDPOINTS_NAME="xxx" \
    \
    -e  HEKETI_HOST_ADDR="xxx" \
    -e  HEKETI_HOST_PORT="xxx" \
    -e  HEKETI_USER="xxx" \
    -e  HEKETI_KEY="xxx" \
    \
    -e  DATAFOUNDRY_HOST_ADDR="xxx" \
    -e  DATAFOUNDRY_ADMIN_USER="xxx" \
    -e  DATAFOUNDRY_ADMIN_PASS="xxx"

oc expose service datafoundryservicevolume --hostname=datafoundry.servicevolume.app.dataos.io

oc start-build datafoundryservicevolume

```



