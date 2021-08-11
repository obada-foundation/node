## Requirements

* Go 1.16
* cgo must be enabled

## Build Node from source
```bash
git clone git@github.com:obada-foundation/node
cd node/src
RUN CGO_ENABLED=1 go build --tags "json1 fts5 secure_delete" -a -o node main.go
```

which should generate a binary **node**, you can test it by running it with:

```bash
node run --url=your.node.url.com
```