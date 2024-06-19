# sls

Simple Lock Service

`sls` works as a web server which has following endpoints

- /lock
- /unlock
- /clear

### lock / unlock

The format of Post request is:

```json
{
  "uuid": "xxx-yyy-zzz-...",
  "path": "file path"
}
```
- `uuid` is used to identify the pair of lock and unlock. 
- `path` is used to identify a resouce name which is locked and unlocked.
- The order in which lock requests for a specific `path` are processed is the order in which they arrive.

### clear

The format of Post request is:

```json
{
  "clear_all": true
}
```

- If `clear_all` is true, all pending lock requests will be aborted and StatusGone(410) will be returned.
