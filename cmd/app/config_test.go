package main

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_configurationLoader_loadConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		profileVal string
		envVal     string
		cmdArgsVal string
		wantErr    bool
		want       string
	}{
		{name: "read from file", profileVal: "key", want: "key"},
		{name: "cmdline override all others", profileVal: "key1", envVal: "key2", cmdArgsVal: "key3", want: "key3"},
		{name: "env override profile", profileVal: "key1", envVal: "key2", want: "key2"},
		{name: "no profile", profileVal: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configFile := filepath.Join(t.TempDir(), "profile.yaml")
			if len(tt.profileVal) > 0 {
				k := koanf.New(".")
				require.Nil(t, k.Set("syncthing.apikey", tt.profileVal))
				if data, err := k.Marshal(yaml.Parser()); err != nil {
					t.Fatalf("failed to marshal configuration file: %s", err.Error())
					return
				} else if err := os.WriteFile(configFile, data, 0644); err != nil {
					t.Fatalf("failed to write configuration file: %s", err.Error())
					return
				}
			}
			fetcher := &MockargumentFetcher{}
			fetcher.EXPECT().GetCommandLineArgs().RunAndReturn(func() []string {
				if len(tt.cmdArgsVal) > 0 {
					return []string{"app", "--syncthing.apikey", tt.cmdArgsVal}
				} else {
					return []string{"app"}
				}
			})
			loader := newConfigurationLoader(fetcher)
			if len(tt.envVal) > 0 {
				t.Setenv(_EnvPrefix+"SYNCTHING_APIKEY", tt.envVal)
			}
			t.Setenv(_ProfileEnv, configFile)

			config, err := loader.loadConfiguration()
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Nilf(t, err, "no error expected, but something went wrong: %s", err)
			assert.EqualValues(t, tt.want, config.Syncthing.ApiKey)
		})
	}
}

func Test_configurationLoader_loadConfigurationFromEnv(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantValues map[string]string
	}{
		{
			name:       "1 config item",
			env:        map[string]string{_EnvPrefix + "SYNCTHING_URL": "https://"},
			wantValues: map[string]string{"syncthing.url": "https://"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher := &MockargumentFetcher{}
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			loader := newConfigurationLoader(fetcher)
			k := koanf.New(".")
			err := loader.loadConfigurationFromEnv(k)

			require.Nil(t, err, "no error expected, but something went wrong")
			for key, v := range tt.wantValues {
				require.EqualValues(t, v, k.String(key))
			}
		})
	}
}

func Test_configurationLoader_loadConfigurationFromCmd(t *testing.T) {
	tests := []struct {
		name       string
		cmdArgs    []string
		wantValues map[string]string
	}{
		{
			name:       "1 config item",
			cmdArgs:    []string{"app", "--syncthing.url", "https://"},
			wantValues: map[string]string{"syncthing.url": "https://"},
		},
		{
			name:       "2 config items",
			cmdArgs:    []string{"app", "--syncthing.url", "https://", "--syncthing.apikey", "key"},
			wantValues: map[string]string{"syncthing.url": "https://", "syncthing.apikey": "key"},
		},
		{
			name:       "no item, take nothing",
			cmdArgs:    []string{"app"},
			wantValues: map[string]string{},
		},
		{
			name:       "key=val style",
			cmdArgs:    []string{"app", "--syncthing.apikey=key"},
			wantValues: map[string]string{"syncthing.apikey": "key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher := &MockargumentFetcher{}
			fetcher.EXPECT().GetCommandLineArgs().Return(tt.cmdArgs)
			loader := newConfigurationLoader(fetcher)

			k := koanf.New(".")
			err := loader.loadConfigurationFromCmd(k)

			require.Nil(t, err, "no error expected, but something went wrong")
			for key, v := range tt.wantValues {
				require.EqualValues(t, v, k.String(key))
			}
			if len(tt.wantValues) == 0 {
				assert.Falsef(t, k.Exists("syncthing.apikey"), "should not take ungiven argument")
			}
		})
	}
}

func Test_configurationLoader_readProfilePath(t *testing.T) {
	tests := []struct {
		name            string
		cmdArgs         []string
		envValue        string
		wantProfilePath string
	}{
		{name: "read from cmdline args", cmdArgs: []string{"app", "--profile", "/a.yaml"}, envValue: "", wantProfilePath: "/a.yaml"},
		{name: "read from env", cmdArgs: []string{"app"}, envValue: "/b.yaml", wantProfilePath: "/b.yaml"},
		{name: "cmdline args should overwrite env", cmdArgs: []string{"app", "--profile", "/a.yaml"}, envValue: "/b.yaml", wantProfilePath: "/a.yaml"},
		{name: "no value given", cmdArgs: []string{"app"}, envValue: "", wantProfilePath: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher := &MockargumentFetcher{}
			fetcher.EXPECT().GetCommandLineArgs().Return(tt.cmdArgs)
			t.Setenv(_ProfileEnv, tt.envValue)
			loader := newConfigurationLoader(fetcher)
			if gotProfilePath := loader.readProfilePath(); gotProfilePath != tt.wantProfilePath {
				t.Errorf("readProfilePath() = %v, want %v", gotProfilePath, tt.wantProfilePath)
			}
		})
	}
}
