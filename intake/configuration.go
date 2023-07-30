package intake

import (
	"encoding/json"
	"fmt"
	"os"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

type IntakeConfiguration struct {
	DbPath string `json:"dbPath"`
}

func ReadIntakeConfiguration() (*IntakeConfiguration, error) {
	var processArgs = os.Args[1:]

	if len(processArgs) == 0 {
		return nil, errors.E("can not find first argument")
	}

	var configPath = processArgs[0]

	var data, err = os.ReadFile(configPath)
	if err != nil {
		return nil, errors.J(err, fmt.Sprintf("reading file failed at %s", configPath))
	}

	var configuration IntakeConfiguration
	err = json.Unmarshal(data, &configuration)
	if err != nil {
		return nil, errors.J(err, fmt.Sprintf("parsing config failed at %s", configPath))
	}

	return &configuration, nil
}
