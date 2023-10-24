package loginit

import (
	"github.com/rosaekapratama/go-starter/constant/env"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log/formatter/gcp"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var (
	Logger         logrus.Ext1FieldLogger
	standardLogger *logrus.Logger
)

func init() {
	// Set default log for all init/Init function logging.
	// This specific loginit.logger need to be used,
	// when log.logger is not initialized yet.
	standardLogger = logrus.StandardLogger()
	SetProjectId(str.Empty)
	Logger = standardLogger
}

func isRunLocally(logger logrus.StdLogger) bool {
	if localRunStr, ok := os.LookupEnv(env.EnvLocalRun); localRunStr != str.Empty && ok {
		localRun, err := strconv.ParseBool(localRunStr)
		if err != nil {
			logger.Print("Failed to parse %s env var '%s' to boolean, %s", env.EnvLocalRun, localRunStr, err.Error())
		} else {
			return localRun
		}
	}

	return false
}

func SetProjectId(projectId string) {
	var jsonFormatter = gcp.JSONFormatter{ProjectId: projectId}
	if isRunLocally(standardLogger) {
		standardLogger.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			ForceQuote:                true,
			EnvironmentOverrideColors: true,
			FullTimestamp:             true,
		})
	} else {
		standardLogger.SetFormatter(&jsonFormatter)
	}
}
