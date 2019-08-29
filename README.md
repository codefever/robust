## Robust

**Robust** makes standalone services to become in cold-standby mode via etcd.

### Usage

```bash
./robust --endpoints ${ETCD_ENDPOINTS} --command './docker/looper.sh 1' --name ${NAME}
```
