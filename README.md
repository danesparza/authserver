# authserver [![CircleCI](https://circleci.com/gh/danesparza/authserver.svg?style=shield)](https://circleci.com/gh/danesparza/authserver)
OAuth 2 based token issue and validation server.  Batteries included: Uses TLS, and comes with its own admin UI to manage users (credentials), resources and roles.  Uses the embedded SQL database [QL](https://github.com/cznic/ql) on the backend.     

Why reimplement an OAuth authorization serice in your app if you don't have to?

## Quick start

* Make sure you have a TLS cert & key for your machine.  If you need one for local development / testing, use [mkcert](https://github.com/FiloSottile/mkcert).  
* Generate a configuration file using `authserver config create > authserver.yml`
* Update the config file with your specific settings.
* Bootstrap the system using `authserver bootstrap`.  This will create the admin password for your system and display it.  Please make a note of it -- you'll only see it once.
* Start the service and admin UI using `authserver start`


