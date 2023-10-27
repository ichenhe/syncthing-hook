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

| Event                               | Description        |
|-------------------------------------|--------------------|
| ex:LocalFolderContentChangeDetected | Listen to an event |

## ex:LocalFolderContentChangeDetected

Parameters:

| Name      | Type   | Default | Description                                                                                                                                                                              |
|-----------|--------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| st-folder | string |         | Folder-id in Syncthing.                                                                                                                                                                  |
| path      | string | /       | Path to the target directory, must start with `/`.                                                                                                                                       |
| tolerance | int64  | 1000    | How many milliseconds to wait before triggering this event after a file under target folder changed detected. In case there are subsequence events of other files. `0` means no waiting. |
| cooldown  | int64  | 500     | The maximum frequency in millisecond that this event can be triggered. All triggering during cooldown period will be discarded. `0` indicates no cooldown needed.                        |