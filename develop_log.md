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