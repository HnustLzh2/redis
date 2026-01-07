# REDIS 复现

## 实现了以下数据结构支持
1. key-value
2. list
3. stream
4. zset(sorted Set) 使用跳表
5. channel(订阅通道)

## 特性
1. 支持并发连接
2. 具有身份验证功能
3. 解析RESP2协议成可读命令行
4. 将结果转成RESP2协议的响应
5. 支持从RDB文件中读取历史记录(解析RDB文件)
6. 支持地理位置的解析 使用zset存储地点和经纬度
7. 支持主从节点 从节点提供数据冗余功能
8. 支持事务处理，提供事务队列

### 启动方式
1. pwd: /app  go run main.go  --port    // 服务端监听的地址
--replicaof // 使用从节点启动方式 填写主节点的地址
--dir   // rdb文件存储的目录
--dbfilename    // rbd文件名称

2. pwd: /app/client go run main.go client 运行client 输入命令并获取响应 注意port设置正确