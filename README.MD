## LWGO - golang package manager for carlt

### Install
```bash
go get github.com/carltd/lwgo
go install
$GOPATH/bin/lwgo -h
```
### Usage
#### Init Project
```bash
cd /path/to/your/project/parent
lwgo init
# flow the steps for create package.json file.( the name must be eq it if project name exists ).
```

#### Add a Dependency

##### Use Cli
all available command for add deps:

```bash
lwgo add src-addr version dst-addr
```

- `src-addr` required, git address(`ssh` or `http[s]`), default repo fetch address and direction's source.
- `dst-addr` optional, git address(`ssh` or `http[s]`), the address used to fetch package, but import path still use src-addr
- `version` optional, default is the repo's default branch with last commit.

```bash
# 1. fetch last commit on default branch
lwgo add src-addr

# 2. fetch the branch named with **dev**
lwgo add src-addr dev

# 3. fetch the Tag named with **release-1.0.0**
lwgo add src-addr release-1.0.0

# 4. fetch the commit named with **bf2120f8a**
lwgo add src-addr bf2120f8a

# 5. fetch use **dst-addr**, but import path not change
lwgo add src-addr beta-1.0.1 dst-addr
```

##### Edit **package.json** file
edit `package.json` file(under project's dir)

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

#### Updates all deps
```bash
cd /path/to/your/project
lwgo gets
```
