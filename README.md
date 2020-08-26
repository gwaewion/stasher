install:

prerequisites:
- go >= 1.14
- git

cd to GOPATH/src
git clone https://github.com/gwaewion/stasher.git stasher
cd stasher && git checkout dev
go get -u github.com/gobuffalo/packr/v2/packr2
GOPATH/bin/packr2
go get && go build

run:

create config file in yaml format (for example, config.yml):
```
stasher:
  address: "127.0.0.1"
  port: "8080"
  salt: "longAndSecureSalt"
  hostname: "127.0.0.1:8080"
  scheme: "http"
couchdb:
  protocol: "https"
  address: "localhost"
  port: "6984"
  dbname: "stasher"
  username: "stasher"
  password: "stasher"
  cert_check: false
```
./stasher -c config.yml
