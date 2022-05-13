# TCP Proxy

## Usage

```
tp REMOTE_ADDRESS [LOCAL_ADDRESS]
```

## Examples

```
tp postgres-host.lan:5432                        listens to 127.0.0.1:5432
tp postgres-host.lan:5432 :8888                  listens to localhost:8888
tp postgres-host.lan:5432 [::1]:8888             listens to [::1]:8888
```
