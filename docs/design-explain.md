English | [简体中文](./design-explain.md_CN.md)

## Function & Design Overview

* [Function &amp; Design Overview](#function--design-overview)
  * [Persistence layer design](#persistence-layer-design)
  * [Application layer function enhancement](#application-layer-function-enhancement)
* [Later iterations](#later-iterations)

### Persistence layer design

> sqlite

```
sqlite3 /tmp/litekube-nc.db
```

- sqlite table network_mgr, token_mgr structure design

```sql
create table if not exists "network_mgr" (
		"id" integer primary key autoincrement,
		"token" text not null unique,
		"state" integer not null,
		"bind_ip" text not null default "",
		"create_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
    "update_time"    timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
)

// for bootstrap token
create table if not exists "token_mgr" (
		"id" integer primary key autoincrement,
		"token" text not null unique,
		"expire_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
		"create_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
    "update_time"    timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
)
```

- update_time trigger
  - Since sqlite does not support the 'on update' keyword, define update_time_trigger to implement the update_time function

```sql
CREATE TRIGGER if not exists update_time_trigger UPDATE OF id,token,state,bind_ip,create_time ON network_mgr
BEGIN
	UPDATE network_mgr SET update_time=datetime(CURRENT_TIMESTAMP, 'localtime') WHERE id=OLD.id;
END

CREATE TRIGGER if not exists update_time_trigger2 UPDATE OF id,token,create_time,expire_time ON token_mgr
BEGIN
  UPDATE token_mgr SET update_time=datetime(CURRENT_TIMESTAMP, 'localtime') WHERE id=OLD.id;
END
```

### Application layer function enhancement

- The control panel realizes communication interaction service based on TCP gRPC+ Protobuf protocol
  - Secure communication, support TLS
- Separate grpc and network certificates
- Get the connected network ip of node (connected and not unregistered)
- The network can be destroyed by setting, the host automatically DHCP to a new ip, and the subsequent reconnection ip is still bound to the host, any server has a stable ip after reconnection
  - After the trusted node applies for the bootstrap token, the newly registered node needs to carry the bootstrap token to register
    - The bootstrap token is time-sensitive, the user can specify the expiration time (unit is min), the default is within 10 minutes
      - goroutine checks for expired tokens every 10 minutes and deletes them
  - A newly registered node must carry the Bootstrap token for registration
    - Server-side unified authentication
      - After checking the validity of the bootstrap token, the network manager generates a unique token and binds it to the networking ip
      - Return node-token and grpc+network certificates
  - The problem here is that the IP needs to be reserved, and there may be a problem of insufficient IP. The LRU strategy is adopted to delete the IP that has not been used for the longest time according to update_time.
    - ippool is equivalent to cache, keeping synchronization with sqlite data
      - The server starts the cache synchronization goroutine and maps the sqlite persistent data to the cache
    - How to distribute the network ip of the new token?
      - First query sqlite, check whether bindIp already exists according to the token
        - Yes: use bindip directly
        - No: find unassigned ip from cache
          - Found: distribute, synchronize cache=1, synchronize sqlite
          - Not found (insufficient ip): According to the LRU strategy, find the longest idle ip from sqlite, release and delete the old entry, no need to synchronize the cache (the ip is assigned to a new token)
- Unregister (close server connection, next time establish connection can generate other IP)
  - Cancel the node and networking ip binding, delete the sqlite entry, synchronize cache=0
- client.yml/server.yml configuration file in YAML format. For the specific field meanings, see the configuration file notes
  -  The network segment can be arbitrarily specified: networkAddr field in the configuration file
- Network status can be queried (disconnected, connected, etc.)
- In order to facilitate the interactive use of LiteKube, a special node-token is initially set to "reserverd" (non-16-bit to avoid being occupied), and there will only be one node-token in the entire network segment. node-token="reserverd" always assumes that bootstrap has been completed, that is, the machine that stores node-token="reserverd" should already have a client certificate
  - When the network server starts, a goroutine check is started, and if there is no such entry, it is inserted
- log rotate
  - The logs are written to the file in days, and the logs are saved only in the last seven days
  - Define the main log file to get the latest logs
- Network-controller gRPC CLI tool ncadm development

## Later iterations

> Definite part

- Multi-subnet
  - Pre-configurable + dynamically increase network segment via grpc request

  - The node carries the node-token and net-token to indicate the identity of the node and the network it expects to join
    - Networks with different net-tokens do not interfere with each other

> Part to be verified

- network-server has multiple replicas, high availability, and supports automatic switching of network-server on the client side
  - Involving data migration and "split brain" issues, the later work will verify the rationality and feasibility