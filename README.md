# authserver [![CircleCI](https://circleci.com/gh/danesparza/authserver.svg?style=shield)](https://circleci.com/gh/danesparza/authserver)
OAuth 2 based grant issue and validation server

## Quick start

* Make sure you have a TLS cert & key for your machine.  If you need one for local development / testing, use [tls-keygen](https://www.npmjs.com/package/tls-keygen).  
* Generate a configuration file using `authserver config create`
* Update the config file with your specific settings.
* Start the service and admin UI using `authserver start`
