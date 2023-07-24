# kube-exec-perf-test
Performance test for kube exec api

## history
```shell
go mod init github.com/xcheng85/kube-exec-perf-test

go get k8s.io/client-go@latest

go mod tidy

go build ./cmds/exec/
# Authenticating outside the cluster
```

## reference
https://github.com/zhiminwen/MQQDScraper/blob/master/qde