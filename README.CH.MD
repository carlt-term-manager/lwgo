## LWGO - 卡尔特golang包管理工具

### 一、安装
```bash
go get github.com/carltd/lwgo
go install
$GOPATH/bin/lwgo -h
```
### 二、使用
#### 2.1 初始化项目
```bash
cd /path/to/your/project/parent
lwgo init
# 按提示操作即可创建配置文件（若项目已存在，name中必须填写一致)
```

#### 2.2 添加依赖

##### 2.2.1 命令行添加
下面是所有可用的`命令行`添加依赖方式，详见注释

```bash
# address 为git地址(ssh形式)
# [branch#] 可省略， 默认为master分支

# 1. 匹配address库branch分支任意版本(始终拉取最新版)
lwgo add address [branch#]*

# 2. 匹配address库branch分支固定tag版本
lwgo add address [branch#]tag

# 3. 匹配address库branch分支tag版本中的主版本
#    例如 ^1.2.3, 则匹配1.*.*
lwgo add address [branch#]^tag

# 4. 匹配address库branch分支tag版本中的次版本
#    例如 ~1.2.3, 则匹配1.2.*
lwgo add address [branch#]~tag
```

##### 2.2.1 **package.json**添加
直接编辑项目下的package.json文件。位置 deps 小节

#### 2.3 安装依赖
```bash
cd /path/to/your/project
lwgo gets # 依赖安装到项目下vendor目录,禁止向git库中提交该目录
```
