package log

import (
	"fmt"
	"os"

	"github.com/MrReality255/turbo-go/tg/utils"

	log "github.com/sirupsen/logrus"
)

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

func LogDebug(msg string, params ...any) {
	log.Debug(fmt.Sprintf(msg, params...))
}

func LogIfError(msg string, err error) {
	if err != nil {
		log.Error(msg, err)
	}
}

func LogInfo(msg string, params ...any) {
	log.Info(fmt.Sprintf(msg, params...))
}

func LogError(msg string, err error) {
	log.Error(msg, err)
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
