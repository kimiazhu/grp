GRP (Go Reverse Proxy) 是一个用Golang实现的反向代理服务器。

### 现有方案

目前已有的相对简单的反向代理网站的解决方案包括使用nginx的方案,还有一个基于PHP的g7host。

前者非常灵活,但是需要root权限,网上已有专门针对google的Nginx模块,但需要自行编译, 如果要自己配置Nginx其实也比较复杂。

后者简单,并且有简易的配置界面,基于PHP的方案也可以找到一些免费主机来支持, 但是我尝试了一下, 发现这个方案至少访问Google是有问题的, 并没有办法正常使用。

### GRP简介

GRP 希望让部署和配置简单化, 启动可以无需root权限, 并且能够在一台服务器上代理多个网站。

需要同时代理多个网站的时候, 提供多个不同的子域名, 可以在DNS服务器中指定这些子域名执行您部署GRP的服务器。

如果你没有自己的DNS服务器, 或者无法指定多个子域名,也可以通过顶级域名或者IP地址访问,但由于顶级域名和IP地址只有一个,所以这时候只能支持代理单个网站。

> 现在网上其实可以申请到免费的空间以及免费的二级域名解析服务。大家可以自行搜索。所以这其实都不是问题。;)
> 当然有条件的, 使用自己购买的主机和域名, 质量还是会更有保证。

配置示例文件在[conf.yml](https://github.com/kimiazhu/grp/blob/master/conf.yml)中。

### 使用

分发版本包括三个文件:

* grp -- 可执行文件, 包括Linux版, Mac版和Windows版;
* conf.yml -- 配置文件, 本地要使用哪个端口, 你通过哪个子域名代理哪个站点, 都在这里进行配置;
* log4go.xml -- 日志配置文件, 控制日志如何输出, 这个你可以选择忽略。

将以上三个文件上传到你的服务器中, 直接执行即可, 例如:

```bash
$ nohup ./grp > /dev/null 2>stderr.log &
```

### 支持网站列表

已测试网站列表:

- Google搜索
- 百度搜索
- Wikipedia(中文页面)

### TODO

- 支持一个本地子域名对应多个远程站点，无需DNS配置多个子域名的情况代理多个远程地址。(P0)
- SSL支持，当本地服务不是起在80和443端口时，在做域名替换的时候要根据远程服务器的Schema或者端口号智能切换本地端口号。这个问题目前在登录accounts.google.com时会出现。(P1)
- 登录和会话保持。(P1)
- on-file-change auto reload
- multi backend services and load balance support
- 自动感知网站是否使用https? 避免手工在配置文件中进行配置。
- 支持通过顶级域名代理网站同一顶级域名下对应的所有子域名。

### Change Log

##### 2016/08/10

- Initial Commit
- 支持Google和百度搜索, zh.wikipedia.org

##### 2016/08/14

- grp.sh script for Linux

##### 2016/08/18

- refactor && bug fixes
- handle compressed response, replaced all target host name.