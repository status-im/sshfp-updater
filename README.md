# Description

SSHFP Tool is a tool created in Golang to glue Consul DB and Cloudflare.

Main purpose is creating SSHFP records to get rid of "host key verification failed".

## Building

```
go mod -vendor
go build -mod vendor
```

## Usage

Supported env variables:
`DOMAIN_NAME` - Domain name which will be working on
`CF_TOKEN` - CloudFlare Token with write access to above domain
`HOST_LIVENESS_TIMEOUT` - number in seconds after which host is 
considered as removed and dns records can be deleted

It's possible to create json formatted config file (example in `testcfg`)

As it has been designed to work with `consul watches` passing proper .json file
to STDIN is required. Ex:
`cat watches.dump | ./infra-sshfp-cf`

## Current state
- CloudFlare integration is fully implemented
- SSHFP Record creation based on tag in Consul form.
- Implemented Consul watches integration
- Implemented logic to manipulate states (merging config, etc)

## TODO:
- A few major changes
