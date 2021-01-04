
## 多节点测试


### 安装LevelDB
```
yum -y install snappy-devel
wget https://github.com/google/leveldb/archive/v1.20.tar.gz && \
  tar -zxvf v1.20.tar.gz && \
  cd leveldb-1.20/ && \
  make && \
  sudo cp -r out-static/lib* out-shared/lib* /usr/local/lib/ && \
  cd include/ && \
  sudo cp -r leveldb /usr/local/include/ && \
  sudo ldconfig && \
  rm -f v1.20.tar.gz
export LD_LIBRARY_PATH=/usr/local/lib
```



### 编译

```shell
make
```



### 初始化

```shell
build/xchain init --home n1
build/xchain init --home n2
```

复制创世块

```shell
cp n1/config/genesis.json n2/config/
```

获取n1节点id

```shell
build/xchain show_node_id --home n1
```

修改n2/config/config.toml

```toml
# 在同一ip下测试时，需要修改端口
proxy_app = "tcp://127.0.0.1:36658"
laddr = "tcp://127.0.0.1:36657"
laddr = "tcp://0.0.0.0:36656"
# 前面的peer id是上面n1节点的id
persistent_peers = "b2c82964b2c67236f94a84aa19b0fda6e91869a0@127.0.0.1:26656"
```

修改各节点的config/config.toml

```toml
create_empty_blocks = false
# 私有网络或单机测试需要设置这个
addr_book_strict = false
```



### 启动节点

```shell
build/xchain node --home n1
build/xchain node --home n2
```



### 新建用户密钥

```
build/xcli init --home users/user1
```



### 启动http服务

```
build/xcli http 8080 users
```



### 提交交易

```shell
build/xcli deal --home users/user1 xxxx
```



### 请求授权

```shell
build/xcli authReq --home users/user1 qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs= eea272cb-74ad-4289-aac4-07f84d3284dc
```



### 响应授权

```shell
build/xcli authResp --home users/user1 6b292c1d-2963-4308-86cb-99fc41c9cd45
```



### 查询交易

```shell
build/xcli queryDeal --home users/user1
build/xcli queryAuth --home users/user1
build/xcli queryTx --home users/user1 qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs= eea272cb-74ad-4289-aac4-07f84d3284dc
build/xcli queryRaw --home users/user1 qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs= eea272cb-74ad-4289-aac4-07f84d3284dc
```



### 查询验证节点信息

```shell
curl localhost:26657/validators
```



### 查询网络信息

```shell
curl localhost:26657/net_info
```
