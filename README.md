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
| syncthing.url    | `--syncthing-url`     | `STHOOK_SYNCTHING_URL`    |
| syncthing.apikey | `--syncthing-apikey`  | `STHOOK_SYNCTHING_APIKEY` |

For example：

```shell
./sthook --syncthing-url http://localhost:8384 --syncthing-apikey aaabbbccc
```

# Events

> In active development, all [Syncthing native events](https://docs.syncthing.net/dev/events.html#event-types) will be
> supported soon.

| Event                               | Description                                                       |
|-------------------------------------|-------------------------------------------------------------------|
| st:FolderCompletion                 | [doc](https://docs.syncthing.net/events/foldercompletion.html)    |
| st:LocalChangeDetected              | [doc](https://docs.syncthing.net/events/localchangedetected.html) |
| ex:LocalFolderContentChangeDetected | Based on `st:LocalChangeDetected`, with file path matcher.        |

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