# Caravana

Deploy to Nomad using templates

## Quick start

### Download caravana 

from https://github.com/majklovec/caravana/releases

 `sudo cp caravana* /usr/local/bin/caravana`

 `sudo chmod +x /usr/local/bin/caravana`

### Create your config directory

 `mkdir configs`

### Add sample repository:

 `caravana repo add https://github.com/majklovec/caravana-amd64`

 `caravana repo ls`

### Create config files:

* cd configs
* add config files named <domain>.yaml

## Repositories

### Add

* `caravana repo add https://github.com/majklovec/caravana-amd64` - will clone services rom the given git repository
Name of the repo will not include `caravana-` if present. So repo above will end up in templates/amd64

### Remove

* `caravana repo del https://github.com/majklovec/caravana-amd64`
* `caravana repo del amd64`

### Update all repos

* `caravana repo update`

## Jobs

### start

* `caravana start git.domain.com`

There needs to be `config/git.domain.com.yaml` with defined `SERVICE` , that has to be present in `services/<SERVICE>`

* `TEMPLATE=majkl/gitea DOMAIN=test caravana start`

### stop

* `caravana stop git.domain.com`

### status

* `caravana status git.domain.com`

## API

* `caravana api`

Will create http server on port 8080, and will listen for HTTP POST to `/deploy`

Payload is in JSON format: `{"DOMAIN":"xxxx.domain.com", "TEMPLATE":"majkl/gitea"}`

There is a html page at http://localhost:8080 for testing
