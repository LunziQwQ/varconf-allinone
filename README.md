# varconf-allinone
Edit by varconf/varconf-server, varconf/varconf-ui, varconf/varconf-docker. Add some feature

## Origin Repos
* [varconf-server](https://github.com/varconf/varconf-server)
* [varconf-ui](https://github.com/varconf/varconf-ui)
* [varconf-docker](https://github.com/varconf/varconf-docker)

Thanks [@yawenok](https://github.com/yawenok). Really clean and powerful config center platform.

## Different
* Refactor cli by cobra.
* (TODO) Make db init by go not sql script
* (TODO) Add config init feature.
* (TODO) Add jsonschema on frontend.

# Build
```
go build -mod=vendor
```

# start
1. Edit config.json
2. Start server
```
varconf start -c ./config.json
```

