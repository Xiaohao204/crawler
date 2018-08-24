# crawler
介绍：go语言版分布式爬虫项目


###流程

启动elasticsearch

```
docker run -d -p 9200:9200 elasticsearch
```

开启数据储存服务端

```
go run itemsaver.go --port=1234
```

开启工作线程服务端，可以开多个(需要打开多个终端)

```
go run worker.go --port=9000
go run worker.go --port=9001
```

运行后端

```
go run crawler_distributed/main.go --itemsaver_host=":1234" --worker_hosts=":9000,:9001"
```

运行前端

```
go run crawler_single/frontend/view/starter.go
```

访问http://localhost:8888/
