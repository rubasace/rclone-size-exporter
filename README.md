## Description
Rclone Size Exporter generates [Prometheus](https://prometheus.io/) metrics for getting size and number of objects on [Rclone](https://rclone.org/) remotes.

Under the cover, it mainly executes a `rclone size {remote}` and stores the response, so it can be scrapped by Prometheus.

## Why
Rclone supports Prometheus metrics, but it doesn't expose remote sizes or number of objects out-of-the-box. This exporter makes it possible to get that information from any remote configured on Rclone, even if it isn't mounted.


## Metrics example
It provides all standard golang metrics regarding memory usage, gc, etc. Rclone metrics are as follow:

```
# HELP rclone_size_exporter_connection_error Flags if there are current issues connecting to rclone.
# TYPE rclone_size_exporter_connection_error gauge
rclone_size_exporter_connection_error{remote="drive:"} 0
# HELP rclone_size_exporter_objects_number Number of elements on the remote volume.
# TYPE rclone_size_exporter_objects_number gauge
rclone_size_exporter_objects_number{remote="drive:"} 706
# HELP rclone_size_exporter_total_size Size of the remote volume.
# TYPE rclone_size_exporter_total_size gauge
rclone_size_exporter_total_size{remote="drive:"} 4.3483827897038e+10
```


## Parameters (Passed as Environment Variables)

| Name | Default | Description |  |
| ---- | ---- | ----------- | -------- |
| PORT | `8080` (On Docker) | Port where the exporter will be serving the `/metrics` endpoint. | &nbsp; |
| RCLONE_CONFIG  | `/config/rclone/rclone.conf` (On Docker) | Location of the Rclone config file  | &nbsp; |
| DELAY | `300`  | Time, in seconds, between executions of `rclone size` to retrieve new values | &nbsp; |
| REMOTE | No Default | Name of the remote to retrieve info from (`drive:`). It can also be a remote + path (`drive:/folder1`). Colon is important! | &nbsp; |

## Run in Docker
```bash
docker run -d --name rclone-size-exporter \
        --restart=unless-stopped \
        -v /path/to/config/rclone.conf:/config/rclone/rclone.conf \
        -e PORT=8181 \
        -e RCLONE_CONFIG=/config/rclone/rclone.conf \
        -e REMOTE="drive:" \
        -e DELAY=600 \
        rubasace/rclone-size-exporter
```
## Requirements
* Go lang 1.14+
* Rclone