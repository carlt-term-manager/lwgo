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

```bash
lwgo add src-addr version dst-addr
```

- `src-addr` git地址(ssh或http[s]), 默认库拉取地址，及目录构建依据
- `dst-addr` git地址(ssh或http[s]), 若存在，则从该地址拉取库，否则使用src-addr拉取， 存储路径由src-addr确定
- `version` 可省略， 默认为master分支

下面是所有可用的`命令行`添加依赖方式，详见注释

```bash
# 1. 拉取src-addr库默认分支最新版本
lwgo add src-addr

# 2. 拉取src-addr库dev分支
lwgo add src-addr dev

# 3. 拉取src-addr库release-1.0.0 Tag
lwgo add src-addr release-1.0.0

# 4. 拉取src-addr库bf2120f8a提交
lwgo add src-addr bf2120f8a

# 5. 从备份库拉取，保持路径不变
lwgo add src-addr beta-1.0.1 dst-addr
```

##### 2.2.1 **package.json**添加
直接编辑项目下的package.json文件。位置 deps 小节

```json
{
    "deps":[
    	{
		"src": "https://github.com/xxx/xx.git",
		"dst": "git@git.example.com:another/xxx.git",
		"ver": "release-1.0.0"
        }
    ]
}
```

#### 2.3 安装依赖
```bash
cd /path/to/your/project
lwgo gets # 依赖安装到项目下vendor目录
```
