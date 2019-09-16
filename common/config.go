package common

import (
	"flag"
	"github.com/koding/multiconfig"
	"os"
)

type FlagConfig struct {
	ConfigFile string `default:"config.json"`
}
type LoggerConfig struct {
	Enabled    bool `default:true`
	LogFile    string
	TraceLevel int `default:3`
}
type ServerConfig struct {
	Listen      string `default:":5000"`
	RuntimePath string `default:"runtime"`
	Db struct {
		Host        string
		Port        string `default:"3306"`
		Name        string
		User        string
		Password    string
		SlaveConfig struct {
			User     string
			Password string
		}
		Slaves []struct {
			Host string
			Port string
			Name string
		}
		MaxOpenConns int  `default:"0"`
		ShowSQL      bool `default:"false"`
	}
	Log struct {
		LogFile    string `default:""`
		SaveType   string `default:"d"`
		TraceLevel int    `default:"3"`
		Logger     struct {
			Trace LoggerConfig
			Info  LoggerConfig
			Warn  LoggerConfig
			Error LoggerConfig
		}
	}
	System struct{
		Debug bool `default:"false"`
	}
}

func (c *FlagConfig) load() error {
	t := &multiconfig.TagLoader{}
	f := &multiconfig.FlagLoader{}
	m := multiconfig.MultiLoader(t, f)
	if err := m.Load(c);err == flag.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		return err

	}
	return nil
}
func (c *ServerConfig) load() error {
	f := &FlagConfig{}
	err := f.load()
	if err == flag.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		return err
	}
	t   := &multiconfig.TagLoader{}
	j   := &multiconfig.JSONLoader{Path:f.ConfigFile}
	m   := multiconfig.MultiLoader(t, j)
	err =m.Load(c)
	return err
}
