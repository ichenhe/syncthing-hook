package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/basicflag"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"strings"
)

const _EnvPrefix = "STHOOK_"
const _ProfileEnv = _EnvPrefix + "PROFILE"

var errNoProfileSpecified = errors.New("no profile specified")

type argumentFetcher interface {
	GetCommandLineArgs() []string
}

type defaultArgumentFetcher struct {
}

var _ argumentFetcher = (*defaultArgumentFetcher)(nil)

func (f *defaultArgumentFetcher) GetCommandLineArgs() []string {
	return os.Args
}

type configurationLoader struct {
	flagSet    *flag.FlagSet
	argFetcher argumentFetcher
}

func newConfigurationLoader(argFetcher argumentFetcher) *configurationLoader {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() {
	}
	fs.String("profile", "", "path to .yaml config file")
	fs.String("syncthing.url", "", "the url of Syncthing api, starts with http(s)://")
	fs.String("syncthing.apikey", "", "the api key of Syncthing, can be obtained from webui")
	return &configurationLoader{
		flagSet:    fs,
		argFetcher: argFetcher,
	}
}

// loadConfiguration loads the configurations according to the following priorities (high to low):
// 1.cmd; 2.env; 3.file
func (l *configurationLoader) loadConfiguration() (*domain.AppProfile, error) {
	k := koanf.New(".")

	// read profile path
	profilePath := l.readProfilePath()
	if len(profilePath) == 0 {
		return nil, errNoProfileSpecified
	}

	// from file
	if err := k.Load(file.Provider(profilePath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config from file '%s': %w", profilePath, err)
	}

	if err := l.loadConfigurationFromEnv(k); err != nil {
		return nil, fmt.Errorf("failed to read config from env variable: %w", err)
	}
	if err := l.loadConfigurationFromCmd(k); err != nil {
		return nil, fmt.Errorf("failed to read config from cmd line: %w", err)
	}

	// parse
	appProfile := &domain.AppProfile{}
	if err := k.Unmarshal("", appProfile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return appProfile, nil
}

// readProfilePath reads profile path from cmd line and env (priority high to low).
// If all fail, the empty string is returned.
func (l *configurationLoader) readProfilePath() (profilePath string) {
	// from cmd line
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() {
	}
	fs.StringVar(&profilePath, "profile", "", "path to .yaml config file")
	_ = fs.Parse(l.argFetcher.GetCommandLineArgs()[1:])
	if len(profilePath) > 0 {
		return
	}

	// from env
	profilePath = os.Getenv(_ProfileEnv)
	if len(profilePath) > 0 {
		return
	}
	return
}

func (l *configurationLoader) loadConfigurationFromEnv(k *koanf.Koanf) error {
	return k.Load(env.Provider(_EnvPrefix, "_", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, _EnvPrefix))
	}), nil)
}

func (l *configurationLoader) loadConfigurationFromCmd(k *koanf.Koanf) error {
	_ = l.flagSet.Parse(l.argFetcher.GetCommandLineArgs()[1:])
	set := make(map[string]struct{})
	l.flagSet.Visit(func(f *flag.Flag) {
		set[f.Name] = struct{}{}
	})
	provider := basicflag.ProviderWithValue(l.flagSet, ".", func(key string, value string) (string, interface{}) {
		if _, ex := set[key]; ex {
			return key, value // ignore flags not given
		}
		return "", ""
	})
	return k.Load(provider, nil)
}

func (l *configurationLoader) printUsage() {
	println("Usage:")
	l.flagSet.PrintDefaults()
}
