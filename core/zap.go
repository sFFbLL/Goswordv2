package core

import (
	"fmt"
	"os"
	"project/global"
	"project/utils"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Level zapcore.Level

func Zap() (logger *global.NewLogger) {
	if ok, _ := utils.PathExists(global.GSD_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.GSD_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GSD_CONFIG.Zap.Director, os.ModePerm)
	}

	switch global.GSD_CONFIG.Zap.Level { // 初始化配置文件的Level
	case "debug":
		Level = zap.DebugLevel
	case "info":
		Level = zap.InfoLevel
	case "warn":
		Level = zap.WarnLevel
	case "error":
		Level = zap.ErrorLevel
	case "dpanic":
		Level = zap.DPanicLevel
	case "panic":
		Level = zap.PanicLevel
	case "fatal":
		Level = zap.FatalLevel
	default:
		Level = zap.InfoLevel
	}
	logger = &global.NewLogger{
		ZapLog: &zap.Logger{},
	}
	if Level == zap.DebugLevel || Level == zap.ErrorLevel {
		logger.ZapLog = zap.New(getEncoderCore(), zap.AddStacktrace(Level))
	} else {
		logger.ZapLog = zap.New(getEncoderCore())
	}
	if global.GSD_CONFIG.Zap.ShowLine {
		logger.ZapLog = logger.ZapLog.WithOptions(zap.AddCaller())
	}
	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.GSD_CONFIG.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case global.GSD_CONFIG.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case global.GSD_CONFIG.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case global.GSD_CONFIG.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case global.GSD_CONFIG.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if global.GSD_CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer, err := utils.GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, Level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(global.GSD_CONFIG.Zap.Prefix + "2006/01/02 - 15:04:05.000"))
}
