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

### clear

The format of Post request is:

```json
{
  "clear_all": true
}

