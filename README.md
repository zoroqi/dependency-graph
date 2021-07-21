# dependency graph

树形打印`go mod graph`


## 安装

```
go get -v -u github.com/zoroqi/dependency-graph

cd ${YOUR PROJECT PATH}

dependency-graph
```

## 使用

```
  -ex value
    	exclude package, 排除部分包, 可以多个 -ex xxx -ex yyy
  -expre value
    	exclude package, prefix match, 排除部分包前缀匹配, 可以多个 -expre xxx -expre yyy
  -l int
    	max level 最大打印深度
  -list
    	filter the package in the 'list -m all' result, 基于 `list -m all` 进行过滤
  -p string
    	default: tree print
    	rl: reverse line print
    	rt: reverse tree print
    	wt: whole tree print
    	dot: graphviz print, `xxx | dot -Tsvg -o test.svg` 

  -s string
    	search pkg name, 只打印固定包
```

打印
* default
```
 a
 |-b
   |-c
 |-d
   |-e
```
* rl
```
c -> b -> a
e -> d -> a
```
* rt
```
 c
 |-b
   |-a
 e
 |-d
   |-a
```
* wt
```
 a
 |-b
   |-c
 a
 |-d
   |-e
```
* dot
```
digraph godeps {
0 [label="github.com/oliver006/redis_exporter@" style="filled"]
1 [label="cloud.google.com/go@v0.34.0" style="filled"]
0 -> 1
}
```

## 一些新知识

`go mod graph` 可以查看依赖关系

`go list -m all` 查看准确依赖版本, 并不完全准确.

`go list -m -u -json all` 依赖详细信息

`go mod why -m all` 查看依赖路径

`go mod why -m github.com/xxx/xxx` 指定package依赖路径

`incompatible`代表包没有按照golang的规范进行版本管理 [挺好的文档](https://github.com/RainbowMango/GoExpertProgramming)
 
```
github.com/xxx/xxx@v2.0.0 就是不规范的 
github.com/xxx/xxx/v2@v2.0.0 就是合规的
github.com/xxx/xxx.v2@v2.0.0 就是合规的
```

[测试项目snake](https://github.com/1024casts/snake), 代码不多依赖不少. 特别声明, 尽量设置打印层数.

