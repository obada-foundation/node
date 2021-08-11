<div align="center">
  <img class="logo" src="https://www.obada.io/assets/images/obada-logo.svg" alt="OBADA | Blockchain Node"/>
</div>

The “blockchain node”. Connects the OBADA API to a synchronized server copy of the “real” node, currently hosted in Amazon QLDB. It’s middleware, designed to abstract blockchain protocols from the client and API. It is distributed as a single binary or as a docker container.

- Single binary distribution
- Docker container distribution
- Automatic SSL termination with Let's Encrypt
- Support of SSL certificates
---

<<< badges here >>>
Starting node requires to have a Node url, AWS credentials and SQS ARN to subscribe from OBADA tech team. The rest of parameters are strictly optional and have sane default.

Examples:

- as a binary: `node run --aws.key=AKIA435245 --aws.secret=435453453466634 --pubsub.topic-arn=arn:aws:sns:us-east-1:271164744603:yourorg.fifo --url=http://localhost`
- as a binary with automatic SSL `node run --aws.key=AKIA435245 --aws.secret=435453453466634 --pubsub.topic-arn=arn:aws:sns:us-east-1:271164744603:yourorg.fifo --url=https://yournode.com --ssl.type=auto --ssl.acme-email=youradminemail@company.com`
- as a docker container `docker up -p 80:3000 obada/node:develop --aws.key=AKIA435245 --aws.secret=435453453466634 --pubsub.topic-arn=arn:aws:sns:us-east-1:271164744603:yourorg.fifo --url=http://localhost`

## Install

OBADA Node distributed as a small self-contained binary as well as a docker image. Both binary and image support multiple architectures and multiple operating systems: linux_x86_64, windows_x86_64.

- for a binary distribution download the proper file in the [release section](https://github.com/obada-foundation/node/releases)
- docker container available on [Docker Hub](https://hub.docker.com/r/obada/node). I.e. `docker pull obada/node`.

See more details about installation options in [INSTALL.md](https://github.com/obada-foundation/node/blob/master/INSTALL.md)

## Default ports

In order to eliminate the need to pass custom params/environment, the default `--api.address` is dynamic and trying to be reasonable and helpful for the typical cases:

- If anything set by users to `--api.address` all the logic below ignored and host:port passed in and used directly.
- If nothing set by users to `--api.address` and node runs outside of the docker container, the default is `127.0.0.1:3000`.
- If nothing set by users to `--api.address` and node runs inside the docker, the default is `0.0.0.0:3000`.

## Options

Each option can be provided in two forms: command line or environment key:value pair. All options have the long form, i.e `--web-api-host=localhost:3000`. The environment key (name) [listed](#all-application-options) for each option as a suffix, i.e. `[$WEB]`.

## **node run** Application Options

```
Usage:
  node [OPTIONS] run [run-OPTIONS]

Application Options:
      --url=                            url to OBADA API Node [$NODE_URL]

Help Options:
  -h, --help                            Show this help message

[run command options]

    api:
          --api.port=                   port (default: 3000)
          --api.address=                listening address (default: 127.0.0.1)

    aws:
          --aws.region=                 AWS region (default: us-east-1)
          --aws.key=                    AWS credential key (default: us-east-1)
          --aws.secret=                 AWS credential secret (default: us-east-1)

    qldb:
          --qldb.database=              QLDB database name (default: obada)

    ssl:
          --ssl.type=[none|static|auto] ssl (auto) support (default: none) [$SSL_TYPE]
          --ssl.port=                   port number for https server (default: 3443) [$SSL_PORT]
          --ssl.cert=                   path to cert.pem file [$SSL_CERT]
          --ssl.key=                    path to key.pem file [$SSL_KEY]
          --ssl.acme-location=          dir where certificates will be stored by autocert manager (default:
                                        ./var/acme) [$SSL_ACME_LOCATION]
          --ssl.acme-email=             admin email for certificate notifications (default: techops@obada.io)
                                        [$SSL_ACME_EMAIL]

    sql:
          --sql.path=                   path to SQLite database (default: obada.db)

    pubsub:
          --pubsub.timeout=             pubsub timeout (default: 5s)
          --pubsub.queue-url=
          --pubsub.topic-arn=

```

## Integration

OBADA provides integration solutions for PHP, Python and C#. You also can check simple applications that show how to use these solutions.
All integration solutions are generated from [OpenAPI specification](https://github.com/obada-foundation/node/tree/master/openapi).

### PHP

- [Integration solution](https://github.com/obada-foundation/node-api-library)
- [Example application](https://github.com/obada-foundation/example-client-system)

### Python

- [Integration solution](https://github.com/obada-foundation/node-api-library-python)
- [Example application](https://github.com/obada-foundation/integration-scenarios/tree/master/python/simple-application)

### C#

- [Integration solution](https://github.com/obada-foundation/node-api-library-csharp)
- [Example application](https://github.com/obada-foundation/integration-scenarios/tree/master/csharp/simple-application/SimpleApplication)
