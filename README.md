# Syncthing Hook

[![Release](https://img.shields.io/github/v/release/ichenhe/syncthing-hook?style=flat-square)](https://github.com/ichenhe/syncthing-hook/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/ichenhe/syncthing-hook?style=flat-square)](https://goreportcard.com/report/github.com/ichenhe/syncthing-hook)
[![Build State](https://img.shields.io/github/actions/workflow/status/ichenhe/syncthing-hook/go_test.yml?style=flat-square)](https://github.com/ichenhe/syncthing-hook/actions)


Do you want to notify your programs to reload when a file is updated? Maybe also need to limit the trigger frequency?
That is what SyncthingHook does.

This project adds hook function to [Syncthing](https://syncthing.net/) (open source file synchronization program), but
works as a standalone program. The basic principle is to constantly request
the [Syncthing API](https://docs.syncthing.net/dev/rest.html) to get the latest events and perform specified actions. So
SynchingHook is compatible with Syncthing on any platform, and can even listen to remote service.

**Features:**

- Single executable file, which is very easy to deploy.
- Supports many kinds of event, including Syncthing native event and extra event.
- Supports multiple filters depending on the kind of event.

# Quick Start

1. Download the executable file from the [release page](https://github.com/ichenhe/syncthing-hook/releases).
2. Copy the [example config file](https://github.com/ichenhe/syncthing-hook/blob/main/config/config.example.sthook.yaml)
   and write your version.
3. Run it with profile: `./SyncthingHook --profile /path/to/config.yaml`

# Configuration

## Profile

SyncthingHook relies on configuration file (profile), which specifies events of interest and corresponding
actions. You **must** specify the path of profile by one of the following ways (highest to lowest priority):

- Command line argument: `--profile /path/to/profile`
- Environment variable: `STHOOK_PROFILE=/path/to/profile`

Please refer to `config/config.example.sthook.yaml` for the format of profile. The detailed specification is also given
through `conofig/sthook.schema.json` (json-schema).

## Other Configuration Methods

### Priority

The priority of different configuration methods from high to low is as follows:

- Command line arguments
- Environment variables
- Profile

The value of the high-priority configuration method will override the configuration elsewhere, **even if its value is
empty ""**.

### Supported Configuration Items

> ⚠️ Warning: do not configure other items expect those listed below by command line arguments or environment as it
> may lead to undefined behavior and fail at any time!

Some items in profile can be configured by command line arguments or environment variables. The following table lists
all of them:

| Config Item      | Command Line Argument | Environment Variable      |
|------------------|-----------------------|---------------------------|
| syncthing.url    | `--syncthing.url`     | `STHOOK_SYNCTHING_URL`    |
| syncthing.apikey | `--syncthing.apikey`  | `STHOOK_SYNCTHING_APIKEY` |

For example：

```shell
./sthook --syncthing.url http://localhost:8384 --syncthing.apikey aaabbbccc
```

# Events

| Event                               | Description                                                                     |
|-------------------------------------|---------------------------------------------------------------------------------|
| st:ConfigSaved                      | [Syncthing doc](https://docs.syncthing.net/events/configsaved.html)             |
| st:DeviceConnected                  | [Syncthing doc](https://docs.syncthing.net/events/deviceconnected.html)         |
| st:DeviceDisconnected               | [Syncthing doc](https://docs.syncthing.net/events/devicedisconnected.html)      |
| st:DeviceDiscovered                 | [Syncthing doc](https://docs.syncthing.net/events/devicediscovered.html)        |
| st:DevicePaused                     | [Syncthing doc](https://docs.syncthing.net/events/devicepaused.html)            |
| st:DeviceResumed                    | [Syncthing doc](https://docs.syncthing.net/events/deviceresumed.html)           |
| st:DownloadProgress                 | [Syncthing doc](https://docs.syncthing.net/events/downloadprogress.html)        |
| st:Failure                          | [Syncthing doc](https://docs.syncthing.net/events/failure.html)                 |
| st:FolderCompletion                 | [Syncthing doc](https://docs.syncthing.net/events/foldercompletion.html)        |
| st:FolderErrors                     | [Syncthing doc](https://docs.syncthing.net/events/foldererrors.html)            |
| st:FolderPaused                     | [Syncthing doc](https://docs.syncthing.net/events/folderpaused.html)            |
| st:FolderResumed                    | [Syncthing doc](https://docs.syncthing.net/events/folderresumed.html)           |
| st:FolderScanProgress               | [Syncthing doc](https://docs.syncthing.net/events/folderscanprogress.html)      |
| st:FolderSummary                    | [Syncthing doc](https://docs.syncthing.net/events/foldersummary.html)           |
| st:FolderWatchStateChanged          | [Syncthing doc](https://docs.syncthing.net/events/folderwatchstatechanged.html) |
| st:ItemFinished                     | [Syncthing doc](https://docs.syncthing.net/events/itemfinished.html)            |
| st:ItemStarted                      | [Syncthing doc](https://docs.syncthing.net/events/itemstarted.html)             |
| st:ListenAddressesChanged           | [Syncthing doc](https://docs.syncthing.net/events/listenaddresseschanged.html)  |
| st:LocalChangeDetected              | [Syncthing doc](https://docs.syncthing.net/events/localchangedetected.html)     |
| st:LocalIndexUpdated                | [Syncthing doc](https://docs.syncthing.net/events/localindexupdated.html)       |
| st:LoginAttempt                     | [Syncthing doc](https://docs.syncthing.net/events/loginattempt.html)            |
| st:PendingDevicesChanged            | [Syncthing doc](https://docs.syncthing.net/events/pendingdeviceschanged.html)   |
| st:PendingFoldersChanged            | [Syncthing doc](https://docs.syncthing.net/events/pendingfolderschanged.html)   |
| st:RemoteChangeDetected             | [Syncthing doc](https://docs.syncthing.net/events/remotechangedetected.html)    |
| st:RemoteDownloadProgress           | [Syncthing doc](https://docs.syncthing.net/events/remotedownloadprogress.html)  |
| st:RemoteIndexUpdated               | [Syncthing doc](https://docs.syncthing.net/events/remoteindexupdated.html)      |
| st:Starting                         | [Syncthing doc](https://docs.syncthing.net/events/starting.html)                |
| st:StartupComplete                  | [Syncthing doc](https://docs.syncthing.net/events/startupcomplete.html)         |
| st:StateChanged                     | [Syncthing doc](https://docs.syncthing.net/events/statechanged.html)            |
| ex:LocalFolderContentChangeDetected | Based on `st:LocalChangeDetected`, with file path matcher.                      |

## Common Parameters

Unless otherwise noted, the following parameters apply to all events:

| Name      | Type  | Default | Description                                                                                                            |
|-----------|-------|---------|------------------------------------------------------------------------------------------------------------------------|
| tolerance | int64 | 0       | How long to wait before triggering this event (ms). Typically used to get only the latest event. `0` means no waiting. |
| cooldown  | int64 | 0       | The maximum frequency in millisecond that this event can be triggered. `0` indicates no limitation.                    |

Example:

```yaml
hooks:
  - event-type: "st:FolderCompletion"
    parameter:
      tolerance: 500
      cooldown: 3000
```

## ex:LocalFolderContentChangeDetected

This event based on `st:LocalChangeDetected`. It will be triggered only when the folder id equals the given one and the
file belongs to the given path (equal or subdirectory). Path pattern `/` matches all events.

**Parameters:**

| Name      | Type   | Default | Description                                        |
|-----------|--------|---------|----------------------------------------------------|
| st-folder | string |         | Folder-id in Syncthing. Cannot be omitted.         |
| path      | string | `/`     | Path to the target directory, must start with `/`. |

# Actions

The only action type currently supported is `exec`. However, since it can run any command, you may do whatever you
want (e.g. run a python script).

Example:

```yaml
action:
  type: "exec"
  cmd: [ "python", "do.py" ]
  # cmd: [ "python do.py" ] BAD!
```

Please note:

- The first element of `cmd` must be the program, which means you can't put parameters in it.
- SyncthingHook won't wait for the command to return. So it may execute the command (triggered by next event) while the
  previous one is still running.
