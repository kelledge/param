## Purpose
Allows dereferencing of specific versions of parameters stored in AWS SSM. Could
potentially be extended to solutions like Hashicorp Vault.

## Scope
 1. Dereference versioned parameters in SSM
 2. Manage deployment of parameters to SSM
 3. Introspection of all parameter value verions

## NFR
 1. Simple addition to docker images

## Workflow
```
+------------+                +---------------------+              +-----------------+
| Repository | >--[deploy]--> | SSM Parameter Store | >--[exec]--> | Process Environ |
+------------+                +---------------------+              +-----------------+
```

## Intended Usage
Intended usage (docker-compose):

Assuming the following parameters are stored in SSM:

| Path                                  | Version | Value     |
| ------------------------------------- |:-------:| ---------:|
| /dev/front-end-api/DATABASE_PASSWORD  | 1       | oldpass   |
| /dev/front-end-api/DATABASE_PASSWORD  | 2       | insecure  |
| /dev/front-end-api/DATABASE_PASSWORD  | 3       | pa55w0rd  |
| /dev/front-end-api/TWILIO_API_KEY     | 1       | XXXXXX    |
| /dev/front-end-api/TWILIO_API_SECRET  | 1       | YYYYYY    |
| /dev/front-end-api/TWILIO_ACCOUNT_SID | 1       | ZZZZZZ    |


And given the following service definition:
```
version: '3'
services:
  print:
    image: busybox
    entrypoint: /usr/bin/local/param exec
    command: /usr/bin/env
    environment:
      - "SSM_PARAM_PATH=/dev/front-end-api"
      - "SSM_PARAM_DATABASE_PASSWORD=3"
      - "SSM_PARAM_TWILIO_API_KEY=1"
      - "SSM_PARAM_TWILIO_API_SECRET=1"
      - "SSM_PARAM_TWILIO_ACCOUNT_SID=1"
```

`/usr/bin/env` would report the following to stdout:
```
...
DATABASE_PASSWORD=pa55w0rd
TWILIO_API_KEY=XXXXXX
TWILIO_API_SECRET=YYYYYY
TWILIO_ACCOUNT_SID=ZZZZZZ
...
```

## Deployment Usage
TODO

Goals:
 * Encrypted at rest using KMS
 * Leverage VCS
 * Declarative

## State
Proof-of-concept
