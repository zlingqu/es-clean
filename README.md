# es-harbor

清理es存储中的过期索引index。

index名称需要符合 ***YYYY-MM-DD格式，比如k8s-dev-devops-2020-09-14

支持一键清除所有。

需要harbor的定时清理任务配合。

# 1. 二进制执行

## 1.1 打包二进制
```shell
go build 
#或者
make build
```


## 1.2 查看帮助

```shell
./es-clean -h
es-clean 用于清理es中的索引，以释放存储资源

Usage:
  es-clean [flags]

Examples:
es-clean --ip 1.1.1.1  --port 9200 --indexName k8s-dev* --keepTimeDay 200

Flags:
  -h, --help               help for es-clean
  -n, --indexName string   密码
  -i, --ip string          例如：1.1.1.1
  -k, --keepTimeDay int    保留索引的天数，单位是天，比如60
  -p, --port string        端口，例如：9200
```
## 1.3 手动执行清理
Linux服务器上操作
```shell
./harbor-clean --ip ** --port ** --indexName ** --keepTimeDay **
```

## 1.4 crontab定时任务执行

```shell
# crontab -l
0 2 * 7 * /root/es-harbor --ip ** --port ** --indexName ** --keepTimeDay ** 100 >> /var/log/es-clean`date "+%Y-%m-%d-%H:%M:%S"`.log
```


# 2. docker执行


## 2.1 制作镜像
```shell
# 做镜像,比如
docker build . -t harbor.abc.com/devops/es-clean:v1
```
## 2.2 docker run形式执行
```
# 执行清理
docker run harbor.abc.com/devops/es-clean:v1 /data/es-clean --ip ** --port ** --indexName ** --keepTimeDay **
```
## 2.3 k8s中CronJob形式执行

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: es-clean
spec:
  schedule: "*/1 * * * *"
  successfulJobsHistoryLimit: 1
  failedJobsHistoryLimit: 1
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: Never
          #imagePullSecrets:
          #- name: regsecret
          containers:
          - name: es-clean
            image: harbor.abc.com/devops/es-clean:v1
            args:
            - "/bin/sh"
            - "-c"
            - "--ip ** --port ** --indexName ** --keepTimeDay **"
```


# 3. 输出内容
正常输入如下类似格式
## 3.1 删除某一类index
```shell
8s-dev* 不包含 k8s-test-cum-2020-08-31, 跳过删除...
k8s-dev* 不包含 k8s-test-cum-2020-09-02, 跳过删除...
k8s-dev* 不包含 k8s-test-cum-2020-09-01, 跳过删除...
k8s-dev* 不包含 k8s-test-cum-2020-09-04, 跳过删除...
k8s-dev* 不包含 k8s-test-cum-2020-09-03, 跳过删除...
k8s-dev* 不包含 k8s-test-cum-2020-09-06, 跳过删除...
k8s-dev* 不包含 k8s-test-cum-2020-09-05, 跳过删除...
k8s-dev* 【包含】 k8s-dev-xmc-tk-2020-08-26,距离过期还剩10.53天, 跳过删除...
k8s-dev* 不包含 k8s-prd-mis-2020-08-25, 跳过删除...
k8s-dev* 不包含 k8s-prd-mis-2020-08-24, 跳过删除...
k8s-dev* 【包含】 k8s-dev-xmc-tk-2020-08-24,距离过期还剩8.53天, 跳过删除...
k8s-dev* 【包含】 k8s-dev-xmc-tk-2020-08-25,距离过期还剩9.53天, 跳过删除...
k8s-dev* 【包含】 k8s-dev-xmc-tk-2020-09-05,距离过期还剩20.53天, 跳过删除...
k8s-dev* 【包含】 k8s-dev-xmc-tk-2020-09-06,距离过期还剩21.53天, 跳过删除...
k8s-dev* 不包含 k8s-test-mis-2020-09-14, 跳过删除...
k8s-dev* 【包含】 k8s-dev-xmc-tk-2020-09-07,距离过期还剩22.53天, 跳过删除...
```

## 3.2 使用all一键删除
```shell
all* 【包含】 k8s-test-mis-2020-09-01,距离过期还剩8.53天, 跳过删除...
all* 【包含】 k8s-test-cum-2020-09-08,距离过期还剩15.53天, 跳过删除...
all* 【包含】 k8s-uat-xmc-tk-2020-08-20,已存在25.47天,需要保留 22天,  开始删除索引k8s-uat-xmc-tk-2020-08-20
all* 【包含】 k8s-test-cum-2020-09-07,距离过期还剩14.53天, 跳过删除...
all* 【包含】 k8s-uat-xmc-tk-2020-08-21,已存在24.47天,需要保留 22天,  开始删除索引k8s-uat-xmc-tk-2020-08-21
all* 【包含】 k8s-uat-xmc-tk-2020-08-22,已存在23.47天,需要保留 22天,  开始删除索引k8s-uat-xmc-tk-2020-08-22
all* 【包含】 k8s-test-cum-2020-09-09,距离过期还剩16.53天, 跳过删除...
all* 【包含】 k8s-dev-px-mlcloud-core-2020-09-11,距离过期还剩18.53天, 跳过删除...
all* 【包含】 k8s-test-ex-dialogue-2020-09-03,距离过期还剩10.53天, 跳过删除...
all* 【包含】 k8s-prd-cum-2020-09-07,距离过期还剩14.53天, 跳过删除...
all* 【包含】 k8s-prd-cum-2020-09-08,距离过期还剩15.53天, 跳过删除...
all* 【包含】 k8s-prd-cum-2020-09-05,距离过期还剩12.53天, 跳过删除...
all* 【包含】 k8s-dev-px-mlcloud-core-2020-09-14,距离过期还剩21.53天, 跳过删除...
all* 【包含】 k8s-prd-cum-2020-09-06,距离过期还剩13.53天, 跳过删除...
all* 【包含】 k8s-prd-cum-2020-09-03,距离过期还剩10.53天, 跳过删除...
all* 【包含】 k8s-prd-cum-2020-09-04,距离过期还剩11.53天, 跳过删除...
```




