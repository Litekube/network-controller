## 功能 & 设计概述

* [功能 &amp; 设计概述](#功能--设计概述)
  * [持久层设计](#持久层设计)
  * [应用层功能增强](#应用层功能增强)
* [后期迭代](#后期迭代)

### 持久层设计

> sqlite

```
sqlite3 /tmp/litevpn.db
```

- sqlite表vpn_mgr结构设计

```sql
create table if not exists "vpn_mgr" (
		"id" integer primary key autoincrement,
		"token" text not null unique,
		"state" integer not null,
		"bind_ip" text not null default "",
		"create_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
    "update_time"    timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
)
```

- update_time触发器
  - 由于sqlite不支持on update关键字，故定义update_time_trigger实现update_time功能

```sql
CREATE TRIGGER if not exists update_time_trigger UPDATE OF id,token,state,bind_ip,create_time ON vpn_mgr
BEGIN
	UPDATE vpn_mgr SET update_time=datetime(CURRENT_TIMESTAMP, 'localtime') WHERE id=OLD.id;
END
```

### 应用层功能增强

- 控制面板基于tcp gRPC+protobuf实现通信交互服务
  - 安全通信，支持tls
- 分离grpc和vpn两套证书

- 获取本机已连接的vpn ip（连接过且未取消注册的）
- 可通过设置摧毁网络，主机自动DHCP到全新的ip，后续重连ip依旧与主机绑定，任意一台服务器，重连以后具有稳定的ip 
  - network manager端生成unique token并和组网ip绑定
    - server端统一认证
    - token的具有时效性，默认10min内
      - goroutine每隔1min检查过期token，并删除
  - 此处的问题在于，需要保留ip，可能存在ip不足的问题，采用LRU策略，根据update_time删除最久没使用的ip
    - ippool相当于cache，和sqlite数据保持同步
      - server启动开启cache同步协程，将sqlite持久化数据映射到cache中
    - 分发新token的vpn ip的处理办法
      - 首先查询sqlite，根据token检查是否已经存在bindIp
        - 有：直接用bindip
        - 无：从cache中找未分配的ip
          - 找到：则分发，并同步cache=1，同步sqlite
          - 没找到（ip不足）：按照LRU策略，从sqlite中找到idle最久的ip，释放，删除旧条目，无需同步cache（ip分给新的token了
- 可注销注册(删除服务器连接，下次建立连接可以生成别的ip)
  - 取消node和组网ip绑定，删除sqlite条目，同步cache=0
- yaml格式的配置文件 client.yml/server.yml，具体字段含义见配置文件注释
  - 网段可任意指定，配置文件中vpnAddr字段
- 可查询网络状态(失联，连接等)
- 为了便于LiteKube交互使用，初始设置一个特殊的node-token，设置值为 "reserverd" (非16位避免被占用），整个网段只会有一台
  node-token="reserverd"总是假定为已经完成了bootstrap，即认为存储有node-token="reserverd"的机器，应该已经具备了client证书
  - vpn server启动时goroutine校验，无此条目则插入

## 后期迭代

> 确定部分

- 多网段

  - 可配置的，通过grpc请求动态增加网段

  - 节点携带node-token和net-token表明节点身份和期待加入的网络
    - 具有不同的net-token的网络互不干扰

> 待验证部分

- vpn-server多副本、高可用，支持client端自动切换vpn-server
  - 涉及到数据迁移、“脑裂”问题，后期工作验证合理性、可行性

- （如果需要）litekube-vpn gRPC CLI 工具
  - 此处可以使用开源grpcurl工具代替