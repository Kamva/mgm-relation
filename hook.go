package mgmrel

import "github.com/Kamva/mgm/v3"

// SyncingHook is the interface to implement hook to call before sync your model.
type SyncingHook interface {
	Syncing() error
}

// SyncedHook is the interface to implement hook to call after sync your model.
type SyncedHook interface {
	Synced() error
}

func callToBeforeSyncHooks(m mgm.Model) error {
	if hook, ok := m.(SyncingHook); ok {
		return hook.Syncing()
	}
	return nil
}

func callToAfterSyncHooks(m mgm.Model) error {
	if hook, ok := m.(SyncedHook); ok {
		return hook.Synced()
	}
	return nil
}
