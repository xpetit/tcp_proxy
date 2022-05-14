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

## Performance

```console
user@host:~$ socat /dev/null,ignoreeof tcp-listen:1234,fork,reuseaddr &
user@host:~$ tp :1234 :2345 &
Forwarding from [::]:2345 to :1234
user@host:~$ dd if=/dev/zero status=progress bs=1M | netcat 127.0.0.1 2345
9458155520 bytes (9.5 GB, 8.8 GiB) copied, 7 s, 1.4 GB/s^C
9464+0 records in
9463+0 records out
9922674688 bytes (9.9 GB, 9.2 GiB) copied, 7.34455 s, 1.4 GB/s
user@host:~$ kill $(jobs -p)
```
