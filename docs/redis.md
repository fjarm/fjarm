# Redis

This document describes how [Redis](https://redis.io/docs/latest/operate/oss_and_stack/management/) is managed and used in Fjarm.

## Redis security

**External network access**

In general, Redis should not be exposed to the internet. Access to Redis should be denied to everybody but trusted
clients.

In practice, this looks like binding Redis to a single interface by adding the following to the `redis.conf` file:
```
bind 127.0.0.1
```

**Limiting commands**

Specific commands like `CONFIG` can also be disabled by adding the following to the `redis.conf` file:
```
rename-command CONFIG ""
```

**Code execution**

Redis doesn't require root privileges to run. So, running it as an unprivileged `redis` user is recommended.

**ACLs**

ACLs allow limiting certain connections to specific commands that can be executed or specific keys that can be accessed.

Redis starts with a "default" user that can be configured to provide a specific subset of functionalities to connections
that aren't authenticated.

The following is an example of a default-configured Redis instance's ACLs:
```
> ACL LIST
1) "user default on nopass ~* &* +@all"
```

Users can be stored in the Redis configuration using either:
1. The `redis.conf` file (good for simple use cases)
2. An external ACL file (better for multiple users in a complex environment)

The format in both files is exactly the same in the format `user <username> ... acl rules ...`:
```redis
user worker +@list +@connection ~jobs:* on >ffa9203c493aa99
```

Using the ACL file requires specifying the configuration directive `aclfile` like:
```redis
aclfile /etc/redis/users.acl
```

Sentinel users minimally need the following permissions:
```redis
ACL SETUSER sentinel-user on >somepassword allchannels +multi +slaveof +ping +exec +subscribe +config|rewrite +role +publish +info +client|setname +client|kill +script|kill
```

And Redis replicas minimally need the following commands to be allowed on the primary instance:
```redis
ACL setuser replica-user on >somepassword +psync +replconf +ping
```

## Restarting Redis without downtime

The `CONFIG SET` command can be used to modify Redis's configuration parameters without restarting the server.
But, not all configuration parameters can be modified using `CONFIG SET`.
For example, the `port` parameter cannot be modified using `CONFIG SET`.

To avoid downtime when restarting Redis, the following steps can be taken:
* Set up the new Redis instance as a replica of the current Redis instance
    * If using a single server, the replica must start on a different port
* After replica synchronization is done, use the `INFO` command to ensure the replica and the primary have the same number of keys
* Allow writes to the replica using `CONFIG SET slave-read-only no`
* Tell clients to use the new instance
    * `CLIENT PAUSE` can be used to pause clients for a short time to allow the switch
* Use `MONITOR` to ensure no writes are happening to the old primary
* Elect the replica to the new primary using `REPLICA OF NO ONE`
* Use `SHUTDOWN` to stop the old primary

## Redis Cluster vs Redis Sentinel

### Redis Sentinel
- Focuses on high availability through monitoring and automatic failover
- Uses primary-replica architecture
- No data partitioning
- Suitable for setups where data fits on a single Redis instance

### Redis Cluster
- Provides both data sharding and high availability
- Automatically partitions data across multiple nodes
- Built-in failover for individual shards
- Handles large datasets through horizontal scaling
- Uses 16384 hash slots for data distribution

## Documentation
* [Redis OSS management](https://redis.io/docs/latest/operate/oss_and_stack/management/)
* [Upgrading or restarting Redis without downtime](https://redis.io/docs/latest/cluster-tutorial#upgrading-or-restarting-redis-without-downtime)
* [Redis security overview](https://redis.io/docs/latest/operate/oss_and_stack/management/security/)
* [Redis ACLs overview](https://redis.io/docs/latest/operate/oss_and_stack/management/security/acl/)
* [Redis TLS overview](https://redis.io/docs/latest/operate/oss_and_stack/management/security/encryption/)
* [Redis Streams intro](https://redis.io/topics/streams-intro)

## Demos and tutorials
* [Redis tutorials](https://redis.io/learn/operate/redis-at-scale)
* [CQRS with Redis Streams](https://redis.io/learn/howtos/solutions/microservices/cqrs)
* [Event driven architecture with Redis Streams - Harness blog](https://www.harness.io/blog/event-driven-architecture-redis-streams)
* [Streams in Redis - LogRocket](https://blog.logrocket.com/why-are-we-getting-streams-in-redis-8c36498aaac5/)
* [Redis Cluster tutorial](https://redis.io/docs/latest/cluster-tutorial)
