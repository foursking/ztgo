
errs 是 qdgo 统一的业务码存放处，将错误码集中定义，有利于归类和查阅。

## 新增错误码

如果在一个文件定义所有的错误码，当错误码越来越多的时候，势必会面临容易命名重复的问题。如项目A定义了一个 `ErrUserNotFound` 来表示用户不存在，项目B也想使用这个错误码，这种情况下就不得不重新起名或加前缀后缀。

为了尽可能减少类似场景，本仓库错误码，建议按业务分级存放，比如 `account` 表示账户相关的错误码。

新建错误码很简单，只需要指定 code 和 message 即可：

```go
var ErrUserNotFound = New(10086, "用户不存在")
```

## 判断错误码是否相同

```go
if errs.Equal(err1, err2) {
    println("两个错误码 code 和 message 均相同")
}
```

## 为错误码携带更多的信息

为一个错误码携带更多信息，将会复制这个错误码，加上附加信息后，生成一个新的错误码对象，从而保证对原错误码不产生修改影响。

**以下提供几个真实示例：**

```go
// 请求参数错误，但是要告诉用户具体是什么类型的错误
err := errs.ErrBadRequest.WithDetail("detail", "用户名长度不符合规范")
```

此时 err 信息如下：

```javascript
{
    "code": -400,
    "message": "Bad Request",
    "details": {
        "detail": "用户名长度不符合规范"
    }
}
```

上面这种情况还支持另一种写法，提供了 `WithDetails` 方法`：

```go
details := map[string]string{
    "detail": "用户名长度不符合规范",
    "max": "10",
    "min", "2"
}
err := errs.ErrBadRequest.WithDetails(details)
```
