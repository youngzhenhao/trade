# Login

## 获取用户名和密码

### 用户名

func GetNPublicKey() string 

###  密码

用手机的唯一标识码作为密码。

## 登录流程

### 1. Login 接口

请求参数：

- username：用户名
- password：密码

func Login(username, password string) (string, error)

返回一个token ，需要保存到客户端，每次请求都需要带上这个token。
错误返回空值，错误信息捕捉异常

### 2.refresh
暂时用Login接口代替


# 托管账户：
### 1.开户（开发票）
第一次开发票的时候会在服务器给用户生成一个托管账户

POST 开发票： /custodyAccount/invoice/apply
开发票请求参数：
请求头：
Authorization ： "Bearer" token

请求体：
格式 Json
"amount": int 请求发票金额
"memo": string 发票备注信息

返回值：
"invoice": string 发票号码
"error": string 错误信息

### 2.查询发票

POST 查询发票： /custodyAccount/invoice/queryinvoice
查询发票请求参数：
请求头：
Authorization ： "Bearer" token

请求体：
{
"asset_id":"00"
}

asset_id：string 资产ID,00表示比特币

返回值：
```json
{
    "invoices": [
        {
            "invoice": "lnbcrt1222220n1pn9jlmxpp5f0yjutf4t2z7hsrp826s4q27f964npgjhjfrtva2pzk240weqt2sdq6dysxcmmkv5sx7umnv3shxerpyqcqzzsxqyz5vqsp5hrtme76j03gaaxle3a3tvd83u86va6q6pltcy2fzta7082ju698q9qyyssqmcey7racq2gu03v54j7jujv2fq7ypkqgj74pcvjpv6p9h5r5lfqplqu9c28lv35x4wrvxvw6hdhjjpppnreqk36he2wyfkxyyknmhqgqxxhtrf",
            "asset_id": "00",
            "amount": 122222,
            "status": 0
        }
    ]
}
```

invoice： string 发票号码
asset_id：string 资产ID,预留，00表示比特币
amount：int 发票金额
status： int 发票状态，0表示未支付，1表示已支付，2表示已失效

"error": string 错误信息
### 3.查询余额


POST 查询余额： /custodyAccount/invoice/querybalance
查询余额请求参数：
请求头：
Authorization ： "Bearer" token

请求体：
无

返回值：
"balance": int 账户余额
"error": string 错误信息

### 4.转账


POST 转账  /custodyAccount/invoice/pay
转账请求参数：

请求头：
Authorization ： "Bearer" token

请求体：
```json 
{
    "invoice":"lnbcrt44440n1pn9mq3fpp5mvnjhz8tguz4e5qene5ztnw4w2rlyvcwecqf53lzyxxvl74aedyqdqqcqzzsxqyz5vqsp5a0gxaku4n3klkjk0x62377u9g97az7mgrmqajp2lludwx70hv3eq9qyyssqq63zjpahhy3r0nlqduyjpjrttlfwxvm052qmehupxfpflyg4fj6yrytfwyl546xd3ptupusr8gazky76f30jny8zhtvx4vrf0ngs72cq98qtry",
    "feeLimit": 0
}
```

invoice： string 发票号码
feeLimit：int 转账手续费限制，默认0，单位satoshi，0表示不限制手续费


返回值:

"error": string 错误信息
"payment": 当成功时返回success



### 5.交易记录

POST 查询交易记录： /custodyAccount/invoice/querypayment

转账请求参数：

请求头：
Authorization ： "Bearer" token


请求体：
```json 
{
  "asset_id": "00"
}
```

返回值：
```json
{
  "payments": [
    {
      "timestamp": 1717645790,
      "bill_type": 1,
      "away": 1,
      "invoice": null,
      "amount": 2500,
      "asset_id": "00",
      "state": 1
    },
    {
      "timestamp": 1717645730,
      "bill_type": 0,
      "away": 0,
      "invoice": "lnbcrt1300u1pnxzvv2pp502gytjrc55sj8nql8wy6mh3x0ta32uluq3km407umkq2447ga5csdqld96zwueqvys8getnwssxjmnkda5kxegcqzzsxqyz5vqsp5ltn6r3vl7hzad90dcpyndd8qxzx0ledgl827maacgdmgz5enk4rq9qyyssq6amcxvcz9p5qd9d04j9q2jx3ayut0tdfmwca77kz4drv8kcxc5y4s0hcl8g4p3wcm09wwj2sv575ex5x2ggv7tpwwqens4h623ar82qqamsmuq",
      "amount": 130000,
      "asset_id": "00",
      "state": 1
    },
    {
      "timestamp": 1717644988,
      "bill_type": 1,
      "away": 1,
      "invoice": null,
      "amount": 4500,
      "asset_id": "00",
      "state": 1
    },
    {
      "timestamp": 1717644987,
      "bill_type": 1,
      "away": 1,
      "invoice": null,
      "amount": 3500,
      "asset_id": "00",
      "state": 1
    }
  ]
}
```

payments： 交易记录数组
            timestamp：int 交易时间戳
            bill_type：int 交易类型(预留)
            away：int 收/付(0表示收入，1表示支出)
            invoice：string 交易相关联的发票号码(如果有)
            amount：int 交易金额
            asset_id：string 资产ID,预留，00表示比特币
            state：int 交易状态(0表示挂起中，1表示成功，2表示失败)
error：string 错误信息





