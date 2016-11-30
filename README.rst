Frack -- Frack it those fractionee!
-----------------------------------

A tool to get the interval times for a multiple runners with different Vo2Max.

The name is a play on world between F***ck it (The usual word you spit when
looking for the first at tonight's coach program) and fractioneés means intervals
in french. It supposed to be funny but yet clever (fail)

USAGE
=====

This has been a work in progress but the basic yaml config source should work,
look over the samples/yaml-config.yaml file for an example and specify on frack
command line for example :

$ chmoufrack -y my-new-config-file.yaml 3x800/2x1000/1x2000 > /tmp/frackitlikeinthe90s.html

This will generate a static html page of your different workout looking something like this
on desktop :



There is other options on the command for example to adjust the VMA ranges.


UI
==

The UI is in a heavy work in progress, you launch it with :

chmoufrack -rest

and it will expose itself to localhost:8080, just access directly in :

https://localhost:8080/static/html

as mentioned things are moving quite heavily in there.

Requirements
============

- go get -u gopkg.in/yaml.v2
