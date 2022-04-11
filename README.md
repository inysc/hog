# Qog

> 打印单一形式的日志库

```log
[2022-03-26 21:24:52.311][qog][TRACE][qog/qog_test.go.TestLevel:38][goid:6]abc
[2022-03-26 21:24:52.311][qog][WARN][qog/qog_test.go.TestLevel:39][goid:6]abc
[2022-03-26 21:24:52.311][qog][DEBUG][qog/qog_test.go.TestLevel:40][goid:6]abc
[2022-03-26 21:24:52.311][qog][INFO][qog/qog_test.go.TestLevel:41][goid:6]abc
[2022-03-26 21:24:52.311][qog][ERROR][qog/qog_test.go.TestLevel:42][goid:6]abc
```

> 运行前需要在 `GOROOT/src/runtime` 文件夹下写入一个新函数

```Go
func Goid() int64 {
	return getg().goid
}
```