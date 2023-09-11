# Gi



## Gin wrapper

提供一系统预定义Option，便于以 [Functional Option](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) 方式创建 gin.Engine:

```golang

router := gi.New(
    gi.WithRecovery(),
    gi.WithPprof(),
    gi.WithCors(),
    gi.WithHSTS(),
    ...
)

```
