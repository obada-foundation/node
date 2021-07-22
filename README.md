<div align="center">
  <img class="logo" src="link to logo" alt="OBADA | Blockchain Node"/>
</div>

The “blockchain node”. Connects the OBADA API to a synchronized server copy of the “real” node, currently hosted in Amazon QLDB. It’s middleware, designed to abstract blockchain protocols from the client and API. It is distributed as a single binary or as a docker container.

- Single binary distribution
- Docker container distribution
---

<<< badges here >>>

Starting node requires to have AWS credentials and SQS ARN to subscribe from OBADA tech team. The rest of parameters are strictly optional and have sane default.

Examples:

- with a static provider: `node --aws-key=AKIA435245 --aws-secret=435453453466634 --pubsub-topic-arn=arn:aws:sns:us-east-1:271164744603:yourorg.fifo`
- as a docker container `docker up -p 80:3000 obada/node:develop --aws-secret=435453453466634 --pubsub-topic-arn=arn:aws:sns:us-east-1:271164744603:yourorg.fifo`

## Install

OBADA Node distributed as a small self-contained binary as well as a docker image. Both binary and image support multiple architectures and multiple operating systems: linux_x86_64, windows_x86_64.

- for a binary distribution download the proper file in the [release section](https://github.com/obada-foundation/node/releases)
- docker container available on [Docker Hub](https://hub.docker.com/r/obada/node). I.e. `docker pull obada/node`.

## Default ports

In order to eliminate the need to pass custom params/environment, the default `--web-api-host` is dynamic and trying to be reasonable and helpful for the typical cases:

- If anything set by users to `--web-api-host` all the logic below ignored and host:port passed in and used directly.
- If nothing set by users to `--web-api-host` and node runs outside of the docker container, the default is `0.0.0.0:3000`.
- If nothing set by users to `--web-api-host` and node runs inside the docker, the default is `0.0.0.0:3000`.

## Options

Each option can be provided in two forms: command line or environment key:value pair. All options have the long form, i.e `--web-api-host=localhost:3000`. The environment key (name) [listed](#all-application-options) for each option as a suffix, i.e. `[$WEB]`.

## All Application Options

```
  --web-api-host/$WEB_API_HOST                  <string>    (default: 0.0.0.0:3000)
  --web-debug-host/$WEB_DEBUG_HOST              <string>    (default: 0.0.0.0:4000)
  --web-read-timeout/$WEB_READ_TIMEOUT          <duration>  (default: 5s)
  --web-write-timeout/$WEB_WRITE_TIMEOUT        <duration>  (default: 5s)
  --web-shutdown-timeout/$WEB_SHUTDOWN_TIMEOUT  <duration>  (default: 5s)
  --sql-sqlite-path/$SQL_SQLITE_PATH            <string>    (default: obada.db)
  --zipkin-reporter-uri/$ZIPKIN_REPORTER_URI    <string>    (default: http://zipkin:9411/api/v2/spans)
  --zipkin-service-name/$ZIPKIN_SERVICE_NAME    <string>    (default: obada-node)
  --zipkin-probability/$ZIPKIN_PROBABILITY      <float>     (default: 0.05)
  --aws-region/$AWS_REGION                      <string>    (default: us-east-1)
  --aws-key/$AWS_KEY                            <string>    (noprint)
  --aws-secret/$AWS_SECRET                      <string>    (noprint)
  --qldb-database/$QLDB_DATABASE                <string>    (default: obada)
  --pubsub-timeout/$PUBSUB_TIMEOUT              <duration>  (default: 5s)
  --pubsub-queue-url/$PUBSUB_QUEUE_URL          <string>    (default: https://sqs.us-east-1.amazonaws.com/271164744603/obada-tradeloop.fifo)
  --pubsub-topic-arn/$PUBSUB_TOPIC_ARN          <string>    (default: arn:aws:sns:us-east-1:271164744603:obada.fifo)
  --help/-h                                     
  display this help message
  --version/-v  
  display version information
```