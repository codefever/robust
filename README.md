## Robust

**Robust** makes standalone services to become in cold-standby mode via etcd.

### Usage

```bash
GOPROXY=https://goproxy.io go build
./robust --endpoints ${ETCD_ENDPOINTS} --command './docker/looper.sh 1' --name ${NAME}
```

You can set up a local cluster for etcd as below.

```bash
cd docker && docker-compose up
```
