## Node minimal hardware requirements

### With docker
* 2 x CPU cores
* 2GB RAM
* 20GB disk space, we recommend SSD or NVMe

### Without docker 
* 1 x CPU cores
* 2GB RAM
* 20GB disk space, we recommend SSD or NVMe

## Before installation
To run a node you need to have two things:
* You should have access to [Node GitHub repo](https://github.com/obada-foundation/node)
* You should have AWS secrets for QLDB and SQS

## Linux binary installation
1. Download a [binary](https://github.com/obada-foundation/node/releases/download/v0.0.3/node.v0.0.3.linux-amd64.tar.gz)
2. Upload it on your server where you want to run a Node
3. Un-archive
```bash
tar -xf node.0.0.3.linux-amd64.tar.gz
```
4. Move executable to **/usr/bin**
```bash
mv node.linux-amd64 /usr/bin
```
5. Test that binary was discovered by **$PATH** and works
```bash
node -h
```
should print:
```bash
OBADA-NODE v0.0.3
Usage:
  node [OPTIONS] <init | run>

Application Options:
      --url=  url to OBADA API Node [$NODE_URL]

Help Options:
  -h, --help  Show this help message

Available commands:
  init
  run
```
6. Initialize node
```bash
DB_PATH=~/obada/db
DB_FILE=$DB_PATH/obada.db
mkdir -p $DB_PATH
touch $DB_FILE
node init --sql.path=$DB_FILE --url=http://localhost
```
where **~/obada/db/obada.db** is path to sqlite database file

After completing first 6 steps and with no errors you are ready to run a Node.

### Run Node on **localhost** interface
```bash
node run \
  --aws.key=AKIA435245 \
  --aws.secret=435453453466634 \
  --pubsub.topic-arn=https://sqs.us-east-1.amazonaws.com/243546434403/your-node.fifo \
  --url=http://localhost \
  --sql.path=~/obada/db/obada.db
```
### Run public Node with LetsEncrypt enabled
Before executing this command, please make sure that you have a DNS record properly configured. 
To test run:
```bash
dig your.node.com a
```
the IP address in response should match your server public IP.
Then you can run a Node:

```bash
node run \
  --api.address=0.0.0.0 \
  --api.port=80 \ 
  --ssl.type=auto \ 
  --ssl.port=443 \
  --url=https://your.node.com \
  --aws.key=AKIA435245 \
  --aws.secret=435453453466634 \
  --pubsub.queue-url=https://sqs.us-east-1.amazonaws.com/243546434403/your-node.fifo \
  --sql.path=~/obada/db/obada.db
```

Test Node with **curl**:
```bash
curl https://your.node.com/ping
```
Should respond back **pong** to you. To test SSL you can use SSL Labs: https://www.ssllabs.com/ssltest/analyze.html?d=your.node.com

### Systemd supervisor

Most Linux distributions use **systemd** as a supervisor to control application.
You can use configuration below as a template for you setup.
```
[Unit]
Description=OBADA Node service
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=root
ExecStart=/usr/bin/node run --aws.key=AKIA435245 --aws.secret=435453453466634 --pubsub.topic-arn=https://sqs.us-east-1.amazonaws.com/243546434403/your-node.fifo --url=http://localhost --sql.path=~/obada/db/obada.db

[Install]
WantedBy=multi-user.target
```

## Windows binary installation
...