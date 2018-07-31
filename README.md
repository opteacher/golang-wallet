# 钱包系统
基于区块链的钱包项目，配备充币、归集和提币三大块服务。并可以使用多种（目前只实现了RPC）API调用，使用前需配置各个币种的钱包节点，并在该项目的配置信息中指明。一个项目实例只能跑一个币种，充提归三个服务可以按自身需求开关。
> 使用前最好具备一定的区块链知识，**如因为操作不当导致资产丢失，本平台不担负任何责任**
## 环境配置
> 如果使用docker配置服务，可不手动做如下配置
* golang1.10.2
* mysql：存储流水和交易
* redis：缓存交易状态
## 安装和配置
下载项目文件

`> git clone git@github.com:opteacher/golang-wallet.git`

配置项目
* `config/settings.json` 全局配置
```json
{
	"env": "dev",
	"services": [
		"withdraw", "deposit"
	]
}
```
| - | - | - |
|---|---|---|
| env | 使用环境 | |
| services | 开启服务 | withdraw（提币）</br>deposit（充币）</br>collect（归集） | |


* `config/coin.json` 币种相关配置
```json
{
	"name": "ETH",
	"url": "http://127.0.0.1:8545",
	"decimal": 8,
	"stable": 2,
	"collect": "0xb78f085e2759baf782c705cd3a9fcb5d39fa7b3c",
	"minCollect": 0.0001,
	"collectInterval": 30,
	"tradePassword": "FROM",
	"unlockDuration": 20000,
	"withdraw": "0x1c6f567e577a351917615fb1c8f1222dc96ba18d"
}
```

- | -
---|---
name | 币种简称
url | 钱包节点
decimal | 币种精度
stable | 转账到账最低稳定块高（防止支链追赶主链）
collect | 归集地址
minCollect | 最小归集金额
collectInterval | 归集间隔
tradePassword | 充值账户的交易密钥
unlockDuration | 解锁充值账户的时间
withdraw | 提币账户/地址

* `config/(dev/prod/..).json` 环境配置，可以自定义名字，并在settings.json指定
```json
{
	"db": {
		"url": "127.0.0.1:3306",
		"name": "test",
		"username": "root",
		"password": "12345",
		"max_conn": 20
	},
	"redis": {
		"password": "12345",
		"time_format": "RFC3339",
		"clusters": [
			{
				"name": "main",
				"url": "127.0.0.1:6379"
			}
		]
	}
}
```

- | - | - | -
---|---|---|---
db | 数据库配置
- | url | 数据库位置
- | name | 数据库名
- | username | 登陆用户名
- | password | 登陆用户密码
- | max_conn | 连接池最大连接数
redis | redis缓存配置
- | password | 查询操作密码
- | time_format | 存储的时间格式
- | clusters | 集群列表
- | - | name | 节点名
- | - | url | 节点URL

> 如果集群列表clusters只有一个节点，则会以单客户端形式调用redis

