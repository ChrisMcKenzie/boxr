package utils

import (
	"os"

	"github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	"%{color}[%{level:.4s}] %{id:03x}%{color:reset} %{message}",
)

func Logger() *logging.Logger {
	log := logging.MustGetLogger("boxr")
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	// Only errors and more severe messages should be sent to backend
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled, backendFormatter)
	return log
}
