package logger

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	Key       = "logger"
	RequestID = "requestId"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string))) // добавляем в логгер requestID
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) addLogMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	ctx, _ = New(ctx)
	reply, err := handler(ctx, req)
	l.l.Info("gRPC top-level log demonstration!")
	return reply, err
}
