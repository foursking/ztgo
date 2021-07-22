# Event


## API

注册同步事件
```
event.On(event string, l ListenerFunc)
```

注册异步事件（不保证执行顺序）

```
event.AsyncOn(event string, l ListenerFunc)
```

触发事件
```
event.Emit(ctx context.Context, event string, v interface{})
```

移除事件监听器
```
event.Off(event string, l ListenerFunc)
```

获取事件监听器列表
```
event.Listeners(event string) []Listener
```

获取所有事件以及监听器列表
```
event.AllListeners() map[string][]Listener
```


## Demo

```go
import "github.com/foursking/ztgo/event"

event.On("user.created", func(v ...interface{}) {
    log.Infof("User '%s' was created", v[0].Username)
})

event.Emit(context.TODO(), "user.created", new(User))
```
