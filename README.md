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

# Build
```
go build -mod=vendor
```

# Run server
1. Edit config.json
2. Start server
```
varconf start -c ./config.json
```

# Run server with init config
> The `-i` args will create app and config only if app not exist. If you need to recreate all data. use `-r` args
1. Edit init.yaml
2. Start server
```
varconf start -c config.json -i .\init.yaml
```
