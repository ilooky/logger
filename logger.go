package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	logger  Logger
	wrapper *loggerWrapper
)

type Config struct {
	Style   string // consul or json
	Path    string // absolute path or filename
	Level   string // log level: debug\info\warn\error
	Release bool   // false=develop or true=product
}

type loggerWrapper struct {
	Logger
	zp *zap.Logger
}

func init() {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(getDevelopEncoder()),
		getWriteSync(""),
		getLogLevel("debug"))
	zp := zap.New(core, zap.AddCallerSkip(1))
	wrapper = &loggerWrapper{Logger: zp.Sugar(), zp: zp}
	logger = wrapper
}

func InitLogger(conf Config) {
	var encoder zapcore.Encoder
	var encoderConfig zapcore.EncoderConfig
	if conf.Release {
		encoderConfig = getProductEncoder()
	} else {
		encoderConfig = getDevelopEncoder()
	}
	if conf.Style == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(
		encoder,
		getWriteSync(conf.Path),
		getLogLevel(conf.Level))
	zp := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	wrapper = &loggerWrapper{Logger: zp.Sugar(), zp: zp}
	logger = wrapper
}

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
}

func getLogLevel(level string) zapcore.Level {
	if zapLevel, ok := levelMap[level]; ok {
		return zapLevel
	}
	return zapcore.InfoLevel
}

func getDevelopEncoder() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
func getProductEncoder() zapcore.EncoderConfig {
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	return encoder
}

//logOutPath + string(os.PathSeparator) + "us_diagram.log",
func getWriteSync(logOutPath string) zapcore.WriteSyncer {
	var syncer zapcore.WriteSyncer
	if logOutPath != "" {
		hook := lumberjack.Logger{
			Filename:   logOutPath,
			MaxSize:    50,
			MaxBackups: 15,
			MaxAge:     7,
			Compress:   true,
		}
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	return syncer
}

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Panic(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
	Panicf(fmt string, args ...interface{})
}

func SetLogger(log Logger) {
	logger = log
}

func GetLogger() Logger {
	return logger
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Debugf(fmt string, args ...interface{}) {
	logger.Debugf(fmt, args...)
}
func Infof(fmt string, args ...interface{}) {
	logger.Infof(fmt, args...)
}

func Warnf(fmt string, args ...interface{}) {
	logger.Warnf(fmt, args...)
}

func Errorf(fmt string, args ...interface{}) {
	logger.Errorf(fmt, args...)
}

func Panicf(fmt string, args ...interface{}) {
	logger.Panicf(fmt, args...)
}

func DebugKV(msg, key string, val interface{}) {
	wrapper.zp.Debug(msg, zap.Any(key, val))
}
func DebugKv(msg string, fields ...zap.Field) {
	wrapper.zp.Debug(msg, fields...)
}

func ErrorKV(msg, key string, val interface{}) {
	wrapper.zp.Error(msg, zap.Any(key, val))
}

func ErrorKv(msg string, fields ...zap.Field) {
	wrapper.zp.Error(msg, fields...)
}

func InfoKV(msg, key string, val interface{}) {
	wrapper.zp.Info(msg, zap.Any(key, val))
}

func InfoKv(msg string, fields ...zap.Field) {
	wrapper.zp.Info(msg, fields...)
}

func PanicKV(msg, key string, val interface{}) {
	wrapper.zp.Panic(msg, zap.Any(key, val))
}

func PanicKv(msg string, fields ...zap.Field) {
	wrapper.zp.Panic(msg, fields...)
}
