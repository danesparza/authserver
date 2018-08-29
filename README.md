# authserver [![CircleCI](https://circleci.com/gh/danesparza/authserver.svg?style=shield)](https://circleci.com/gh/danesparza/authserver)
OAuth 2 based token issue and validation server.  Batteries included: Uses TLS, and comes with its own admin UI to manage users (credentials), resources and roles.  Uses the embedded SQL database [QL](https://github.com/cznic/ql) on the backend.     

Why reimplement an OAuth authorization serice in your app if you don't have to?

## Quick start

* Make sure you have a TLS cert & key for your machine.  If you need one for local development / testing, use [mkcert](https://github.com/FiloSottile/mkcert).  
* Generate a configuration file using `authserver config create > authserver.yml`
* Update the config file with your specific settings.
* Bootstrap the system using `authserver bootstrap`.  This will create the admin password for your system and display it.  Please make a note of it -- you'll only see it once.
* Start the service and admin UI using `authserver start`

## Interacting with the service

First get a token for the admin user:
```
curl -X POST \
  https://localhost:3001/token/client \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json'
  -d '{
	"grant_type": "client_credentials",
	"client_id": "admin",
	"client_secret": "your_admin_password",
	"scope": "*"
}'
```

Next, create a new resource and a role:

Then create users and associate them to the new resource and role:

Finally, verify the user has been assigned to the new resource and role.
