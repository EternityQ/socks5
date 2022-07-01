### 1. 简介



那些有效地将组织内部的网络结构与外部网络（如互联网）隔离的网络防火墙和系统的使用正变得越来越流行。 这些防火墙系统通常充当网络之间的应用层网关，通常提供受控的 TELNET、FTP 和 SMTP 访问。随着旨在促进全球信息发现的更复杂的应用层协议的出现，需要为这些协议提供一个通用框架，以透明和安全地穿越防火墙。

此外，还需要以实际可行的、细致的方式对这种穿透进行强认证。 这一要求源于这样一种认识，即客户端-服务器关系出现在各种组织的网络之间，并且这种关系需要受到控制并经常被严格验证。

此处描述的协议旨在为使用TCP和UDP的客户端-服务器应用程序提供一个框架，以便方便且安全地使用网络防火墙的服务。 该协议在概念上是应用层和传输层之间的“填充层”，因此不提供网络层网关服务，例如转发 ICMP 消息。

### 2. 已有实践

当前存在一个SOCKSv4协议，它为基于TCP的客户端-服务器应用程序提供不安全的防火墙穿透，包括TELNET、FTP 和流行的信息发现协议，如 HTTP、WAIS 和 GOPHER。

新协议扩展了SOCKS4模型，增加了对UDP协议的支持、扩展了框架以增加对通用强身份验证方案的规定，并在地址解析方法增加了对域名和 IPv6地址的使用。

为实现SOCKS协议，通常要重新编译或重新链接基于TCP的客户端应用程序，以便使用SOCKS库中对应的封装接口。

注意：

除非另有说明，本文数据包格式图中出现的十进制数表示相应字段的长度，以字节为单位。 当指定的字节必须采用特定值时，语法X'hh'用于表示该字段中单个字节的值。 当使用“变量（Variable）”一词时，它表示相应字段为可变长度，具体长度由其关联的长度字段（一个或两个字节）或数据类型字段定义。

### 3.  基于TCP的客户端的处理过程

当TCP客户端希望建立一个连接而目标为只能通过防火墙访问的对象时（这种确定由实现决定），它必须创建一个连接到SOCKS服务器相应SOCKS端口的TCP连接。 SOCKS服务通常位于TCP端口1080上。如果连接请求成功，客户端将进入协商阶段以便选择要使用的身份验证方法，而后使用所选方法进行身份验证，然后发送中继请求。 SOCKS服务器会评估这些请求，然后建立适当的连接或拒绝它。

除非另有说明，本文数据包格式图中出现的十进制数表示相应字段的长度，以字节为单位。 当指定的字节必须采用特定值时，语法X'hh'用于表示该字段中单个字节的值。 当使用“变量（Variable）”一词时，它表示相应字段为可变长度，具体长度由其关联的长度字段（一个或两个字节）或数据类型字段定义。

客户端连接到服务器，并发送版本标识符/方法选择消息（version identifier/method selection message）：

| VER  | NMETHODS | METHODS  |
| ---- | -------- | -------- |
| 1    | 1        | 1 to 255 |

对于此版本的协议，VER字段被固定设置为 X'05'

NMETHODS字段表示在METHODS字段中的方法标识符（以字节为单位）的数量。

服务器从 METHODS 中给出的方法之一中进行选择，并发送METHOD选择消息（METHOD selection message）：

2

| VER  | METHOD |
| ---- | ------ |
| 1    | 1      |

如果选择的 METHOD 是 X'FF'，则表示客户端列出的方法（服务器）都不接受，而且客户端必须关闭连接。

当前METHOD定义的值是：

- X'00' 无需认证（NO AUTHENTICATION REQUIRED）
- X'01' GSSAPI
- X'02' 用户名/密码（USERNAME/PASSWORD）
- X'03' 至 X'7F' IANA 分配（IANA ASSIGNED）
- X'80' 到 X'FE' 保留用于私人方法（RESERVED FOR PRIVATE METHODS）
- X'FF' 没有可接受的方法（NO ACCEPTABLE METHODS）

然后客户端和服务器进入一个特定于方法的子协商步骤。 与方法相关的子协商的有关描述放在单独的备忘录中。

开发此协议新METHOD支持的人员应联系IANA获取（未分配的）METHOD编号。 有关方法编号及其相应协议的当前列表，应参考分配编号文件（ASSIGNED NUMBERS document）。

一个具有兼容性的实现必须支持GSSAPI并且应该支持USERNAME/PASSWORD身份验证方法。

### 4. 请求

一旦完成依赖于方法的子协商，客户端就会发送请求详细信息。 如果在协商的方法中包括了出于完整性检查和/或机密性目的的封装，则这些请求必须被封装在和方法相关的封装中。

SOCKS请求格式如下：

| VER  | CMD  | RSV   | ATYP | DST.ADDR | DST.PORT |
| ---- | ---- | ----- | ---- | -------- | -------- |
| 1    | 1    | X'00' | 1    | Variable | 2        |

其中：

- VER 协议版本号: X'05'
- CMD
  - CONNECT X'01'
  - BIND X'02'
  - UDP ASSOCIATE X'03'
- RSV 保留未用（RESERVED）
- ATYP 后续地址（即DST.ADDR）的地址类型（address type of following address）
  - IP V4 address: X'01'
  - DOMAINNAME: X'03'
  - IP V6 address: X'04'
- DST.ADDR 期望的目标地址
- DST.PORT 期望的目标端口，网络字节序

SOCKS服务器通常会根据源地址和目标地址评估请求，并根据请求类型返回一条或多条回复消息。

### 5. 地址解析

在地址字段（DST.ADDR、BND.ADDR）中，ATYP 字段指定字段中包含的地址类型：

- X'01' 该地址是IPv4地址，长度为4个字节。
- X'03' 地址字段包含一个完全限定域名(Full qulified domain name)。 地址字段的第一个字节为其后域名的长度（以字节为单位），长度不包括终止符NULL。
- X'04' 该地址是IPv6地址，长度为16个字节。

(关于完全限定域名Full qulified domain name，FQDN总是以主机名开始并且以顶级域名结束，“.”是指根域名服务器，当给出的名字是baidu而不是baidu.的时候，他通常是只主机名，而名字后面带.的通常认为是全域名，eg：www.baidu.com ，其中www是主机名，baidu.是二级域，.com是顶级域。
### 6. 回复

客户端在与 SOCKS 服务器建立连接并完成身份验证协商后立即发送 SOCKS 请求信息。 服务器评估请求，并返回如下形式的回复：



| VER  | REP  | RSV   | ATYP | BND.ADDR | BND.PORT |
| ---- | ---- | ----- | ---- | -------- | -------- |
| 1    | 1    | X'00' | 1    | Variable | 2        |

其中：

- VER 协议版本：X'05'
- REP 回复字段：
  - X'00' 成功
  - X'01' 普通SOCKS服务器故障
  - X'02' 连接不被规则集允许
  - X'03' 网络不可达
  - X'04' 主机不可达
  - X'05' 连接被拒绝
  - X'06' TTL 过期
  - X'07' 不支持的命令
  - X'08' 不支持的地址类型
  - X'09' 到 X'FF' 未分配
- RSV 保留未用，总是为0
- 后续地址(即BND.ADDR)的ATYP地址类型
  - IP V4 地址：X'01'
  - 域名：X'03'
  - IP V6 地址：X'04'
- BND.ADDR 服务器绑定的地址
- BND.PORT 服务器绑定的端口（网络字节序）

标记为保留（RESERVED (RSV)） 的字段必须设置为 X'00'。

如果所选方法因为身份验证、完整性和/或机密性的目的需要对数据包进行封装，则回复将封装在和此方法有关的封装数据包中。

CONNECT



在对CONNECT的回复中，BND.PORT包含服务器分配用于连接目标主机的端口号，而 BND.ADDR包含相关联的IP地址。 提供的BND.ADDR通常与客户端用于访问SOCKS服务器的IP地址不同，因为此类服务器通常是多宿主的（多网卡/IP？）。 在期望（的实现）中，SOCKS服务器将使用DST.ADDR和DST.PORT，以及客户端源地址和端口来评估CONNECT请求。



BIND

BIND请求用于那些要求客户端接受来自服务器的连接的协议（译注：即server-to-client的连接）。 FTP 是一个众所周知的例子，它主要使用“客户端到服务器”的连接来执行命令和获取状态报告，但也可以使用“服务器到客户端”的连接来按需传输数据（例如 LS、GET、PUT）。

Socks5协议期望的情况是，在客户使用CONNECT建立一个主要的（“客户端-服务器”的）连接后，BIND仅用于去建立一些次要的（"服务器-客户端"的）连接。在期望（的实现）中，SOCKS服务器将使用DST.ADDR和DST.PORT来评估BIND请求。

在一个BIND操作期间，两个回复将从SOCKS服务器发送向客户端。第一个在服务器创建并绑定新套接字后发送。 BND.PORT字段的值是 SOCKS服务器分配的用于监听传入连接的端口号。 BND.ADDR字段是其关联的IP地址。 客户端通常会使用这些信息来通知（通过主连接或控制连接）集合点地址的应用服务器（The client will typically use these pieces of information to notify (via the primary or control connection) the application server of the rendezvous address.）。 第二个回复仅在预期的进入连接（incoming connection）成功或失败后发生。



在第二个回复中，BND.PORT和BND.ADDR 字段包含正在连接的主机的地址和端口号。



UDP ASSOCIATE



UDP ASSOCIATE 请求用于在 UDP转发过程中建立偶联以处理UDP数据报(注：“偶联”原文是association，因为UDP是无连接的，所以这里翻译为“偶联”)。 DST.ADDR和DST.PORT字段包含客户端所希望的在偶联上发送UDP数据报时的地址和端口。 服务器可以根据这些信息来限制对偶联的访问。 如果客户端在UDP ASSOCIATE时没有该信息（译注：是不做限制的意思？），则客户端必须使用全为零的端口号和地址。（If the client is not in possesion of the information at the time of the UDP ASSOCIATE, the client MUST use a port number and address of all zeros.）

当发送UDP ASSOCIATE请求的TCP连接被终止时，UDP偶联也结束。

在对UDP ASSOCIATE请求的回复中，BND.PORT和BND.ADDR字段指示客户端转发UDP消息时的目标端口号/地址。



Reply Processing



当回复一个表示失败的值时（REP值不是X'00'），SOCKS服务器必须在发送回复后立即终止TCP连接。 从检测到导致故障的情况起（到断开连接），此时间不得超过10秒。

如果回复代码（X'00' 的 REP 值）表示成功，且请求是BIND或 CONNECT，则客户端现在可以开始传递数据。 如果所选的身份验证方法出于完整性、身份验证和/或机密性的目的封装了数据，则需要使用依赖于方法的封装算法来封装数据。 类似地，当数据到达了此客户端对应的SOCKS服务器时，服务器必须根据所使用的身份验证方法使用适当的算法解封数据。

### 7. 基于UDP的客户端的处理过程



基于UDP的客户端在将其数据报发送到 UDP 中继服务器时，必须使用UDP ASSOCIATE请求响应中由BND.PORT指示的UDP端口。 如果选定的认证方法为了真实性、完整性和/或机密性的目的提供封装，则必须使用适当的封装算法来封装数据报。 每个UDP数据报都携带一个UDP请求头：

| RSV  | FRAG | ATYP | DST.ADDR | DST.PORT | DATA     |
| ---- | ---- | ---- | -------- | -------- | -------- |
| 2    | 1    | 1    | Variable | 2        | Variable |

UDP请求头中各字段的值:

- RSV 保留未用 X'0000'
- FRAG 当前分片号(Current fragment number)
- ATYP 地址类型（address type of following addresses）
  - IP V4 address: X'01'
  - DOMAINNAME: X'03'
  - IP V6 address: X'04'
- DST.ADDR 期望目标地址（desired destination address）
- DST.PORT 期望目标端口（desired destination port）
- DATA 用户数据（user data）



当UDP转发服务器决定转发UDP数据报时，它会静默进行而不会通知请求客户端。

同样它会丢弃它不能或不会中继的数据报。 当 UDP中继服务器收到来自远程主机的回复数据报时，它必须使用上述UDP请求头和依赖于身份验证方法的封装算法来封装该数据报。

UDP转发服务器必须从SOCKS服务器获取期望的客户端IP地址（译注：通常UDP转发服务器和SOCKS 服务器是同一个程序），将向该地址的特定端口发送数据报，端口号为UDP ASSOCIATE的回复中给出的BND.PORT。 它必须丢弃除了来自特定偶联记录的IP地址之外的其他任何源IP地址的数据报（译注：即除了那个发送UDP ASSOCIATE请求的客户端，其他来源都不响应，因为在UDP是无连接的，所以服务器上可能会收到五花八门的数据包）。

FRAG字段指示该数据报是否是多个片段之一。 如果实现了，高位表示片段序列的结束，而值 X'00' 表示该数据报是独立的。 1到127之间的值表示片段序列中的片段位置。 每个接收者将具有与这些片段相关的并包队列（REASSEMBLY QUEUE）和并包计时器（REASSEMBLY TIMER）。 每当REASSEMBLY TIMER超时，或者新数据报到达且该数据报的值小于为该片段序列处理的最高FRAG值时（译注：即乱序了，如前面收到了FRAG值等于100的包后续又收到了FRAG值为99的包），必须重新初始化并包队列并放弃相关的片段。 重组计时器必须不少于5秒。 建议应用程序尽可能避免碎片化。

分片的实现是可选的； 不支持分片的实现必须丢弃任何 FRAG 字段不是 X'00' 的数据报。

一个SOCKS的UDP编程界面(The programming interface for a SOCKS-aware UDP)必须报告当前可用UDP数据报缓存空间小于操作系统提供的实际空间。

- if ATYP is X'01' - 10+method_dependent octets smaller
- if ATYP is X'03' - 262+method_dependent octets smaller
- if ATYP is X'04' - 20+method_dependent octets smaller



### 8. 安全性考虑
本文档描述了IP网络防火墙的应用层穿越协议。 这种遍历的安全性高度依赖于特定实现中提供的特定身份验证和封装方法，并在 SOCKS 客户端和 SOCKS 服务器之间协商期间选择。