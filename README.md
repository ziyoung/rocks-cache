## rocks-cache

根据 Redis 协议规范，先解析 Redis 各个命令。待实现的命令如下：

### string（字符串）

> 过期机制暂先不做

- [ ] get
- [ ] set
- [ ] del
- [ ] mset
- [ ] mget
- [ ] exists

### list（列表）

- [ ] rpush
- [ ] lpop
- [ ] rpop
- [ ] lindex
- [ ] lrange
- [ ] ltrim

### hash（字典）

- [ ] hset
- [ ] hgetall
- [ ] hlen

### set（集合）

- [ ] sadd
- [ ] smembers
- [ ] scard
- [ ] spop
