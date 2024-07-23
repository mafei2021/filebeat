package initialize

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

// initLog 日志的统一结构处理
func initLog() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// JSON 格式输出，并启用 Caller 信息
	//log.Logger = log.Output(os.Stdout).With().Caller().Logger()

	// 切换到文本格式输出，并启用 Caller 信息，不需要输出全路径，只需要输出文件名+代码行号就可以
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.TimeOnly}).With().Caller().Logger()

	switch viper.GetString("log.level") {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

}
