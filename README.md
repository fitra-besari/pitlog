[![N|Solid](https://cldup.com/dTxpPi9lDf.thumb.png)](https://nodesource.com/products/nsolid)

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

#### Testing Machine:
    - MacOs 12.6.7
    - go version go1.23.1 darwin/amd64 
#  Pitlog: A Logging Middleware for Echo Framework

`Pitlog` is a logging middleware designed to integrate with the Echo framework in Go. It provides structured logging for both requests and responses with support for masking sensitive fields.

##  Installation
To install the package, use `go get`:
```
go get github.com/fitra-besari/pitlog
```
for example you can see on directory /example and you can run ` go mod download ` and ` go run main.go`, will be run on port `8080` you can hit `localhost:8080/ping` and you have see log api on your terminal and you config file on root repository on your project.

### parameters
you can config application name on log

`appName  :=  "ExampleApp"`

you can config application version on log

`appVersion  :=  "1.0.0"`

you can config application level on log, if you use `production` you will be get minimal indentation on log 

`appLevel  :=  "development"`

you can config directory of log file , will be on dir  `./logs` from root repository

`logDir  :=  "logs"`

you can config view log on your terminal console for enable log on console ,default will be outo enable on log file only . optional value is `"1"` or `"0"`

`enableLogConsole  :=  "1"`

you can config border of log view . optional value is `"1"` or `"0"`
`useSeparate  :=  "1"`

if you set  `"1"` you will be see your log data on object view, if use `"0"` will be on string view

`objectView  :=  "1"`

#### Note

- jika ada kendala di invite project github / postman mohon hubungi afitra - 085230010042

### http://apitoong.com
 