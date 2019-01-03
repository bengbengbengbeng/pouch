# 内部插件逻辑

[插件能力文档](../docs/features/pouch_with_plugin.md)

## 容器插件

### 容器创建插件点：PreCreate

#### 1. env中获取网络4元组信息

```
-e RequestedIP=192.168.5.100 请求容器ip
-e DefaultRoute=192.168.5.1 默认网关
-e DefaultMask=255.255.255.0 掩码
-e DefaultNic=bond0.701 主机网卡
```

根据上述网络的四元组信息来创建容器的network，一般情况下都是创建alinet的网络，详细的创建方法可见代码。

#### 2. env中设置admin uid

```
-e ali_admin_uid=0
```

如果设置了-e ali\_admin\_uid=0，那么会生成一个默认的uid给容器中admin用户使用，设置的规则是：

```plain
500+ip的最后一位数
```

例如： -e RequestedIP=192.168.10.121 -e ali\_admin\_uid=0， 那么最终到容器中uid为500+121=621，可以在容器的/etc/passwd中可见。

#### 3. env中设置的富容器开关

```
-e ali_run_mode=common_vm
或者
-e ali_run_mode=vm
```

#### 4. 容器user变更

如果开启了富容器模式 -e ali\_run\_mode=common\_vm，那么如果容器设置的User非root，那么容器配置的User会被强制替换为root

#### 5. DiskQuota的数据结构的转换

##### 5.1 对-l DiskQuota的转换

内部场景DiskQuota的设置是通过容器的label设置的，详细的设置方法可见：[https://yuque.antfin-inc.com/pouchcontainerblog/documents/dy6utt](https://yuque.antfin-inc.com/pouchcontainerblog/documents/dy6utt)，插件中需要将label内容转换成为pouchd中的数据结构，也就是map[string]string的格式，map的key是表示容器中的volume挂载路径，map的value表示要设置的quota的大小。详细实现见代码。例如：

```plain
例一：
-l DiskQuota=100g
转换之后为：
map[string]string {
    ".*":"100g"
}

例二：
-l DiskQuota=/mnt=10g;.*=100g
转换之后为：
map[string]string {
    "/mnt":"10g"
    ".*":"100g"
}
```

##### 5.2 对label中AutoQuotaId和QuotaId的转换

如果设置了-l AutoQuotaId=true或者-l DiskQuota=10g的格式（不包含分号; 和等号=，也就是不是这种格式 -l DiskQuota=/=10g;/mnt=20g），那么会去判断是否设置了-l QuotaId=1234567，如果设置了QuotaId，那么容器config中QuotaID字段设置为label中的QuotaId，如果没有设置，那么容器config中QuotaID设置为-1。

#### 6. 增加额外的capabilities

插件中为容器增加了额外的capabilities，列表如下：

```plain
"SYS_RESOURCE",
"SYS_MODULE",
"SYS_PTRACE",
"SYS_PACCT",
"NET_ADMIN",
"SYS_ADMIN"
```

#### 7. 在富容器模式下特殊的设置

##### 7.1 将label部分数据转换到env中

将label的中如下3个key值转换到env中

```plain
ali_host_dns
com_alipay_acs_container_server_type
ali_call_scm
```

同时增加额外key的前缀标签：label\_\_
例如：

```
-e ali_run_mode=common_vm
-l ali_host_dns=true
转换后为
-e label__ali_host_dns=true
```

##### 7.2 不对网络hostname相关的文件进行bind

同时富容器情况下，不将容器config中HostnamePath、HostsPath、ResolvConfPath三个配合路径bind到容器中的如下三个文件

```
/etc/hosts /etc/hostname /etc/resolv.conf
```

##### 7.3 ShmSize特殊设置

如果ShmSize没有设置或者设置为0，那么会默认设置为Memory大小的一半。

#### 8. 对hostname的设置

如果容器config中设置了Hostname，同时env中没有设置HOSTNAME，那么会将容器config.Hostname转换到env中的HOSTNAME字段

#### 9. 对volumesFrom字段去除/的前缀

在upgrade的场景下，sigma传递的volumesFrom字段会在开始增加/，例如 volumesFrom=/testcontainer，插件中会去除/，变更为volumesFrom=testcontainer

#### 10. NetPriority的设置

在config.SpecAnnotation["net-priority"]设置对应的值

#### 11. 将label中设置的annotation.的前缀字段设置到SpecAnnotation中

#### 12. 容器cgroup rw属性的设置

如果label中含有pouch.SupportCgroup=true，那么转换成env: pouchSupportCgroup=true
runc中会解析pouchSupportCgroup=true， 将容器cgroup的readonly mount option去掉，最终达到
容器中对cgroup可写的需求。
该label的使用方式为了兼容alidocker上的使用，兼容后，sigma可以直接通过设置label的方式同时
支持alidocker和pouch。alidocker上该功能的代码提交如下：
[label alipay.SupportCgroup=true](http://gitlab.alibaba-inc.com/docker/docker/commit/87bda17027515c0cf60993f9c7d454ecf2ec84cd)

### 容器启动插件点：PreStart

设置runc中prestart hook要执行的工具：

```
/opt/ali-iaas/pouch/bin/prestart_hook
```

## Daemon插件

### 启动插件：PreStartHook

执行daemon启动定制脚本：

```
/opt/ali-iaas/pouch/bin/daemon_prestart.sh
```

检查宿主机相关是否支持pouch的启动运行

加载相关依赖的组件，例如collectd，alinet网络插件，nvidia-docker 插件

### 关闭插件：PreStopHook

执行关闭定制脚本

```
/opt/ali-iaas/pouch/bin/daemon_prestop.sh
```

## CRI插件

### 业务容器创建插件：PreCreateContainer

#### 1. 更新业务容器中网络4元组相关的env

从sandbox 容器中将ip， gateway，mask设置到业务容器的env中，为了给业务容器是富容器情况下，富容器启动是依赖的相关env的信息，分别是：

```
RequestedIP
DefaultMask
DefaultRoute
```

#### 2. 将CRI中env设置的DiskQuota转换到daemon的数据结构

## API插件

### 新增的API

- /host/exec
- /host/exec/result

### Hook

判断是否为sigma2.0的证书，如果是，hook生效

#### /version

version返回 `1.12.6`

#### 对于带有容器名的请求

兼容swarm的做法。

请求来的时候，将容器名前面的`/`去掉；返回时，在容器名前添加`/`
