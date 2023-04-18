spa-server
==========

Simple static file server for single-page application

This fork introduces a bunch of configuration and code cleanup. You can specify a custom listening hostname and port,
as well as custom base directory for the files. It is as simple as:

```sh
PORT=3000 HOST=127.0.0.1 BASE_DIRECTORY=/home/ubuntu/application spa-server
```

Please build the binary yourself.

```sh
# Assuming you already have Go
go build .
./spa-server
```

```sh
$ tree
.
|-- index.html
`-- assets
    |-- js
    |   `-- main.js
    `-- css
        `-- main.css
$ sap-server 5050
...
$ curl http://localhost:5050/
=> ./index.html
$ curl http://localhost:5050/assets/js/main.js
=> ./assets/js/main.js
$ curl http://localhost:5050/index.html
=> ./index.html
$ curl http://localhost:5050/page1
=> ./index.html
$ curl http://localhost:5050/page2/123
=> ./index.html
```
