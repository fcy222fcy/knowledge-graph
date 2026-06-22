Redis是一个基于内存的key-value 键值存储的、可持久化的数据库，并且提供了非常丰富的数据结构，同时还支持非常丰富的功能特性。

redis的特点:
- **性能极高：** Redis 以其极高的性能而著称，能够支持每秒数十万次的读写操作24。这使得Redis成为处理高并发请求的理想选择，尤其是在需要快速响应的场景中，如缓存、会话管理、排行榜等。
- **丰富的数据类型：** Redis 不仅支持基本的键值存储，还提供了丰富的数据类型，包括字符串、列表、集合、哈希表、有序集合等。这些数据类型为开发者提供了灵活的数据操作能力，使得Redis可以适应各种不同的应用场景。
- **原子性操作：** Redis 的所有操作都是原子性的，这意味着操作要么完全执行，要么完全不执行。这种特性对于确保数据的一致性和完整性至关重要，尤其是在高并发环境下处理事务时。
- **持久化：** Redis 支持数据的持久化，可以将内存中的数据保存到磁盘中，以便在系统重启后恢复数据。这为 Redis 提供了数据安全性，确保数据不会因为系统故障而丢失。
- **支持发布/订阅模式：** Redis 内置了发布/订阅模式（Pub/Sub），允许客户端之间通过消息传递进行通信。这使得 Redis 可以作为消息队列和实时数据传输的平台。
- **单线程模型：** 尽管 Redis 是单线程的，但它通过高效的事件驱动模型来处理并发请求，确保了高性能和低延迟。单线程模型也简化了并发控制的复杂性。
- **主从复制：** Redis 支持主从复制，可以通过从节点来备份数据或分担读请求，提高数据的可用性和系统的伸缩性。
- **应用场景广泛：** Redis 被广泛应用于各种场景，包括但不限于缓存系统、会话存储、排行榜、实时分析、地理空间数据索引等。
- **社区支持：** Redis 拥有一个活跃的开发者社区，提供了大量的文档、教程和第三方库，这为开发者提供了强大的支持和丰富的资源。
- **跨平台兼容性：** Redis 可以在多种操作系统上运行，包括 Linux、macOS 和 Windows，这使得它能够在不同的技术栈中灵活部署。



是nosql?数据库,不是以行为单位存储的

基本数据类型
- 字符串 String
- 列表 List
- 集合 Set
- 有序集合 SortedSet
- 哈希 Hash
高级数据类型
- 消息队列 Stream
- 地理空间 Geospatial
- HyperLogLog
- 位图 Bitma
- 位域 Bitfield

命令行界面 CLI (Command Line Interface)
应用程序接口 APl (Application Programming Interface)
图形化界面 GUl (Graphical User Interface)

redis优点
- 性能极高
- 数据类型丰富，单键值对最大支持512M大小的数据
- 简单易用，支持所有主流编程语言
- 支持数据持久化、主从复制、哨兵模式等高可用特性

> 注意
> redis中默认使用字符串来处理数据,而且是二进制安全的---不支持中文(会将中文转换成十六进制形式)
> redis中区分大小写

```bash

# linux上安装redis
yum install redis

# 启动
redis-server
```

```bash

# windows启动
redis-server.exe

# 停止
ctrl + C
```

```bash
# 启动客户端
redis-cli

# 输入ping,回复pong则连接成功

# 退出客户端
quit
```


## Redis通用命令
```bash
# 设置键
set key value

# 获取值--如果不存在返回nil
get key

# 删除键--返回1代表删除
del key

# 判断键是否存在--如果不存在返回0,存在返回1
exists key

# 查看数据库中都有哪些键
keys *
# 查找所有以me结尾的键
keys *me

# 清屏
clear

# 一键删除
flushall

# 查看键的过期时间(time to leave)--返回 -1 表示没有设置过期时间;返回-2 表示过期, 返回其他的数值表示剩余过期时间
# 过期之后再使用get查看该键就没有输出了
TTL key

# 设置键的过期时间
expire key time

# 

# 设置一个带有过期时间的键值对
setex key time value

# 当键不存在设置键的值
setnx key value

# 序列化给定key,并返回序列化的值
dump key
```

```bash
127.0.0.1:6379> set runoob 菜鸟教程
OK
127.0.0.1:6379> get runoob
菜鸟教程
127.0.0.1:6379> exists runoob
1
127.0.0.1:6379> dump runoob  ##序列化给定的key,并返回结果

菜鸟教程6��9�L�
127.0.0.1:6379> expire runoob 30
1
127.0.0.1:6379> ttl runoob
23
127.0.0.1:6379> ttl runoob
17
127.0.0.1:6379> pttl runoob
8688
127.0.0.1:6379> keys *

127.0.0.1:6379> 

```



## 字符串String
```bash



```
## 列表List
```bash
# 头部(左侧)添加元素
lpush key element
# 尾部(右侧)添加元素
rpush key element

# 一次性添加多个元素---从左到右依次添加
lpush letter c d e
# 查看后输出 e d c

# 尾部添加一个字母列表,值是a
rpush letter a


# 获取列表内容
# 起始位置和结束位置都是从0开始数的,stop为-1 表示结尾
lrange key start stop

# 列表中的所有元素
lrange key 0 -1


# 从头删除元素,返回被删除的元素
lpop key
# 从尾删除元素,返回被删除的元素
rpop key

# 指定删除的个数
# 从头删除2个元素
lpop letter 2

# 获取列表长度
llen key

# 组合命令实现先进先出队列
rpoplpush

# 删除列表中指定范围以外的元素--只保留 start 和 stop之间的元素
ltrim key start stop

```

## 集合Set
set是一种无序集合
列表中的元素是可以重复的
集合中的元素是不可以重复的

相关命令是以s开头

```bash

# 创建集合并添加元素
sadd key member

# 查看集合
smemmbers key

# 添加重复元素,输出0 -->添加失败
(integer) 0

# 判断member是否在集合中
sismember key member

# 删除集合中的指定元素
srem key member

# 集合的运算 交集,并集
sinter key key

sunion key [key]



```


## 有序集合SortedSet

有序集合的每个元素都会关联一个浮点类型的分数,然后按照这个分数对集合中的元素按照从小到大的排序

有序集合的元素是唯一的,但是分数是可以重复的
相关指令是 z 有关的


## 哈希表Hash
哈希是一个字符类型的字段和值的映射表,简单来说就是一个键值对的集合,特别适合用来存储对象
相关命令是以h开头
```bash
# 向哈希中添加一个键值对
hset key field value



```


## 发布订阅模式
(消息无法持久化,无法记录历史消息)

publish channel message

subscribe channel

ctrl + C 退出

## 消息队列Stream
是一个轻量级的消息队列,用于解决发布订阅的一些局限性,比如消息无法持久化,无法记录历史消息

相关命令以 x 开头

xadd key 
// 这里的* 表示会自动生成一个消息的ID,回显信息就是生成的id
xadd geekhour * course  redis

// 手工指定id  ,id的格式是一个整数加上一个短横线再加上一个整数
// 第一个整数表示时间戳,第二个整数表示序列号
// 如果使用 * 的话,redis会保证id是自增的;手工指定的话就需要自己来保持自增

xadd geekhour 1-0 course git
xadd geekhout 2-0 course docker

// 查看stream中消息的数量
xlen geekhour

// -和+ 表示查看所有消息
xrange geekhour - +

// 删除消息
xdel

// 删除消息,,,,其中 maxlen 0 表示删除了所有的消息----回显数字表示删除的消息
xtrim geekhour maxlen 0

// 消息创建之后的消费
xread count 2 block 1000 streams geekhour 0
这种可以重复读取

// count 2 一次读两条信息
// block 1000 如果没有消息就阻塞1000毫秒,也就是1秒
// 0 表示从头开始读取;如果是1 表示从第二条消息开始读取;如果超过 最大id ,会阻塞1秒,然后返回nil
// 如果是想获取从现在开始以后的最新的消息,用 $ 

创建消费者组
xgroup create key group id 

geekhour是消费者组的名称, 
xgroup create geekhour group1 0

查看创建的组的信息
xinfo groups key

添加两个消费者
xgroup createconsumer key group consumer

读消费者
xreadgroup goup consumer count 2 block 3000 key >
// group 组名
// consumer 消费者名
// count 2 表示读取两条消息
// block 3000 表示没有消息就阻塞 3000 毫秒,也就是3秒
// key 是消息的名字
// > 表示从消息中读取最新的消息


## 地理空间Geospatial

相关命令以 GEO 开头
```sql
// 添加一个地理位置信息
geoadd key 经度 纬度 城市名字

//key是这个地理位置信息的名字



// 例如
geoadd city 116.405285 39.904989 beijing

// 一次性添加多个物理位置信息---返回几就是添加成功了几个
geoadd city ... ... shanghai ... ... shenzhen ... ... guangzhou ... ... hangzhou

// 获取某个地点的经纬度--返回值第一个表示经度,第二个表示纬度
geopos city beijing

// 获取两个城市间的距离--默认单位是米
geodist city beijing shanghai

geodist city beijing shanghai M|KM|FT|MI
// M 表示米
// KM是千米
// FT 是英尺
// MI 是英里

// 搜索指定范围内的成员并返回
geosearch city frommember shanghai byradius 300 KM

//  shanghai 表示起始范围
// byradius 表示一个圆形的范围
// 300 KM 表示圆形的半径

扩展:
georadius
georadiusbymember
bybox


```

## HyperLogLog
是一个用于做基数统计的算法,不是一个redis特有的算法
**基数**:如果集合中的每个元素都是唯一且不重复的,那么这个集合的基数就是集合中元素的个数

原理:使用随机算法来计算,通过牺牲一部分精确度,来换取更少的内存消耗
优点:内存消耗小
缺点:有误差

适合:精确度不高,而且数据量非常大的统计工作,比如统计某个网站的uv,统计某个词的搜索次数等等

命令都以pf 开头
```mysql

// 添加元素
pfadd key element

// 查看基数
pfcount key

// 合并,结果放在destkey里面
pfmerge destkey key1 key2

```


## 位图Bitmap
位图是字符串类型的扩展,可以使用一个string类型来模拟一个bit数组,数组下标就是偏移量,值只有0和1,支持一些位运算,比如 与 或 非 异或 等等,常用于用户的在线状态,签到情况等等

命令都以bit开头

// 设置偏移量的值
setbit key offset value

// 例如
setbit dianzan 0 1
setbit dianzan 1 0
这样就设置了一个长度为2的位图

// 获取某个偏移量的值
getbit dianzan 1

本质上是一个字符串,我们可以直接使用字符串的命令来设置他的值

一次设置多个值,2进制的 1111 = 16 进制的 0xF  相当于一次性设置了8位 即 11110000
set dianzan "\xF0"

// 统计某一个key的值里面有多少个bit是1
bitcount key start 

// 例如,统计dianzan 这个里面有多少个bit是1
bitcount dianzan 











例如:
集合12345的基数就是5
如果把集合12345的后面再加上12345变成一个10元素的集合,他的基数还是5























redis的中文问题
redis不支持中文,但是可以以十六进制的形式存储
```bash
# 如果想在redis中显示中文
# --raw 表示以原始的形式显示内容
redis-cli --raw


```












