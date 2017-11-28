## LWGO - golang package manager for carlt

### Install
```bash
go get github.com/carltd/lwgo
go install
$GOPATH/bin/lwgo -h
```
### Usage
#### init project
```bash
cd /path/to/your/project/parent
lwgo init
# flow the steps for create package.json file.(if project name exists, the name must be eq it. )
```

#### add deps

##### use cli
all available command for add deps:

```bash
# address the address of repo(start with git:|http:|https|domain|ip)
# [branch#] option, default is `master`

# 1. fetch any version(last commit)
lwgo add address [branch#]*

# 2. fetch fixed (branch and) tag
lwgo add address [branch#]tag

# 3. fetch main version in tag(tag must match d.d.d)
#   eg. ^1.2.3, match 1.*.*
lwgo add address [branch#]^tag

# 4. fetch sub version in tag
#   eg. ~1.2.3, match 1.2.*
lwgo add address [branch#]~tag
```

##### edit **package.json**
edit `package.json` file(under project's dir)

#### gets all deps
```bash
cd /path/to/your/project
lwgo gets
```
