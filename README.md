# Frack - An interval program generator!

A running interval training generated for different vVO2max (or MAS).

See informatbout about your MAS [here](http://www.scienceforsport.com/maximal-aerobic-speed-mas/) and a tool to calculate it [here](https://www.2peak.com/tools/mas.php)

The website in action is located here :

https://frack.chmouel.com

## REST Server

To install the REST Server you simply need to :

```shell
go get -u github.com/chmouel/chmoufrack
```

The REST server is a GO service connected to a MySQL database, see the file [etc/systemd.service](etc/systemd.service) for an example of how to run it. Make sure to adjust the path to where is your GOLANG path (in the example it's located in /usr/local/go/src

The database would be automatically propogated.

Some REST resources like POST/DELETE needs to have a facebook token to make sure we delete only the one we have created.

## Client

The client is an AngularJS webpage. You can serve it directly from your web server (recommended) or let the REST server doing it.

For example here is a Nginx snippet that pass all /v1 path to the rest server and the serve directly the HTML files :

```
    location /v1 {
        root /usr/local/go/src/github.com/chmouel/chmoufrack/client
        index index.html index.htm;

        proxy_read_timeout    90;
        proxy_connect_timeout 90;
        proxy_redirect        off;

        proxy_set_header      X-Real-IP $remote_addr;
        proxy_set_header      X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header      X-Forwarded-Proto https;
        proxy_set_header      X-Forwarded-Port 443;
        proxy_set_header      Host $host;

        proxy_pass http://localhost:9091/v1;
    }


    location / {
        root /usr/local/go/src/github.com/chmouel/chmoufrack/client;
        index index.html index.htm;
    }
```

This is assuming you are running the server on 9091 like in the [etc/systemd.service](etc/systemd.service) example.

You need to adjust the [config.js](client/js/config.js) file to your facebook application ID.

If you want to add a Program you need to be logged to Facebook and then you can click on the Menu to "add a program"


## Build

``frack`` use [https://github.com/Masterminds/glide](glide) for dependency management. After installing glide (see README) just issue this to build chmoufrack :

```
glide up
go build -o chmoufrack -v github.com/chmouel/chmoufrack/server/cli
```

To run the test you need to have a MySQL DB access specified in the environment variable `FRACK_TEST_DB`. You can setup access via the .my.cnf mechanism
