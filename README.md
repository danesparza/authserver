# authserver [![CircleCI](https://circleci.com/gh/danesparza/authserver.svg?style=shield)](https://circleci.com/gh/danesparza/authserver)
OAuth 2 based grant issue and validation server.  Batteries included: Uses TLS, and comes with its own UI to manage credentials, resources and roles.  Uses [BoltDB](https://github.com/coreos/bbolt) on the backend.     

Why reimplement OAuth in your app if you don't have to?

## Quick start

* Make sure you have a TLS cert & key for your machine.  If you need one for local development / testing, use [tls-keygen](https://www.npmjs.com/package/tls-keygen).  
* Generate a configuration file using `authserver config create`
* Update the config file with your specific settings.
* Start the service and admin UI using `authserver start`

## Data model
![Data model](https://github.com/danesparza/authserver/raw/master/data_model.svg?sanitize=true)
