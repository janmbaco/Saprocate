package config

import (
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/janmbaco/go-reverseproxy-ssl/disk"
)

type(
	pairCertFiles struct{
		PublicKeyFile string `json:"public_key_file"`
		PrivateKeyFile string `json:"private_key_file"`
	}
	
	logInfo struct {
		ConsoleLevel cross.LogLevel `json:"console_level"`
		FileLevel    cross.LogLevel `json:"file_level"`
		LogsDir      string         `json:"logs_dir"`
	}
	config struct{
		Port string `json:"port"`
		AuthorizeCertificate pairCertFiles `json:"autorize_certificate"`
		ServerCertificate pairCertFiles `json:"server_certificate"`
		LogInfo logInfo `json:"log_info"`
		DataDir string `json:"data_dir"`
	}
)

var Config *config

func init(){
	// business defaults
	Config = &config{
		Port:                 ":5555",
		AuthorizeCertificate: pairCertFiles{
			PublicKeyFile:  "../certs/ca.pem",
			PrivateKeyFile: "../certs/ca.key",
		},
		ServerCertificate:    pairCertFiles{
			PublicKeyFile:  "../certs/server.crt",
			PrivateKeyFile: "../certs/server.key",
		},
		LogInfo:              logInfo{
			ConsoleLevel: cross.Trace,
			FileLevel:    cross.Info,
		},
		DataDir: "../data/blochchain",
	}

	disk.ConfigFile.ConstructorContent = func() interface{} {
		return &config{}
	}

	disk.ConfigFile.CopyContent = func(from interface{}, to interface{}) {
		fromConfig := from.(*config)
		toConfig := to.(*config)
		*toConfig = *fromConfig
	}
}

