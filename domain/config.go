package domain

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type AppProfile struct {
	Syncthing struct {
		Url    string `koanf:"url"`
		ApiKey string `koanf:"apikey"`
	} `koanf:"syncthing"`
	Hooks []Hook `koanf:"hooks"`
}

type HookParameters map[string]any

type Hook struct {
	Name       string         `koanf:"name"`
	EventType  string         `koanf:"event-type"`
	Parameter  HookParameters `koanf:"parameter"`
	Conditions []struct {
		Type  string `koanf:"type"`
		Var   string `koanf:"var"`
		Value string `koanf:"value"`
	} `koanf:"conditions"`
	Action HookAction `koanf:"action"`
}

type HookAction struct {
	Type string   `koanf:"type"`
	Cmd  []string `koanf:"cmd"`
}

// HookDefinition represents a Hook item in the configuration.
type HookDefinition struct {
	Name  string
	Index int
}

func NewHookDefinition(name string, index int) *HookDefinition {
	return &HookDefinition{
		Name:  name,
		Index: index,
	}
}

// AddToLogger adds hook definition info to the given logger, returns new logger.
func (d *HookDefinition) AddToLogger(logger *zap.SugaredLogger) *zap.SugaredLogger {
	return logger.With(zap.String("hookName", d.Name), zap.Int("hookIndex", d.Index))
}

func LoadAppProfile(path string) (*AppProfile, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, err
	}
	appProfile := &AppProfile{}
	if err := k.Unmarshal("", appProfile); err != nil {
		return nil, err
	}
	return appProfile, nil
}

func (p HookParameters) ContainsKey(key string) (ex bool) {
	_, ex = p[key]
	return
}

func (p HookParameters) GetString(key string, defaultValue string) string {
	if v, ex := p[key]; !ex {
		return defaultValue
	} else if t, ok := v.(string); !ok {
		return defaultValue
	} else {
		return t
	}
}

func (p HookParameters) GetInt64(key string, defaultValue int64) int64 {
	if v, ex := p[key]; !ex {
		return defaultValue
	} else {
		switch i := v.(type) {
		case int64:
			return i
		case int:
			return int64(i)
		case int32:
			return int64(i)
		default:
			return defaultValue
		}
	}
}

func (p HookParameters) ExtractIntIfExist(key string, target *int) {
	if v, ex := p[key]; ex {
		if s, ok := v.(int); ok {
			*target = s
		}
	}
}

func (p HookParameters) ExtractInt64IfExist(key string, target *int64) {
	if v, ex := p[key]; ex {
		switch i := v.(type) {
		case int64:
			*target = i
		case int:
			*target = int64(i)
		case int32:
			*target = int64(i)
		}
	}
}

func (p HookParameters) ExtractStringIfExist(key string, target *string) {
	if v, ex := p[key]; ex {
		if s, ok := v.(string); ok {
			*target = s
		}
	}
}
