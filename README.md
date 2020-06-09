CLI
=====
![](https://github.com/projecteru2/cli/workflows/goreleaser/badge.svg)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c4328f07835a43148ef8d2a87dbe5c85)](https://www.codacy.com/app/projecteru2/cli?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=projecteru2/cli&amp;utm_campaign=Badge_Grade)

CLI for Eru.

Modify resources for eru pods / nodes, manipulate containers and images.

### Usage

* Use `cli -h` to show commands and subcommands.
* Currently supported commands are:
	* `container` subcommands:
		* `container get {id} ...`
		* `container remove {id} ...`
		* `container realloc --cpu {cpu} --memory {memory}`
		* `container deploy`
	* `pod` subcommands:
		* `pod list`
		* `pod add`
		* `pod nodes {podname}`
		* `pod networks {podname}`
	* `node` subcommands:
		* `node get {nodename}`
		* `node remove {nodename}`
		* `node set-status [--available] {nodename}`
	* `image` subcommands:
		* `image build`
	* `lambda`

### Develop

We use `go mod` to manage requirements.  
Start developing:

```
$ git clone github.com/projecteru2/cli
$ make test binary
```

Commands' source code in `commands` dir, you can define your own commands inside. Use `make test` to test and `make build` to build. If you want to modify and build in local, you can use `make deps` to generate vendor dirs.

### Dockerized cli

Image: [projecteru2/cli](https://hub.docker.com/r/projecteru2/cli/)

```shell
docker run -it --rm \
  --net host \
  --name eru-cli \
  projecteru2/cli \
  /usr/bin/eru-cli <PARAMS>
```
