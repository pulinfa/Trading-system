- [1. 项目介绍](#1-项目介绍)
- [2. 服务器架构](#2-服务器架构)
  - [2.1 基础结构](#21-基础结构)
  - [2.2 连接封装和业务绑定](#22-连接封装和业务绑定)
  - [2.3 自定义业务模块](#23-自定义业务模块)
  - [2.4 全局配置模块](#24-全局配置模块)
  - [2.5 消息封装](#25-消息封装)
  - [2.6 多路由模块](#26-多路由模块)
  - [2.7 读写分析模块](#27-读写分析模块)
  - [2.8 消息队列和协程池](#28-消息队列和协程池)
  - [2.9 连接管理](#29-连接管理)
  - [2.10 连接属性配置](#210-连接属性配置)
- [3. 业务模块设计](#3-业务模块设计)
  - [3.1 项目系统的协议定义](#31-项目系统的协议定义)
  - [3.2 用户模块](#32-用户模块)
  - [3.3 股票模块](#33-股票模块)
  - [3.4 持仓模块](#34-持仓模块)
  - [3.5 用户上线](#35-用户上线)
  - [3.6 拉取用户的持仓记录](#36-拉取用户的持仓记录)
  - [3.7 客户端逻辑](#37-客户端逻辑)
- [4. 项目结构](#4-项目结构)
- [5. 测试](#5-测试)
  - [5.1 功能测试](#51-功能测试)
  - [5.2 性能测试](#52-性能测试)
    - [5.2.1 客户端响应时间](#521-客户端响应时间)
    - [5.2.2 服务端性能分析](#522-服务端性能分析)


## 1. 项目介绍
设计一个后端的服务器，用来做交易的模拟，目标如下：

1. 本地创建两台虚机（centos， red hat等），安装docker环境，一台虚机用于部署server，一台虚机用于部署模拟的10个client。
2. 使用golang，实现client，server的双向通信，包括以下逻辑：

2.1 client与server首次建立tcp连接后，向server端拉取用户持仓列表，并保持长连接；
2.2 server端初始化100个用户持仓列表，并模拟用户持仓列表更新，包含以下场景（10个client需涵盖：持仓从无到有、从有到无、持仓记录条数减少、持仓记录条数增加、持仓记录条数无变化但是持仓数量减少或新增）：

- 新建持仓：客户没有某个市场的持仓时，新买入该市场代码，从而该客户新增一条持仓记录；
- 清空持仓：客户有某个市场代码的持仓时，清空卖出该市场代码（比如持仓数量200股，卖出成交200股），从而删除该客户对应持仓记录；
- 持仓新增更新：客户有某个市场代码的持仓时，买入该市场代码，更新市值与持仓数量（比如持仓100股、最新价10元、市值100 * 10 = 1000元，买入100股，则持股200股，市值200 * 10 = 2000元）；
- 持仓减少更新：客户又某个市场代码的持仓时，卖出该市场代码，更新市值与持仓数量（比如持仓200股、最新价10元、市值200 * 10 = 2000元，卖出100股，则持仓100股、市值100 * 10 = 1000元）

2.3 将server端模拟的成交记录push到client，client展示成交记录消息，支持server端针对client端展示的成交记录消息回退；
3基于实现的业务服务完成功能测试，包括持仓列表拉取展示、持仓列表更新消息推送、成交记录消息推送；
4基于实现的业务服务完成压力测试，重点测试随着socket连接数量的不断增加、信息推送频率的不断增加，server端的性能表现，发现存在的问题加以解决（基于pprof输出函数调用，火焰图等测试效果）
## 2. 服务器架构
### 2.1 基础结构
在server这端，首先要定义基本的方法和属性，如下：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669481577893-c122fbc7-337c-49bd-8167-2071311251d4.jpeg)
在启动服务器的时候，主要的功能包括：

1. 创建addr；
2. 监听端口；
3. 处理客户端的基本业务

在运行服务器的时候，主要就是调用Start方法，调用之后做阻塞处理；
主要的代码结构如下：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669482657077-951c4b32-9c9f-4a41-95a2-1ef2558ed190.jpeg)
### 2.2 连接封装和业务绑定
在之前的设计中，我们在Start()方法中，如果取到一个连接的时候，就直接新建一个协程来处理这个事务，但是这样的有很大的不足之处，我们每次在重新更改业务的时候，需要去Start()中修改对应的代码，耦合度就很高，同时要不能对连接进行管理，所以我们将连接业务代码抽取出来作为单独的一个模块，将连接进行单独的封装，并且同业务进行绑定。
这部分的设计主要设计到如下的内容：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669637315500-d31b3928-b206-4e5d-a24c-857edb2566df.jpeg)
### 2.3 自定义业务模块
基于上述的方面，连接和业务的方法是解耦的，但是目前的方式是通过固定一个方法的方式进行的，这样如果我切换另一个业务代码，就需要去更改server的start里面的代码，这是违背了设计原则的。
所以我们的想法是对于这个业务模块，能够实现自定义的方法注册。
为了能够达到这个目的，将这个这个业务模块封装成为Router模块，但是由于接口是需要所有都实现的，所以先用一个BaseRouter抽象类实现这个接口。当用户需要自定义自己的方法的时候可以通过继承这个BaseRouter类来实现。并且在server中需要增加一个Router的属性，相当于在服务器中注册对应的业务方法，同时，就需要有AddRouter为服务器注册对应的业务模块。
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669646789240-3d8a034f-f340-4057-979d-8f7110c1e7ef.jpeg)
### 2.4 全局配置模块
随着架构的逐步变大，需要的参数会越来越多，同时，由于学校的IP是自动分配的，每周就会定期的修改分配的IP，这导致我们需要频繁的更改参数。所以接下来，我们做了一个加载配置的模块，还有一个全局获取参数的对象。
在实现中，我们使用json数据格式来作为文件的配置文件格式。
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669687955854-82cbf505-0f72-4f58-95b9-6f7ad9a0723d.jpeg)
我们新建一个功能模块，utils/globalobj用于放置全局配置对象
全局配置模块的实现如下：
在代码的实现中，由于服务器的全局配置只能有一份，所以在这里使用的是单例模式：
init()就相当于是静态方法，而var变量就是表示的一个静态的实例
### 2.5 消息封装
不同的业务需要不同的需求，我们之前在实现Request的时候，定义了一个data，但是没有定义一个具体的业务，在这里，我们主要就是对于这个部分的一个完善，将消息进行封装，其主要内容如下：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669737862031-823ac1d8-7fcd-4241-be6c-6e9df5d1ba52.jpeg)
### 2.6 多路由模块
在之前的路由模块中，是一个基础的模块，在这个模块一个服务器只能服务一个业务，但是我们在开发的时候，需要同时使用多个路由来区分不同的一个业务场景，所以在这里设计一个路由器集合，同时需要增加一个消息管理模块，用来支持多路由api调度管理模块。
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669777961594-d65e79b0-88a9-4618-8e68-9914c14b0fc8.jpeg)
### 2.7 读写分析模块
目前，在connection的模块中，我们的读写协程的模块是放在一起进行实现的，这种方式拓展性很差，如果说我想要在写之前或者写之后增加一些处理内容，就很麻烦，耦合性太高。所以考虑将读写进行分离，其结构如下所示：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669794584879-9eb985a3-7460-49d4-9615-e6f9e14ac879.jpeg)
首先在读写的模块直接，需要一个通过，所以就直接在connection的属性里面，增加一个chan的无缓冲通过作为读写之间沟通的桥梁。并更改cennection的创建方法；
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669794822576-d146a376-0cde-4b5e-a32a-c5324ccd47da.jpeg)
### 2.8 消息队列和协程池
在读写分离和多业务模块下，考虑到服务器可能同时会有大量的客户端连接需要提供服务，这个时候一个连接需要对应一个Reader协程，一个Writer协程和一个业务处理的Handle协程，Reader和Writer是阻塞的，不占用CPU资源，但是Handle需要占用大量的资源，特别是在协程之间进行调度的过程中上下文切换的开销会很大，所以后续为了解决这个问题，考虑通过池化层的方式，为处理业务层的代码，从而减少了调度的开销。
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669910530393-ed6f7002-007b-45e3-b6b0-0d22f1261ea0.jpeg)
### 2.9 连接管理
连接最大化的限制，或者连接之前的工作，之后的工作，以及创建连接等。
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669983723134-0ce24be6-1516-482d-a48d-fe7a23c4119a.jpeg)
在集成连接管理模块的时候：首先，需要将连接管理模块加入到server模块中，最开始需要添加一个connMgr的属性，同时进行创建模块的修改。第二，在客户端连接的时候，在新建connection之前，要判断是否达到最大连接数量，最后，在服务器停止的时候，需要对这个数据进行清空。
然后，在建立了连接之后，需要在NewConnection方法中，将conn加入到对应的connMgr中，所以先添加connection中另一个server属性，然后创建的时候进行初始化。
最后在客户端断开连接的时候，conn.stop()中，将对应的conn删除，同时在server中stop中调用清除函数。
在客户端开始连接的时候我们需要做一些事情，比如在开始连接的时候，通知所有人某人已经上线，在断开连接的时候告诉所有人某人已经下线，具体来看一下该部分所在的位置：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669983425520-04860f1e-28b5-439a-8650-b069b6187110.jpeg)
### 2.10 连接属性配置
让用户可以自定义一些自己的属性，来增加系统的一个可扩展性。
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1669985581482-3c1d4fa7-7ba2-4c3e-b305-7615132cb7a9.jpeg)
## 3. 业务模块设计
在搭建好了我们的服务器的架构之后，需要在架构层上面进一步的构造我们的业务，所以在接下来的工作中，我们先大概的划分一下需要做的工作内容包括哪些？
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1670240609718-ea967874-8860-4048-908f-d467b9862675.jpeg)
### 3.1 项目系统的协议定义
我们在前面的框架中，建立了TLV机制来处理消息的信息，接下来，我们具体的来定义我们的系统中，每个消息所代表的业务描述：

| **MsgID** | **Client** | **Server** | **描述** |
| --- | --- | --- | --- |
|  |  |  |  |
|  |  |  |  |

### 3.2 用户模块
对于目前的情况来讲，我们使用多个客户端来对同一个数据库进行读写的操作，所以在这个过程中就需要记录是哪个用户的请求，所以我们在建立用户的连接的时候，就需要将为用户分配一个Uid来唯一的标识。
目前是针对一个客户端就可以修改所有的持仓数据进行的设计。
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1670411959780-925670d2-0ace-4332-981c-e376f1871918.jpeg)
### 3.3 股票模块
对于一只股票来讲，市场和股票的代码就决定了一只股票，就可以确定价格：（市场，股票代码）--->（价格）
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1670414135768-ecffc0ed-08fd-48fe-9544-5e7df6fe443a.jpeg)
在股票的模块中，除了股票的struct之外，还需要一个结构来保持我的所有股票，所以在这里我们采用的方法是使用一个silece来存储所有的股票，并且提供一个init方法：
### 3.4 持仓模块
再持仓的模块中，对于每个市场的每个股票代码使用同一个持仓模块进行：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1670430915325-af974924-6e54-4a36-b707-460fb594ff16.jpeg)
在持仓模块中，对于每个记录，我需要使用一个表来进行记录，所以同样采用一个slice来存储所有的持仓信息。
### 3.5 用户上线
用户在上线的时候，首先需要同步一个Uid来在后续的过程中实现业务的逻辑，所以在客户端连接到服务器的时候，第一步需要做的就是同步两者之间的Uid，为什么不直接使用客户号来作为这个值呢？如果使用客户号来作为Uid的话，客户号需要由客户端进行发起，所以流程就是，（1）客户端发送客户号，（2）服务端通过客户号进行查询；（3）服务端推送持仓记录；而使用随机的Uid的流程是：（1）服务端发送Uid；（2）客户端使用Uid查询；（3）发送持仓记录。
最后我们考虑使用动态分配Uid的方式，考虑的原因主要是，在服务端，我们需要一个对象，这个对象指明了客户端和连接之间的关系，将两者进行绑定，所以这就是我们的User这个模块。在客户端最开始连接的时候就需要创建这样的一个实例。为了能够让方便管理，就需要一个编号来进行不同的User的Id来进行识别，上线的流程如下：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1670744221523-b1ab6400-3d0f-4856-859d-43024825fa3e.jpeg)
### 3.6 拉取用户的持仓记录
客户端在第一次建立了连接之后，需要从服务器端同步自己的持仓记录，但是此时，服务端还不知道这个对应的用户所对应的客户号是多少，所以我们选择让刚开始的时候，客户端主动发送一个请求，来拉取服务器端的持仓记录。设计的流程如下：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1670811756912-a5b16a7d-6f9b-43ac-b77e-833798730e2b.jpeg)
### 3.7 客户端逻辑
在之前的业务中，我们实现了用户在连接服务器初期的工作，包括连接之后同步两者之间的uid，以及拉取用户的持仓记录等等工作。
在本节中，我们要模拟用户的行为，进行自己持仓股票的买入和卖出等操作。
主要的流程如下：
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1670832533386-5bbc8c75-b973-4477-8669-2eb68e07e883.jpeg)
## 4. 项目结构
项目的目录结构如下所示
![](https://cdn.nlark.com/yuque/0/2022/jpeg/34348486/1671212455057-fc261c38-2369-463c-af7a-67f9709c83c3.jpeg)
## 5. 测试
### 5.1 功能测试
功能测试主要就是对于服务器功能完整性的测试，按照项目要求，我们有如下的功能，并通过测试得到了相应的结果：
(1) 用户使用客户号登录，首次建立连接，拉取客户的持仓记录，并且保持长连接，测试结果如下：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671174520319-38779128-9f0d-40ff-b5b7-9c92fe25a59a.png#averageHue=%23262524&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=366&id=u19f60d35&margin=%5Bobject%20Object%5D&name=image.png&originHeight=457&originWidth=648&originalType=binary&ratio=1&rotation=0&showTitle=false&size=31143&status=done&style=none&taskId=ub455c664-a5fa-41e9-898b-08644b22c83&title=&width=518.4)
其中包含的内容为，告诉用户自己登录的账户，以及是用户的持仓记录（只是自己登录的账户的持仓记录），用户的所有持股的市值总和，以及使用系统的指令菜单，最后>>>提示用户输入指令；
（2）为了测试一个用户从无到有的过程，我们先登录一个没有任何持仓记录的用户，然后对该用户执行如下的测试：

- **登录一个新用户**

![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671175124595-6db3ff47-f18b-4ada-b101-37ee0714295f.png#averageHue=%23282625&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=258&id=u3f1583d5&margin=%5Bobject%20Object%5D&name=image.png&originHeight=323&originWidth=639&originalType=binary&ratio=1&rotation=0&showTitle=false&size=23543&status=done&style=none&taskId=u79686e2b-4213-41a5-95fe-4a7cc68f92d&title=&width=511.2)
可以看到除了持仓记录为零，没有别的问题

- **持仓从无到有**

![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671175457598-4b1724e9-6b00-4cff-8224-a4e8bf15f7e1.png#averageHue=%232b2926&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=99&id=u03a3ef4a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=124&originWidth=323&originalType=binary&ratio=1&rotation=0&showTitle=false&size=5836&status=done&style=none&taskId=u838503ba-8d82-443d-baf2-fd4dc2cf587&title=&width=258.4)
买入两只股票，然后查看一下当前持仓情况：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671175503087-dca76307-21d9-4620-804f-75e2e9707134.png#averageHue=%23272524&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=310&id=u9520f511&margin=%5Bobject%20Object%5D&name=image.png&originHeight=388&originWidth=635&originalType=binary&ratio=1&rotation=0&showTitle=false&size=27048&status=done&style=none&taskId=u87d1fce7-8809-41b6-aa15-244d04e345d&title=&width=508)
可以观察到，目前已经成功买入了两只股票

- **持仓条目增加**

这个时候已经有了一个已经有持仓的账户，在这个基础上，我们在购入五只股票，得到的结果如下：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671175896753-96a5ab1f-1b6b-40e2-bf2d-0a0121b735ff.png#averageHue=%23262523&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=398&id=u850bd986&margin=%5Bobject%20Object%5D&name=image.png&originHeight=497&originWidth=637&originalType=binary&ratio=1&rotation=0&showTitle=false&size=34614&status=done&style=none&taskId=ud3fbf856-34cf-41e6-bf6d-f267695190a&title=&width=509.6)

- **持仓条目减少**

持仓条目的减少需要卖掉一整只的股票，所以卖掉两只股票，得到如下，剩余三只的情况：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671176006476-05b577bb-e9c3-4885-a3a4-3db7f534f063.png#averageHue=%23262523&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=344&id=u2059ca2b&margin=%5Bobject%20Object%5D&name=image.png&originHeight=430&originWidth=648&originalType=binary&ratio=1&rotation=0&showTitle=false&size=28577&status=done&style=none&taskId=u6723ca1a-5eed-4ee9-83be-f84bc8b3bb9&title=&width=518.4)

- **持仓数量增加**

分别购买已经持有的股票，则数量就会增加，我们分别选择两只股票进行购买，得到数量增加后的持仓情况如下：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671176097639-e87c5c80-ad09-4d39-b790-4cf9f7174dd8.png#averageHue=%23262524&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=340&id=u3dd6ed2c&margin=%5Bobject%20Object%5D&name=image.png&originHeight=425&originWidth=635&originalType=binary&ratio=1&rotation=0&showTitle=false&size=28752&status=done&style=none&taskId=u09ca84d2-529e-4773-af54-f9deff057da&title=&width=508)
可以看到第二和第三只股票分别增加了我们需求的300和160只，同时市值也跟着改变。

- **持仓数量减少**

同样，和增加数量相反，现在我们需要卖出已有股票中的部分，同样选择两只股票分别卖出100只和200只，得到的结果如下：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671176252345-4476b405-81ed-4955-b111-5cf8bf7a085c.png#averageHue=%23262524&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=341&id=ub583df8e&margin=%5Bobject%20Object%5D&name=image.png&originHeight=426&originWidth=640&originalType=binary&ratio=1&rotation=0&showTitle=false&size=28545&status=done&style=none&taskId=ufec1f540-46b1-43bf-a218-4b2228ec3b4&title=&width=512)

- **查询持仓记录**

持仓记录的查询直接调用holding命令就可以得到，如下所示：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671176315291-ebd7877d-ace6-4fea-95db-f919cdcc2cc4.png#averageHue=%23272624&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=261&id=ufe649fb1&margin=%5Bobject%20Object%5D&name=image.png&originHeight=326&originWidth=631&originalType=binary&ratio=1&rotation=0&showTitle=false&size=20978&status=done&style=none&taskId=ud0f9e311-dc12-4412-9351-6b84b77d598&title=&width=504.8)

- **持仓从有到无**

现在我们将剩余的三只股票全部卖出，让这个账户清空，如下：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671176400607-184f9ce2-b0df-4c80-9754-0a2d8f2e4abb.png#averageHue=%23262524&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=345&id=ud94d19c0&margin=%5Bobject%20Object%5D&name=image.png&originHeight=431&originWidth=632&originalType=binary&ratio=1&rotation=0&showTitle=false&size=28583&status=done&style=none&taskId=u7f0f76e9-406d-436d-9915-4d61318bdeb&title=&width=505.6)
可以看到，这个时候，用户的账户就已经清空了。

- **查询历史交易记录**

历史记录在后台是保存在txt文件中的，在前端中，用户可以通过history命令得到自己的交易记录：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671176483700-1d86d899-0323-4071-9b25-01f01094ef32.png#averageHue=%232e2b28&clientId=u22bc8f8e-8ee0-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=316&id=u1a6588dd&margin=%5Bobject%20Object%5D&name=image.png&originHeight=395&originWidth=560&originalType=binary&ratio=1&rotation=0&showTitle=false&size=32057&status=done&style=none&taskId=u10c44a66-7819-40d3-b3ad-debede62219&title=&width=448)
可以看到这都是我们之前所执行的命令。
### 5.2 性能测试
#### 5.2.1 客户端响应时间
考察了不同的条件下的客户端响应时间，客户端响应时间是值从客户端发送一个请求（卖出或者买入或者查询）到接收到服务端的响应的时间称为客户端响应时间。为了排除实验中网络不稳定的因素，我们进行了多组实验，求每个组的平均客户端响应时间，按照客户端数量和请求间隔时间进行考察。间隔时间是根据区间进行随机来执行的。
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671211770506-a4ef27df-606b-4c48-9f76-2d94f25b9f66.png#averageHue=%23fbfbfa&clientId=ue81e8829-b860-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=595&id=ue97b22b6&margin=%5Bobject%20Object%5D&name=image.png&originHeight=654&originWidth=1035&originalType=binary&ratio=1&rotation=0&showTitle=false&size=35796&status=done&style=none&taskId=u8ae7339d-1838-40c3-9be6-4d42e99eef4&title=&width=940.909070515436)
从测试结果中，可以随着客户端数量的增加，在平均处理时间上总体上呈现出增长的趋势，这是因为随着客户端数量的增加，占用的系统资源增加，所以操作系统需要调度不同的客户端服务，所以处理时间总体上会增长。但是在10个客户端之内并没有明显的增长趋势，这是因为我们的工作池中设置了10个工作协程，这样的话，在10个以内的时间消耗大致一样，因为对于每个工作协程都至多处理一个客户端，而当客户端数量大于工作池协程数量之后，由于客户端需要等待和调度不同的工作协程为客户端服务，所以时间上会有所增加。在请求频率上，随着请求频率的增加，总的来将，平均处理时间有所上升，这是因为请求太密集导致的服务器压力增加。
#### 5.2.2 服务端性能分析
主要是使用pprof，火焰图等方式进行分析。
使用pprof的流程
```shell
//1. 开启服务器
go run main.go -cpuprofile cpu.prof

//2. 开启客户端的压力测试

//3. 使用pprof进行采样收集
go tool pprof http://localhost:8999/debug/pprof/profile

//将生成的文件拿到windows下来，生成网页
go tool pprof -http=":8081" pprof.main.samples.cpu.xxx.pb.gz
```
我们选择的压力测试的数据为，100个客户端，1-2秒发送一个请求，得到的调用图和火焰图如下所示：
![image.png](https://cdn.nlark.com/yuque/0/2022/png/34348486/1671211836333-6935aba7-bebe-4678-97b3-9dda3eeaf29f.png#averageHue=%23f1ceb0&clientId=ue81e8829-b860-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=555&id=u85f08af2&margin=%5Bobject%20Object%5D&name=image.png&originHeight=611&originWidth=1745&originalType=binary&ratio=1&rotation=0&showTitle=false&size=129025&status=done&style=none&taskId=ua3635784-56a5-424f-b1e1-9be3006a025&title=&width=1586.3636019801313)
从火焰图中我们可以发现，在服务器进行服务的过程中，大量的调用Write，这个函数的功能是在对用户提供服务的时候，响应的时候调用Write来写端口达到回传给客户端信息的功能，这个有很大的时间是在做这个工作，为什么会产生这种情况呢？
这是因为我们在读写分离的时候说到，每个客户端来的时候，我们都会创建一个读的协程和一个写的协程来完成相关的工作，在那个时候为了减少工作开销，我们设计了工作池来减少为每个客户端都建立一个协程来进行服务。但是没有对读写的协程进行池化，如果有100个客户端连接，就依然需要有100个读的协程和100个写的协程，所有这主要就是大量调用write的原因。
如何解决这个问题，目前能想到的方式是，同样的方法对写的协程同样做一个工作池，来进行服务，而不是每个用户都需要一个协程进行写，从而减少系统的协程调度开销，以达到系统优化的目的。
从横向上来看，主要的花销可以分为五个部分，分别是调度、读、写、买入、卖出，出了刚刚说到的写的部分的开销比较大，其实调度的部分花销也不小，这主要考虑的因素我想还是和读写协程数量有关系。
然后是在买入和卖出这两个函数上，在这两个模块，主要都是调用的更新的函数，在更新函数的里面，除了业务的逻辑之外，主要的时间是花费在数据库操作上，一样的道理，每个用户都会使用一个协程来处理数据库的更新操作，这就使得调度和数据库操作的时间成为性能上的瓶颈，一种解决的思路还是使用数据库线程池的方式，使用一组特定的协程来对数据库进行操作，这样有利于较少系统调度的操作系统调度的开销。

若有收获，就点个赞吧

 
