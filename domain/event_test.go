package domain

import (
	"errors"
	"github.com/syncthing/syncthing/lib/events"
	"testing"
)

func TestEventType_ConvertToNative(t *testing.T) {
	tests := []struct {
		name    string
		t       EventType
		want    events.EventType
		wantErr bool
	}{
		{name: "native-wrapper event", t: LocalChangeDetected, want: events.LocalChangeDetected, wantErr: false},
		{name: "non-native event", t: LocalFolderContentChangeDetected, want: events.EventType(-716), wantErr: true},
		{name: "invalid input", t: EventType("abc"), want: events.EventType(-716), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.ConvertToNative()
			if tt.wantErr && !errors.Is(err, ErrNotValidNativeEventType) {
				t.Errorf("ConvertToNative() error = %v, want %v", err, ErrNotValidNativeEventType)
				return
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("ConvertToNative() error = %v, want success", err)
					return
				}
				if got != tt.want {
					t.Errorf("ConvertToNative() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestEventType_IsNativeEvent(t *testing.T) {
	tests := []struct {
		name string
		t    EventType
		want bool
	}{
		{name: "native event", t: FolderCompletion, want: true},
		{name: "non-native event", t: LocalFolderContentChangeDetected, want: false},
		{name: "invalid input", t: EventType("abc"), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsNativeEvent(); got != tt.want {
				t.Errorf("IsNativeEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalEventType(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want EventType
	}{
		{name: "native event", str: "st:FolderCompletion", want: FolderCompletion},
		{name: "ex event", str: "ex:LocalFolderContentChangeDetected", want: LocalFolderContentChangeDetected},
		{name: "invalid input", str: "st:invalid", want: UnknownEventType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnmarshalEventType(tt.str); got != tt.want {
				t.Errorf("UnmarshalEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertFromStEventType(t *testing.T) {
	tests := []struct {
		name string
		stEv events.EventType
		want EventType
	}{
		{name: "normal", stEv: events.LocalChangeDetected, want: LocalChangeDetected},
		{name: "unknown st event type", stEv: events.EventType(-716), want: UnknownEventType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertFromStEventType(tt.stEv); got != tt.want {
				t.Errorf("convertFromStEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}
