# lever

An audited feature flagging service.

To run without any external dependencies, you can run either of these:

* `$ LEVER_RELEASEMODE=false LEVER_STORE=memory go run *.go` (from source)
* `$ LEVER_RELEASEMODE=false LEVER_STORE=memory lever_0.1.0_darwin_amd64` (from binary)

## Development

1. Ensure running Go 1.5+
2. Install [glide](https://github.com/Masterminds/glide#install)
3. `export GO15VENDOREXPERIMENT=1`
4. `glide i`
5. `go run *.go`

### Tests

* `go vet $(glide nv)`
* `go test -v -race $(glide nv)`

## Config

Config options can be defined either via a `config.yml` file alongside the
binary, in `/etc/lever/config.yml`, or as environment variables.

```yaml
# Example YAML config; showing default values
releaseMode: true         # Set to false for debug logs & profiling tools
store: cql                # Set to "memory" to use an in-memory storage backend
http:
  addr: 8500              # The HTTP port to bind to
  cert: ""                # SSL-only
  key: ""                 # SSL-only
cassandra:                # Only used when store == "cql"
  keyspace: lever
  hosts:
  - 127.0.0.1
statsd:                   # Always enabled. The application will not fail to
                          # start if StatsD is unavailable
  addr: 127.0.0.1:8125
  bufferLength: 100       # The maximum num of stats to buffer before flushing
oauth:
  host:                   # The full root host of the OAuth2 provider
  user:                   # HTTP Basic Auth username
  password:               # HTTP Basic Auth password
```

Any of these items can be set via environment variables, prefixed with
`LEVER_`. Dictionaries are converted to underscores, for example:

* `LEVER_RELEASEMODE=false`
* `LEVER_HTTP_ADDR=80`

## API

You can find the full API definition here ([RAML](http://raml.org/)): [lever.raml](lever.raml)

## TODO

The following items still need to be completed before the application is to be
considered full-featured, ordered by priority. Some items are also flagged as
low-hanging fruit, which anyone with or without Go knowledge could knock out in
a few hours.

* H: More tests
* H: Increase configurations; a lot of hardcoded assumptions
* M: Server -> Cassandra caching (low-hanging fruit, in progress)
* M: hystrix-go integration
* M: Server Config model & endpoints (low-hanging fruit)
* L: Angular Web UI

## Credit

I derived the layout of this service on the patterns in [drone](https://github.com/drone/drone).
