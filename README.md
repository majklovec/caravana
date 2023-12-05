# Caravana

Deploy to Nomad using templates

This is proof of concept, could be subject to change.

## Repositories

### Add

* `caravana repo add https://github.com/majklovec/caravana-repo` - will clone services rom the given git repository

### Remove

* `caravana repo del https://github.com/majklovec/caravana-repo`
* `caravana repo del caravana-repo`

### Update all repos

* `caravana repo update`

## Jobs

### start

* `caravana start git.vondracek.dev`

There needs to be `config/git.vondracek.dev.yaml` with defined `SERVICE` , that has to be present in `services/<SERVICE>`

* `TEMPLATE=majkl/gitea DOMAIN=test caravana start`

### stop

* `caravana stop git.vondracek.dev`

### status

* `caravana status git.vondracek.dev`

## API

* `caravana api`

Than you can send HTTP POST `{"DOMAIN":"xxxx.domain.com", "TEMPLATE":"majkl/gitea"}` to http://localhost:8080/deploy to deploy the service

There is a html page at http://localhost:8080 for testing
