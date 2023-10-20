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