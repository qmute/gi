# Gi

- 集合了 gin 的常用middleware、handler、以及对其进行极简封装
- 能直接使用的 middleware 和 handler 直接暴露出来供调用
- 不能直接使用的封装成 [Functional Option](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

```
    router := gi.New(
        gi.WithPprof(),
    )

    router.Use(
        gi.MidCORS(),
        gi.MidRecovery(),
        gi.MidLogger(),
        gi.MidStatic(),
    )

```
