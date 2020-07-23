# dependency graph

树形打印`go mod graph`


## 安装

```
go get -v -u github.com/zoroqi/dependency-graph

cd ${YOUR PROJECT PATH}

dependency-graph
```

## 一些新知识

`go mod graph` 可以查看依赖关系

`go list -m all` 查看准确依赖版本

`go list -m -u -json all` 依赖详细信息

`go mod why -m all` 查看依赖路径

`go mod why -m github.com/xxx/xxx` 指定package依赖路径

`incompatible`代表包没有按照golang的规范进行版本管理 [挺好的文档](https://github.com/RainbowMango/GoExpertProgramming)
 
```
github.com/xxx/xxx@v2.0.0 就是不规范的 
github.com/xxx/xxx/v2@v2.0.0 就是合规的
github.com/xxx/xxx.v2@v2.0.0 就是合规的
```

[测试项目snake](https://github.com/1024casts/snake), 代码不多依赖不少.

## 一些问题

1. 存在循环依赖
    ```
    这两个包就是互相依赖, 不清楚作者为啥要这样写. 做兼容吗?
    github.com/ugorji/go/codec
    github.com/ugorji/go
    ```
2. 重复依赖导致输出过长
