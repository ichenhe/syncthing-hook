package domain

// HookManager registers hook and executes action.
type HookManager interface {
	RegisterHook(hook *Hook, hookDef *HookDefinition) error
	UnregisterAll()
}
