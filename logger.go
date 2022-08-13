package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Init() error {
	c := zap.NewDevelopmentConfig()
	c.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zl, err := c.Build()
	if err != nil {
		return err
	}

	l.SugaredLogger = zl.Sugar()

	return nil
}

// // Config defines the config for middleware
// type Config struct {
// 	// h fiber.Handler
// 	// Next defines a function to skip this middleware when returned true.
// 	//
// 	// Optional. Default: nil
// 	Next func(c *fiber.Ctx) bool

// 	// Logger defines zap logger instance
// 	Logger *zap.Logger
// }

// func NewLoggerMiddleware(cfg Config) fiber.Handler {
// 	// Set variables
// 	var (
// 		once sync.Once
// 		// mu         sync.Mutex
// 		errHandler fiber.ErrorHandler
// 		errPadding = 15
// 	)

// 	return func(c *fiber.Ctx) (err error) {
// 		// Don't execute middleware if Next returns true
// 		if cfg.Next != nil && cfg.Next(c) {
// 			return c.Next()
// 		}

// 		// Set error handler once
// 		once.Do(func() {
// 			// get longested possible path
// 			stack := c.App().Stack()
// 			for m := range stack {
// 				for r := range stack[m] {
// 					if len(stack[m][r].Path) > errPadding {
// 						errPadding = len(stack[m][r].Path)
// 					}
// 				}
// 			}
// 			// override error handler
// 			errHandler = c.App().ErrorHandler
// 		})

// 		var start, stop time.Time

// 		start = time.Now()

// 		// Handle request, store err for logging
// 		chainErr := c.Next()

// 		// Manually call error handler
// 		if chainErr != nil {
// 			if err := errHandler(c, chainErr); err != nil {
// 				_ = c.SendStatus(fiber.StatusInternalServerError)
// 			}
// 		}

// 		stop = time.Now()

// 		buf := bytes.NewBufferString(fmt.Sprintf(
// 			"| %d | \t %s | \t %s | %s \t | %s",
// 			c.Response().StatusCode(),
// 			stop.Sub(start).String(),
// 			c.IP(),
// 			c.Method(),
// 			c.Path(),
// 		))

// 		formatErr := ""
// 		if chainErr != nil {
// 			formatErr = chainErr.Error()

// 			buf.WriteString("\t | " + formatErr)
// 			cfg.Logger.Error("err",
// 				zap.String("time", stop.Sub(start).Round(time.Millisecond).String()),
// 				zap.Int("code", c.Response().StatusCode()),
// 				zap.String("ip", c.IP()),
// 				zap.String("method", c.Method()),
// 				zap.String("path", c.Path()),
// 				zap.String("error", formatErr),
// 			)

// 			return nil
// 		}

// 		cfg.Logger.Info(buf.String(),
// 			zap.String("time", stop.Sub(start).String()),
// 			zap.Int("code", c.Response().StatusCode()),
// 			zap.String("ip", c.IP()),
// 			zap.String("method", c.Method()),
// 			zap.String("path", c.Path()),
// 		)
// 		return nil
// 	}

// }

// type Request struct {
// 	Time   time.Duration
// 	Code   int
// 	IP     string
// 	Method string
// 	Path   string
// 	Error  string
// }

// // func RequestEncoder(c fiber.Ctx, enc zapcore.PrimitiveArrayEncoder) {
// // 	enc.AppendString("[" + c.Method() + "]")
// // 	// zapcore.NameEncoder

// // 	// zapcore.CapitalColorLevelEncoder()
// // }

// func (r Request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
// 	enc.AddString("method", r.Method)
// 	enc.AddDuration("latency", r.Time)
// 	return nil
// }

// func (r Request) MarshalLogArray(enc zapcore.ArrayEncoder) error {
// 	enc.AppendString(r.Method)
// 	enc.AppendInt(r.Code)
// 	enc.AppendDuration(r.Time)
// 	return nil
// }
