# Description

SSHFP Tool is a tool created in Golang to glue Consul DB and Cloudflare.

Main purpose is creating SSHFP records to get rid of "host key verification failed".

# Building

```sh
make
```
Or to build for all platforms
```sh
make release
```

# Usage

Supported env variables:

* `DOMAIN_NAME` - Domain name which will be working on.
* `CF_TOKEN` - CloudFlare Token with write access to above domain.
* `HOST_LIVENESS_TIMEOUT` - After this number of seconds of host being offline DNS records are removed.

It's possible to create JSON formatted config file (example in `testcfg`)

As it has been designed to work with `consul watches` passing proper JSON file to `stdin` is required.
```sh
cat watches.dump | ./sshfp-updater
```

# Current state

- CloudFlare integration is fully implemented
- SSHFP Record creation based on tag in Consul form.
- Implemented Consul watches integration
- Implemented logic to manipulate states (merging config, etc)

# TODO

- A few major changes
