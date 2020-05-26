## GOIM快速入门

### 整体架构(2.0版本)
![2](./2.png)

#### 组件功能划分
+ Zookeeper 支持Kafka的运行
+ Kafka 消息队列启动，业务Push消息到这里
+ Redis 用户在线状态、服务状态的缓存服务启动
+ Discovery 启动服务注册、发现框架
+ logic启动，接收业务Push消息，推送到消息队列
+ comet启动，接受用户注册、连接注册、消息下发
+ job启动，读取消息队列，发送到comet

+ comet / job / logic 支持多实例部署。 同时, push message 消息发布接口从 comet 拆分也有一定的考量, 毕竟多数IM 尤其是 bilibili 的业务场景上来说, 发送量少, 而阅读量多。

### 架构细节(内部逻辑组件与接口关系)

![3](./3.png)




### 问题解决
+ 连接数如何分配 ？
+ 数据如何在多个节点之间转发 ？
+ 节点异常恢复 ？

### 参考资料
+ [《goim 架构与定制》](https://juejin.im/post/5cbb9e68e51d456e51614aab)
+  [《goim架构分析》](https://www.jankl.com/info/goim%20%E6%9E%B6%E6%9E%84%E5%88%86%E6%9E%90)
+ [《goim源码剖析》](https://www.jianshu.com/p/aa8be29397ec)