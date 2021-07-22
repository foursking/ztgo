
ParseDuration 是标准库 `time.ParseDuration` 的增强版，将字符串描述的时长转换成 `time.Duration`

## ParseDuration() 支持的时间单位

| Unit of time | value |
| --- | --- |
| ns | time.Nanosecond |
| us | time.Microsecond |
| ms | time.Millisecond |
| s | time.Second |
| sec | time.Second |
| second | time.Second |
| seconds | time.Second |
| m | time.Minute |
| min | time.Minute |
| minute | time.Minute |
| minutes | time.Minute |
| h | time.Hour |
| hour | time.Hour |
| hours | time.Hour |
| d | time.Hour * 24 |
| day | time.Hour * 24 |
| days | time.Hour * 24 |
| w | time.Hour * 24 * 7 |
| week | time.Hour * 24 * 7 |
| weeks | time.Hour * 24 * 7 |
