# zap-log-wrapper

zap作为golang的日志库，以其卓越的性能已经被众多知名项目采用。但是我在使用zap的时候，官方使用方式总觉得很不习惯，并且还需要自己配置日志切割。因此为了方便日后的使用，依赖zap库的一些代码的方法，对zap库做了一层包装。

## 1. 日志输出到stdout


```golang
NewLogger(&LoggerConfiguration{
		Level:       "debug",
		Development: false,
		Encoding:    "text",
		OutputPath:  StdOut,
	})
```

```yaml
log:
  level: "debug"
  development: false
  encoding: "text"
  output_path: "stdout"
```


## 2. 日志输出到文件，并支持日志轮转

```golang
NewLogger(&LoggerConfiguration{
		Level:       "debug",
		Development: false,
		Encoding:    "json",
		OutputPath:  "/data/logs/std.log",
		Rotate: LogRotate{
			MaxSizeMB:  1,
			MaxBackups: 3,
			Compress:   true,
		},
	})
```

```yaml
log:
  level: "debug"
  development: false
  encoding: "text"
  output_path: "/data/logs/std.log"
  rotate:
    max_size_mb: 1
    max_backups: 3
    compress: true
```