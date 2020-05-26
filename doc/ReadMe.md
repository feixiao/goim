## GOIM快速入门

### 整体架构(2.0版本)
![arch](../docs/arch.png)

### 功能简介（1.x版本）
![arch](./1.jpg)
#### comet
+  长连接服务器，支持长轮训、tcp、websocket连接，有超时机制
+  前端接入可以使用LVS 或者 DNS来转发

#### logic
+ logic 属于无状态的逻辑层，可以随意增加节点，使用nginx upstream来扩展http接口，内部rpc部分，可以使用LVS四层转发。
+ 逻辑处理服务器，消息推送入口，通过userId在router服务器中查找对应的comet服务器，将serverId和消息一起保存到kafka队列
+ 因为**comet需要连接logic**，所以在comet服务器中需要**连接logic的通过LVS的虚IP**，LVS加了real server后不会断开，
所以需要在comet服务触发SIGHUP，重新load配置文件。

#### router
+ 2.0 中router已经不存在，功能合并到哪里了？
+ router 属于**有状态节点**，logic可以使用一致性hash配置节点，增加多个router节点（目前**还不支持动态扩容**），提前预估好在线和压力情况
+ 路由服务器，保存userId和serverId的关系，serverId指的是comet服务器地址。
+ logic**连接router需要一致性hash**，所以不能随意添加router服务器。
+ 不选择redis代替router作者解释是因为有**同一userId多次连接序号分配问题以及原子操作**，我觉得通过key:userId来记录自增，key:userId:seq来记录每个连接，这样也是可以的。



#### job
+ 消息转发服务器，从kafka队列中获取消息，发送给comet服务器，无状态服务器可以随意增删。


### 推送服务具有的特性
+ 1.支持一个端的多个连接，比如同一用户在不同的电脑上登录可以同时接收到消息；
+ 2.不支持客户端的区分（pc/wap/ios/android）
+ 3.目前没有查询某个用户是否在线，某些业务上有显示在线状态的需求，这个可以在logic服务器中添加一个接口；发送消息时如果不在线可能会通过其他推送服务（如信鸽推送），这个可以通过修改push接口的返回值来判断
+ 4.群发只能是在线才能收到，如果不在线就丢弃此消息，当然用户可以下次登录时再通过未读消息接口获取
+ 5.在长连接的时候comet->logic->router，而不是comet->router，这样是为了可以在逻辑服务器中获取其他用户状态、未读消息等，这样确实使得comet简单了，且不会有一致性问题。
+ 6.支持群发功能，所有用户或按照roomId来群发，这个功能太爽了，效率不是提高一点半点。


### 问题解决
+ 连接数如何分配 ？
+ 数据如何在多个节点之间转发 ？
+ 节点异常恢复 ？

### 参考资料
+  [《goim架构分析》](https://www.jankl.com/info/goim%20%E6%9E%B6%E6%9E%84%E5%88%86%E6%9E%90)
+ [《goim源码剖析》](https://www.jianshu.com/p/aa8be29397ec)