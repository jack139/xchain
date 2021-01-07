## 基于国密算法在区块链上进行数据授权交换



### 1. 概述

​		区块链具有分布式、去中心化、防篡改等特性。由于区块链构建的特性，每个节点都存储有全部的区块数据，理论上用户可以查看到所有的区块数据。当使用区块链对敏感数据进行数据交换或数据共享时，需要对用户查看数据的权限进行验证和授权，用户只能查看或共享被授权的数据。同时，因为区块链特有的放篡改性质，将数据授权和数据交换的过程记录在区块链上，可以实现对用户授权、数据交换、数据共享等操作的留痕和追溯。因此，如何在区块链上构建数据授权和交换的方法是具体应用区块链时必不可少的需求。
​		国密算法是我国自主研发创新的一套数据加密处理系列算法。从SM1-SM4分别实现了对称、非对称、摘要等算法功能。特别适合应用于嵌入式物联网等相关领域，完成身份认证和数据加解密等功能。2010年底，国家密码管理局公布了我国自主研制的“椭圆曲线公钥密码算法”（SM2算法）。为保障重要经济系统密码应用安全，国家密码管理局于2011年发布了《关于做好公钥密码算法升级工作的通知》，要求“自2011年3月1日起，在建和拟建公钥密码基础设施电子认证系统和密钥管理系统应使用国密算法。自2011年7月1日起，投入运行并使用公钥密码的信息系统，应使用SM2算法。”
​		现有典型区块链应用实现（例如以太坊、超级账本Fabric）比较完备，但是过于复杂，例如以太坊基于工作量证明实现共识（挖矿），交易吞吐量受限；Fabric基于PBFT共识算法，但是实现过于复杂，所有Fabric应用均需要基于docker容器（称为“链码”的模块），增加节点运维成本。
​		因此，在本次实现中，使用Tendermint作为区块链共识引擎实现轻量化区块链应用，便于维护。去中心化、分布式存储链上信息；使用LevelDB嵌入式数据库存储链表索引数据；使用国密算法进行密钥交换和数据加密。以授权数据交换和共享为场景，实现了：（1）用户数据链上存储，实现防篡改和可追溯；（2）使用SM2非对称密钥作为用户密钥进行用户识别和授权密钥交换；（3）使用SM4对称加密算法进行数据加密；（4）节点分布式去中心化部署，允许1/3节点故障仍可以正常工作。



### 2. 设计概要



### 3. 数据结构设计

```go
// Transx 事务基类
type Transx struct {
	Signature  []byte //发送方对这个消息的私钥签名
	SendTime   *time.Time
	SignPubKey []byte // 签名公钥
	Payload    IPayload
}
```

```go
// 交易信息
type Deal struct {
	ID     uuid.UUID // 交易ID
	UserID []byte //用户的加密公钥
	Data   []byte // 加密交易数据（例如 ipfs hash）
}
```

```go
// 授权操作
// ToUserID 请求 FromUserID 授权，指定 DealID
// FromUserID 加密数据 Data 后 返回 ToUserID
type Auth struct {
	ID          uuid.UUID
	ReqID       uuid.UUID // 授权请求的ID， action==5时使用
	DealID      uuid.UUID // 相关交易ID
	FromUserID  []byte // 数据所有者的用户ID（用户公钥）
	ToUserID    []byte // 被授权的用户ID（用户公钥）
	Data        []byte //  action==4, rb.pub (33 bytes) + rb.priv (B的私钥加密的)
			           //  action==5, ra.pub (33 butes) + data (协商的密钥加密)
	Action      byte // 0x04 请求授权， 0x05 响应授权
}
```


### 4. 使用示例

#### 4.1 本地多节点测试

编译

```shell
$ make build
```


初始化

```shell
$ build/xchain init --home n1
```


启动节点

```shell
$ build/xchain node --home n1 --consensus.create_empty_blocks=false
```

初始化用户密钥
```shell
$ build/xcli init --home users/user1
$ build/xcli init --home users/user2
```

查询验证节点信息

```shell
$ curl http://localhost:26657/validators
```

查询网络信息

```shell
$ curl http://localhost:26657/net_info
```

#### 4.2 数据上链
数据上链
```shell
$ build/xcli deal --home users/user1 "hello world"
```

查询链上数据
```shell
$ build/xcli queryDeal --home users/user1
```

#### 4.3 授权数据交换
> user1 拥有数据，user2向user1发起请求，获取user1的数据
> user1 公钥： qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=

user2 提出数据授权请求
```shell
$ build/xcli authReq --home users/user2 qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs= eea272cb-74ad-4289-aac4-07f84d3284dc
```

user1 查询授权请求
```shell
$ build/xcli queryAuth --home users/user1
```

user1 对授权请求进行授权
```shell
$ build/xcli authResp --home users/user1 6b292c1d-2963-4308-86cb-99fc41c9cd45
```

user2 查询授权结果
```shell
build/xcli queryAuth --home users/user2
```



### 5. 代码工程

https://gitlab.ylzpay.com/guantao/xchain



### 6. 技术栈	

- 区块链中间件：Tendermint 0.34.0（https://github.com/tendermint/tendermint）

- 数据库：LevelDB 1.20（https://github.com/google/leveldb） 

- 开发语言：Go 1.15.6

