# Gi

- 集合了 gin 的常用middleware、handler、以及对其进行极简封装
- 能直接使用的 middleware 和 handler 便直接暴露出来供调用
- 不能直接使用的，才封装成 [Functional Option](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
- 另提供一个公共函数 gi.With ，用以把任何middleware转成 gi.GinOption
- 这样就可以以统一的方式创建 gin.Engine:

```
router := gi.New(
    gi.WithPprof(),
    gi.With(gi.MidHSTS()),
    ...
)

router.Use(
    gi.MidCORS(),
    gi.MidRecovery(),
    gi.MidLogger(gi.LogWithThreshold(200*time.Millisecond)),
    ...
)

```
