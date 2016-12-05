package logging

import (
	"io"

	logging "github.com/op/go-logging"
)

var (
	logger = logging.MustGetLogger("animeshotbot")
	format = logging.MustStringFormatter("%{color}%{time:15:04:05.000}" +
		"%{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")
)

//Logger get a global logger
func Logger() *logging.Logger {
	return logger
}

//Setup setup log level
func Setup(backend io.Writer, debug bool) {
	logBackend := logging.NewLogBackend(backend, "", 0)

	if debug {
		logging.SetBackend(logBackend)
	} else {
		logBackendLeveled := logging.AddModuleLevel(logBackend)
		logBackendLeveled.SetLevel(logging.ERROR, "")
		logging.SetBackend(logBackendLeveled)
	}

}
