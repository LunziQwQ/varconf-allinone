# varconf-allinone
Edit by varconf/varconf-server, varconf/varconf-ui, varconf/varconf-docker. Add some feature

## Origin Repos
* [varconf-server](https://github.com/varconf/varconf-server)
* [varconf-ui](https://github.com/varconf/varconf-ui)
* [varconf-docker](https://github.com/varconf/varconf-docker)

Thanks [@yawenok](https://github.com/yawenok). Really clean and powerful config center platform.

## Different
* Refactor cli by cobra.
* Make Docker use external mysql.
* Make db auto init by go not sql script.
* Add config init feature.
* (TODO) Add jsonschema on frontend.

# Dependencies
* MySQL
* No OS dependencies

# Build
```
go build -mod=vendor
```

# Usage
1. Edit config.json to set mysql dsn.
2. Start server.
```
varconf start -c ./config.json
```
3. All done: Visit `http://127.0.0.1:8088`. The default auth is `admin/123456`


If you want to clear all data and recreate. use `-r` arg.
```
varconf start -c ./config.json -r
```

If you want to start varconf with preset config. Edit `init.yaml` and use `-i` arg.
```
varconf start -c config.json -i .\init.yaml
```
> The `-i` args will create app and config only if app not exist. If you need to recreate all data. use `-r` args