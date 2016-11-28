Frack -- Frack it those fractionee!
-----------------------------------

A tool to get the interval times for a multiple runners with different Vo2Max.

The name is a play on world between F***ck it (The usual word you use when
looking at the coach program) and fractionees with mean intervals in french.

USAGE
=====

This has been a work in progress but the basic yaml config source should work,
look over the samples/yaml-config.yaml file for an example and specify on frack
command line for example :

$ chmoufrack -y my-new-config-file.yaml 3x800/2x1000/1x2000 > /tmp/frackitlikeinthe90s.html

This will generate a static html page of your different workouts,

There is other options on the command for example to adjust the VMA ranges.


Requirements
============

- go get -u gopkg.in/yaml.v2
