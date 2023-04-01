# Qog

> 打印单一形式的日志库

```log
2022-11-10 09:09:07.913|Demo|DEBUG|test/qog_test.go.TestQog:11|goid:6|23333
2022-11-10 09:09:07.913|Demo|INFO |test/qog_test.go.TestQog:12|goid:6|23333
2022-11-10 09:09:07.913|Demo|WARN |test/qog_test.go.TestQog:13|goid:6|23333
2022-11-10 09:09:07.913|Demo|ERROR|test/qog_test.go.TestQog:14|goid:6|23333
```

> 运行前需要在 `GOROOT/src/runtime` 文件夹下写入一个新函数

```Go
func Goid() uint64 {
	return getg().goid
}
```
