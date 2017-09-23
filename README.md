# k-on server

センサの値を記録したり，グラフ表示する為のWebサーバ．

- k-on = 計温
- 値の取得には認証かかっていません
- データストアはLevelDB(ローカルファイル)とElasticsearchが使える気がします
- [Prometheus](https://prometheus.io/) でのmetricsの取得もできます

![chart sample](doc/images/chart01.png)


# Get Started

## Build
``` bash
go get github.com/binzume/k-on-server
cd $GOPATH/src/github.com/binzume/k-on-server
go build
./k-on-server -p 8080
# open http://localhost:8080/
```

### Docker

``` bash
./build-on-docker.sh
docker run -d -v /path_to_data_dir:/data -p 8080:8080 k-on-server
```

## Request samples

``` bash
curl -X POST http://localhost:8080/device --data "name=test&description=hoge&fields=temp,humid&secret=test"
curl -X POST http://localhost:8080/stats/test/values --data "temp=1.0&humid=40.5&_secret=test"
curl -X POST http://localhost:8080/stats/test/values --data "temp=2.0&humid=0&_secret=test"
curl -X POST http://localhost:8080/stats/test/values --data "temp=2.71&humid=3.14&_secret=test&_timestamp=1500222333000"
curl -X GET "http://localhost:8080/stats/test/values?offset=0&limit=10"
curl -X GET "http://localhost:8080/stats/test/values/latest"
curl -X DELETE "http://localhost:8080/stats/test/values/1500222333000?_secret=test"
```

## Register Data


# API

json api.

## GET /status

Retrun "ok".


## POST /device

Register new device or update.


Params:

- name: device name ([a-z0-9_]{1,32})
- description: device description. (text)
- fields: (comma separeted)
- secret: device key.


```
curl -X POST http://localhost:8080/device --data "name=_default&description=test&fields=temp,humid&secret=test"
```


Update secret.

```
curl -X POST http://localhost:8080/device --data "name=_default&description=test&fields=temp,humid&secret=new&_secret=test"
```

## POST /stats/:dev_name/values

Register new value.

Params:

- _secret: device key.
- _timestamp: timestamp in millis (optional).
- device defined fields.

```
curl -X POST http://localhost:8080/stats/_default/values --data "temp=10.0&humid=40.5&_secret=test"
```



## GET /stats/:dev_name/values?from=0&limit=100

Return values.

```
curl http://localhost:8080/stats/_default/values
```


## GET /stats/:dev_name/values/latest

Return latest value.

## DELETE /stats/:dev_name/values/:timestamp

Delete value.

## GET /metrics

Metrics for Prometheus.

# License

MIT license
