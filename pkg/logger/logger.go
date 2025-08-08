package logger

import (
	"kafka-client/pkg/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局Logger
var lg *zap.Logger

// InitLogger 根据提供的配置初始化Logger
func InitLogger(cfg *config.AppConfig) (err error) {
	// 1. 设置日志级别
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(cfg.LogConfig.Level)); err != nil {
		return err
	}

	// 2. 配置 lumberjack
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.LogConfig.LumberjackConfig.Filename,
		MaxSize:    cfg.LogConfig.LumberjackConfig.MaxSize,
		MaxBackups: cfg.LogConfig.LumberjackConfig.MaxBackups,
		MaxAge:     cfg.LogConfig.LumberjackConfig.MaxAge,
		Compress:   cfg.LogConfig.LumberjackConfig.Compress,
	}

	// 3. 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.LogConfig.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 4. 创建Core
	// 同时写入文件和控制台
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger)),
		level,
	)

	// 5. 创建Logger
	// 如果启用了调用者信息，则添加 zap.AddCaller()
	var opts []zap.Option
	if cfg.LogConfig.EnableCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1)) // AddCallerSkip(1) 补偿封装导致的调用栈深度增加
	}

	// 配置错误日志也输出到指定位置
	opts = append(opts, zap.ErrorOutput(zapcore.AddSync(lumberjackLogger)))

	lg = zap.New(core, opts...)

	// 替换zap全局的Logger
	zap.ReplaceGlobals(lg)

	zap.L().Info("logger init success")
	return nil
}

// 提供便捷的日志记录方法
func Debug(msg string, fields ...zap.Field) {
	lg.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	lg.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	lg.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	lg.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	lg.Fatal(msg, fields...)
}
