# 开发日志

## 2020-06-13

基于`go mod graph`打印出依赖树. 结构树形打印 
```
-github.com/xxx/xxx@
    |-github.com/xxx/xxx@v1.3.3
    |-github.com/xxx/xxx@v0.0.0-20180125190556-5a6b3ba71ee6
    |-github.com/xxx/xxx@v2.5.0
    |-github.com/xxx/xxx@v3.2.0
    |-github.com/xxx/xxx@v0.0.0-20180712184237-d95a45783239
    |-github.com/xxx/xxx@v1.4.7
    |-github.com/xxx/xxx@v1.6.3
    |    |-github.com/xxx/xxx@v0.1.0
    |    |    |-github.com/xxx/xxx@v1.3.0
    |    |    |    |-github.com/xxx/xxx@v1.1.0
    |    |    |    |-github.com/xxx/xxx@v1.0.0
    |    |    |    |-github.com/xxx/xxx@v0.1.0
```

tag:v1.0.0

## 2020-06-27

扩充打印方式, 不能只进行简单树形打印. 可以扩展输出方式. 现在输出方式搞不好会特别长, 根本不知道依赖来源于不方便查看.

采用方案
* 在遍历过程中生成新的树
* 存储父节点, 方便遍历

输出方案

1. 依赖来源输出, 方便`grep`发现依赖从哪引入的. 结合`go mod list -m xxxx`查看问题
    ```
    github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx 
    github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx  
    github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx -> github.com/xxx/xxx  
    ```
2. 提供查找方案, 指定包打印输出相应路径, 反向输出

## 2020-07-07

调整参数关系和输出格式. 搜索可以采用树形打印

## 2020-09-12

尝试打印prometheus的依赖树, 然后打印了60万行文件有71Mb. 太多了根本没法看, 需要找到一个方案减少输出. 暂时想到的方案

1. 只打印固定层数的依赖
2. 部分包排除不打印
    * `golang.org/x`这个依赖的居多也很复杂, 而且版本依赖很混乱

## 2021-04-28

1. 发现[godep](github.com/google/godepq)这么个工具, 可以输出graphviz的格式, 这个效果还不错, 在多了以后就会出现问题.
2. 用`list -m all`进行过滤输出, 只打印可能的依赖. 对`mod graph`产生的庞杂依赖进行过滤, 方便发现找到可能的依赖.

> 我发先[oauth2](golang.org/x/oauth2) v0.0.0-20210402161424-2e8d93401602 默认打印可以打印1.7Gb, 太疯狂了.

