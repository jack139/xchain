##  应用层API

###  一、 说明

​		应用层API与区块链节点一起部署，提供给客户端调用，进行基础的区块链功能操作。


### 二、 概念和定义

#### 1. 节点

​		节点是区块链上的一个业务处理和存储的单元，是一个具有独立处理区块链业务的服务程序。节点可以是一台物理服务器，也可以是多个节点共用一个物理服务器，通过不同端口提供各自节点的功能。

#### 2. 链用户

​		链用户是具有提交区块链交易权限的用户。每个链用户通过一对密钥识别（例如下例中的PubKey），同时使用此密钥进行数据的加密解密操作，因此链用户的密钥需要妥善保管。密钥类似如下形式：
```json
{
	"sign_key":{
		"type":"ed25519/privkey",
		"value":"UgM13IPx/BkwfQo8jce6TMR5bRuAv7ZLdBooTZWm2ixLaNitCW91NHW06h8VQw=="
	},
	"CryptoPair":{
		"PrivKey":"tgNfUoYkh9xKs1hVKs+5uXNetCxvDRRHBNmLMs5/NKk=",
		"PubKey":"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs="
	}
}
```

#### 3. 交易区块
​		链上数据存储在区块链的区块中，区块目前分两类：（1）交易区块；（2）授权区块。交易区块用于存储交易数据。交易区块中的数据是加密的。链用户只能查看自己提交的区块上的加密数据。如果要查看其他链用户的区块加密数据，需要向区块所有者（即区块的提交者）进行请求授权。当区块所有者同意并授权后，请求方才能看到相应加密区块的数据。同时，请求和授权过程也会记录在区块链上，用于追溯。

**交易区块内容：**

| 名称       | 类型   | 说明                                       |
| ---------- | ------ | ------------------------------------------ |
| ID         | uuid   | 交易ID，自动生成                           |
| UserID | string | 链用户ID（即，链用户公钥）                 |
| Data       | string | 加密交易数据（只有链用户UserID可解密） |

**授权区块内容：**

| 名称           | 类型   | 说明                                   |
| -------------- | ------ | -------------------------------------- |
| ID             | uuid   | 授权ID，自动生成                       |
| UserID     | string | 数据原始提交者的用户ID（链用户公钥） |
| AuthUserID | string | 请求授权的用户ID（链用户公钥）       |
| Data           | string | 加密交易数据（AuthUserID可以解密） |
| Action         | byte   | 交易类型：4 请求授权， 5 响应授权      |



### 三、 API提供的区块链功能

| 序号 | 接口名称        | 接口功能                                     |
| :--: | :-------------- | -------------------------------------------- |
|  1   | deal            | 提交交易数据                                 |
|  2   | auth_request    | 请求授权查询指定交易                         |
|  3   | auth_response   | 授权查看指定交易                             |
|  4   | query_deals     | 查询自己的所以交易                           |
|  5   | query_auths     | 查询授权请求和授权响应                       |
|  6   | query_block     | 按区块ID查询制定区块                         |
|  7   | query_raw_block | 按区块ID查询制定区块原始数据                 |




### 四、接口定义

#### 1. 全局接口定义

输入参数

| 参数      | 类型   | 说明                          | 示例        |
| --------- | ------ | ----------------------------- | ----------- |
| appid | string | 应用渠道编号                  |             |
| version   | string | 版本号                        | 1 |
| sign_type | string | 签名算法，目前使用SHA256算法 | SHA256 |
| sign_data | string | 签名数据，具体算法见下文      |             |
| timestamp | int    | unix时间戳（秒）              |             |
| data      | json   | 接口数据，详见各接口定义      |             |

> 签名/验签算法：
>
> 1. appid和app_secret均从链用户密钥文件中```sign_key.value```字段生成：appid为```sign_key.value```做MD5（字母小写），app_secret既是```sign_key.value```字段。
> 2. 筛选，获取参数键值对，剔除sign_data参数。data参数按key升序排列进行json序列化。
> 3. 排序，按key升序排序；data中json也按key升序排序。
> 4. 拼接，按排序好的顺序拼接请求参数。
>
> ```key1=value1&key2=value2&...&key=appSecret```，key=app_secret固定拼接在参数串末尾。
>
> 4. 签名，使用制定的算法进行加签获取二进制字节，使用 16进制进行编码Hex.encode得到签名串，然后base64编码。
> 5. 验签，对收到的参数按1-4步骤签名，比对得到的签名串与提交的签名串是否一致。

签名示例：

```json
请求参数：
{
    "appid": "66A095861BAE55F8735199DBC45D3E8E", 
    "version": "1", 
    "data": {
        "test1": "test1", 
        "atest2": "test2", 
        "Atest2": "test2"
    }, 
    "timestamp": 1608904438, 
    "sign_type": "SHA256",  
    "sign_data": "..."
}

密钥：
app_secret="43E554621FF7BF4756F8C1ADF17F209C"

待加签串：
appid=66A095861BAE55F8735199DBC45D3E8E&data={"Atest2":"test2","atest2":"test2","test1":"test1"}&sign_type=SHA256&timestamp=1608948188&version=1&key=43E554621FF7BF4756F8C1ADF17F209C

SHA256加签结果：
"fa72d34eafea3639b0a207bdd7ceb49586f4be92e58ee97b6453b696b0edb781"

base64后结果：
"ZmE3MmQzNGVhZmVhMzYzOWIwYTIwN2JkZDdjZWI0OTU4NmY0YmU5MmU1OGVlOTdiNjQ1M2I2OTZiMGVkYjc4MQ=="
```

返回结果

| 参数      | 类型    | 说明                                                         |
| --------- | ------- | ------------------------------------------------------------ |
| code      | int   | 状态代码，0 表示成功，非0 表示出错                                 |
| msg   | string | 成功时返回success；出错时，返回出错信息                                                     |
| data      | json    | 成功时返回结果数据，详见具体接口                |

返回示例

```json
{
    "code": 0, 
    "msg": "success", 
    "data": {
    }
}
```

全局出错代码

| 编码 | 说明                               |
| ---- | ---------------------------------- |
| 9000 | 签名错误                           |



#### 2. 交易提交接口

请求URL

> http://<host>:<port>/api/deal

请求方式

> POST

输入参数（data字段下）

| 参数      | 类型   | 说明               |
| --------- | ------ | ------------------ |
| data      | string | 交易数据           |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 交易区块id                              |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "data": "zzzzz", 
    }, 
    "timestamp": 1609384008, 
    "appid": "dec213b6aced0336932e272f3faaf9e4", 
    "sign_data": "MmRmOWIyOTNmNjgxMzBlMjU4ZDU3YjIwNDE4MDI2ZDJhNjFiMDhhMzUyOGI1MmE0YzY3YjZlODdlZjIxN2RiNQ=="
}
```

返回示例

```json
{
    "code": 0, 
    "data": {
        "id": "5cee67ef-b86c-4056-816b-f19334afa41f"
    }, 
    "msg": "success"
}
```



#### 3. 请求授权接口

请求URL

> http://<host>:<port>/api/auth_request

请求方式

> POST

输入参数（data字段下）

| 参数             | 类型   | 说明                 |
| ---------------- | ------ | -------------------- |
| deal_id          | string | 要查看的交易的区块ID |
| from_user_id | string | 区块提交者的ID       |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 授权区块id                              |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "deal_id": "7adfe653-afc8-45d8-901a-e001fdbb102d", 
        "from_user_id": "j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg="
    }, 
    "timestamp": 1609384238, 
    "appid": "dec213b6aced0336932e272f3faaf9e4", 
    "sign_data": "NzUxZTkwYmIzMjEwMzIzYzY5MmE2YTZmMDdlMTBjODc5Mjg2OTc1OWFkOTdmNDg4OTVlNzVkMTUyMWM1YWU1Ng=="
}
```

返回示例

```json
{
    "code": 0, 
    "data": {
        "id": "5cee67ef-b86c-4056-816b-f19334afa41f"
    }, 
    "msg": "success"
}
```



#### 4. 响应授权接口

请求URL

> http://<host>:<port>/api/auth_response

请求方式

> POST

输入参数（data字段下）

| 参数    | 类型   | 说明             |
| ------- | ------ | ---------------- |
| auth_id | string | 请求授权的区块ID |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 授权区块id                              |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "auth_id": "ef57fd9e-66c8-4d23-b142-8bc32b57bfcd"
    }, 
    "timestamp": 1609384388, 
    "appid": "dec213b6aced0336932e272f3faaf9e4", 
    "sign_data": "ZWUzYmMxMDQ2MjQyOTY4M2JhMmMzNTliZmM3YjQ4ZjYwNzM1YjRhODY0MDIxNTFkNDhlODkyNzFjNTRlYjJjNA=="
}
```

返回示例

```json
{
    "code": 0, 
    "data": {
        "id": "5cee67ef-b86c-4056-816b-f19334afa41f"
    }, 
    "msg": "success"
}
```



#### 5. 查询接口

**（1）查询所有历史交易**

请求URL

> http://<host>:<port>/api/query_deals

请求方式

> POST

输入参数（data字段下）

| 参数 | 类型 | 说明 |
| ---- | ---- | ---- |
| 无   |      |      |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 交易列表                                |

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {}, 
    "timestamp": 1609384428, 
    "appid": "dec213b6aced0336932e272f3faaf9e4", 
    "sign_data": "N2IzZTBjOGE1NzZlMDM4YjY0Zjg2Y2YwN2NlMjc4ZjdjNWQyYjdkYWI4N2UyYWNmMDg1Y2E2M2YzYWYxMGMzNA=="
}
```

返回示例

```json
{
    "code": 0, 
    "data": {
        "deals": [
            {
                "data": "zzzzz", 
                "user_id": "qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=", 
                "id": "59534f7d-db5b-4792-8937-09996638c3d4", 
                "send_time": "2020-12-31T03:06:48.535213018Z", 
                "type": "DEAL"
            }
        ]
    }, 
    "msg": "success"
}
```



**（2）查询所有授权请求和授权响应**

请求URL

> http://<host>:<port>/api/query_auths

请求方式

> POST

输入参数（data字段下）

| 参数 | 类型 | 说明 |
| ---- | ---- | ---- |
| 无   |      |      |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 授权列表                                |

> 返回的结果包括：
>
> 1. 其他链用户请求查看当前链用户交易数据的授权请求；
> 2. 当前链用户请求查看其他链用户交易数据，对方返回的授权响应。
>
> 当前链用户提交的授权请求和授权响应，会被对方链用户检索到，自己检索不到。

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {}, 
    "timestamp": 1609384479, 
    "appid": "dec213b6aced0336932e272f3faaf9e4", 
    "sign_data": "N2VmNTk2YThhMzFlNGU3MzlkYzVlYjhmOTg2ZGQzYWFhYTg1OTNkYTJiNmI1ZWUxMGUzMmIyYTgzMzAxMmY4OA=="
}
```

返回示例

```json
{
    "code": 0, 
    "data": {
        "auths": [
            {
                "action": 5, 
                "auth_user_id": "qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=", 
                "data": "xxx", 
                "user_id": "j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg=", 
                "id": "6037ce9d-fb6a-494c-9c6c-2ddc39762796", 
                "deal_id": "7adfe653-afc8-45d8-901a-e001fdbb102d", 
                "send_time": "2020-12-31T03:16:05.876390503Z", 
                "type": "AUTH"
            }, 
            {
                "action": 4, 
                "auth_user_id": "j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg=", 
                "data": "", 
                "user_id": "qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=", 
                "id": "ef57fd9e-66c8-4d23-b142-8bc32b57bfcd", 
                "deal_id": "59534f7d-db5b-4792-8937-09996638c3d4", 
                "send_time": "2020-12-31T03:12:06.064601793Z", 
                "type": "AUTH"
            }
        ]
    }, 
    "msg": "success"
}
```




**（3）查询指定区块ID的交易内容**

请求URL

> http://<host>:<port>/api/query_block

请求方式

> POST

输入参数（data字段下）

| 参数        | 类型   | 说明                         |
| ----------- | ------ | ---------------------------- |
| user_id | string | 区块提交者的ID（链用户公钥） |
| block_id    | string | 区块ID                       |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 指定区块的交易/授权数据                 |

> 说明：
>
> 按区块ID查询时没有限制链用户范围。

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "user_id": "qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=", 
		"block_id": "59534f7d-db5b-4792-8937-09996638c3d4"
    }, 
    "timestamp": 1609385156, 
    "appid": "dec213b6aced0336932e272f3faaf9e4", 
    "sign_data": "ZDRlZTcyNWJiYjRmOGEzMjJiMjE2ZDY2ZGJiMjQ1MzQwZTgwNTVlZDI5N2NjOThkMTE5YWJlNjJhYmVkYjEwOQ=="
}
```

返回示例

```json
{
    "code": 0, 
    "data": {
        "blcok": {
            "data": "zzzzz", 
            "user_id": "qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=", 
            "id": "59534f7d-db5b-4792-8937-09996638c3d4", 
            "send_time": "2020-12-31T03:06:48.535213018Z", 
            "type": "DEAL"
        }
    }, 
    "msg": "success"
}
```



**（4）查询指定区块ID的原始区块数据**

请求URL

> http://<host>:<port>/api/query_raw_block

请求方式

> POST

输入参数（data字段下）

| 参数        | 类型   | 说明                         |
| ----------- | ------ | ---------------------------- |
| user_id | string | 区块提交者的ID（链用户公钥） |
| block_id    | string | 区块ID                       |

返回结果

| 参数 | 类型   | 说明                                    |
| ---- | ------ | --------------------------------------- |
| code | int    | 状态代码，0 表示成功，非0 表示出错      |
| msg  | string | 成功时返回success；出错时，返回出错信息 |
| data | json   | 指定区块的原始区块数据                  |

> 说明：
>
> 按区块ID查询时没有限制链用户范围。

请求示例

```json
{
    "version": "1", 
    "sign_type": "SHA256", 
    "data": {
        "user_id": "qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=", 
        "block_id": "59534f7d-db5b-4792-8937-09996638c3d4"
    }, 
    "timestamp": 1609385186, 
    "appid": "dec213b6aced0336932e272f3faaf9e4", 
    "sign_data": "Yjg0ZGQzNmRlYjFhNTIwMTFlMTExYzM2NjUyNTlkYzcyOTEwNTljYTUwYmEzMGJlYzUxMTdmMTYwOThhMzQ2NA=="
}
```

返回示例

```json
{
    "code": 0, 
    "data": {
        "blcok": {
            "data": {
                "txs": [
                    "eyJTaWduYXR1cmUiOiJ5aVFNTGx6bXp3TGZmR2NxMitiajY0OUlDVmU3cHlaVXUwYVhmYTN4eTdjYktFMm43UWg0WW5TdTNERVNtRFUrNFRLWVQrdGt1QWVjSStmUE5lQ3VCZz09IiwiU2VuZFRpbWUiOiIyMDIwLTEyLTMxVDAzOjA2OjQ4LjUzNTIxMzAxOFoiLCJTaWduUHViS2V5Ijp7InR5cGUiOiJlZDI1NTE5L3B1YmtleSIsInZhbHVlIjoiYm5nUm9FdzFRTCsyUzNRYUtFMlZwdG9zUzJqWXJRbHZkVFIxdE9vZkZVTT0ifSwiUGF5bG9hZCI6eyJ0eXBlIjoiZGVhbCIsInZhbHVlIjp7IklEIjoiV1ZOUGZkdGJSNUtKTndtWlpqakQxQT09IiwiQXNzZXRzSUQiOiJNVEl6IiwiRXhjaGFuZ2VJRCI6InF5QnNYblZLS2p2Rk54SEJSdWRjM3RDcDh0OHltcUJTRjFHYThxbGZxRnM9IiwiRGF0YSI6IkIveVBSMVFIUkZwaU0rVk9DMVlsVDFtOENsd2VRTk1yQlcwTHUxbzlBWlBGNGlQM2NSZGdzZ1dxL0t1ZSIsIlJlZmVyIjoiZW5wNmVubz0iLCJBY3Rpb24iOjF9fX0="
                ]
            }, 
            "evidence": {
                "evidence": []
            }, 
            "header": {
                "app_hash": "0000000000000000", 
                "chain_id": "test-chain-mu6R5U", 
                "consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F", 
                "data_hash": "BCD0AF5B8CA23DA3949962C459553057B1C58ADFDEF742850E91F5976C9B1EE0", 
                "evidence_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855", 
                "height": "3", 
                "last_block_id": {
                    "hash": "DA6997132D77FB664891184ED1769D423913F9D735096155F546F04FD80C2D07", 
                    "parts": {
                        "hash": "0FE0A68DED7CA1952D8E0C14FF00D21E88835E0061FFE6E854C6BDEBECBAF2E3", 
                       	"total": 1
                    }
                }, 
                "last_commit_hash": "094D0C3196022282C0AE5816B54C03C2832003C631AE709E8C87AE0254FBA7C2", 
                "last_results_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855", 
                "next_validators_hash": "B1CD94B2455ADA3D2A90EF23C42827C7F341597C1507468B54D2FC3C92CDD4FC", 
                "proposer_address": "0B9CE4D58B9ECFEF233551D7EDA6346360D72C03", 
                "time": "2020-12-31T03:06:04.660214065Z", 
                "validators_hash": "B1CD94B2455ADA3D2A90EF23C42827C7F341597C1507468B54D2FC3C92CDD4FC", 
                "version": {
                    "app": "1", 
                    "block": "11"
                }
            }, 
            "last_commit": {
                "block_id": {
                    "hash": "DA6997132D77FB664891184ED1769D423913F9D735096155F546F04FD80C2D07", 
                    "parts": {
                        "hash": "0FE0A68DED7CA1952D8E0C14FF00D21E88835E0061FFE6E854C6BDEBECBAF2E3", 
                        "total": 1
                    }
                }, 
                "height": "2", 
                "round": 0, 
                "signatures": [
                    {
                        "block_id_flag": 2, 
                        "signature": "wBbAoBO+ODWkxZoGsPs1nnMBbOifQl7PZtSXcLHSLHIyIpTfncx1W6W4YoBwW7CMbzjkBvMh2sQlKyvvIG5jBA==", 
                        "timestamp": "2020-12-31T03:06:04.660214065Z", 
                        "validator_address": "0B9CE4D58B9ECFEF233551D7EDA6346360D72C03"
                    }
                ]
            }
        }
    }, 
    "msg": "success"
}
```


