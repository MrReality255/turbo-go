package log

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/MrReality255/turbo-go/tg/utils"

	log "github.com/sirupsen/logrus"
)

type Context map[string]any

type Config struct {
	Filename string
}

type logFileWriter struct {
	filename string
}

func (l *logFileWriter) Write(p []byte) (n int, err error) {
	_, _ = os.Stdout.Write(p)

	// if file exists, append
	if utils.FileExists(l.filename) {
		var f *os.File
		f, err = os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}

		err = utils.CloseAfter(f, func() error {
			n, err = f.Write(p)
			return err
		})
		return
	}

	return len(p), os.WriteFile(l.filename, p, 0644)
}

func Debug(msg string, params ...any) {
	log.Debug(fmt.Sprintf(msg, params...))
}

func IfError(msg string, err error) {
	if err != nil {
		LogError(msg, err)
	}
}

func LogIfWarnErr(msg string, err error) {
	if err != nil {
		log.Warning(msg, err)
	}
}

func Info(msg string, params ...any) {
	log.Info(fmt.Sprintf(msg, params...))
}

func LogError(msg string, err error) {
	if !strings.Contains(msg, "%") {
		msg = fmt.Sprintf("%v %%v", msg, err)
	}
	log.Errorf(msg, err)
}

func Warn(msg string, params ...any) {
	log.Warn(fmt.Sprintf(msg, params...))
}

func LogWarnErr(msg string, err error) {
	log.Warning(msg, err)
}

func LogIfErrCtx(msg string, err error, ctx Context) {
	if err != nil {
		log.Errorf("%v: %v%v", msg, err, ctx)
	}
}

func Setup(setup *Config) {
	if setup == nil {
		setup = &Config{}
	}
	lw := &logFileWriter{filename: setup.Filename}

	if setup.Filename != "" {
		log.SetOutput(lw)
	} else {
		log.SetOutput(os.Stdout)
	}
	log.SetLevel(log.DebugLevel)
}

func (c Context) String() string {
	if c == nil || len(c) == 0 {
		return ""
	}

	items := make([]string, 0, len(c))
	for k, v := range c {
		items = append(items, fmt.Sprintf("%s=%v", k, v))
	}
	sort.Strings(items)
	return "{" + strings.Join(items, ", ") + "}"

}
