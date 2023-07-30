package intake

import (
	"encoding/json"
	"fmt"
	"os"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

type IntakeConfigurationTransport struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

func (transport *IntakeConfigurationTransport) checkForErrors() error {
	if transport.Network != "unix" && transport.Network != "tcp" {
		return errors.E(fmt.Sprintf("invalid value for network field, got %s", transport.Network))
	}

	return nil
}

type IntakeConfiguration struct {
	DbPath    string                       `json:"dbPath"`
	Transport IntakeConfigurationTransport `json:"transport"`
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

	err = configuration.Transport.checkForErrors()
	if err != nil {
		return nil, errors.J(err, "config transport field check failed")
	}

	return &configuration, nil
}
