### Frack -- Frack it those fractionee! ###

A tool to get the interval times for a multiple runners with different Vo2Max.

The name is a play on world between F***ck it (The usual word you spit when
looking for the first at tonight's coach program) and fractioneés means intervals
in french. It supposed to be funny but yet clever (fail)

There is two component to this, the REST server and the UI based on angular.

It was originally based on sqlite but I changed as yaml since that was good
enough and didn't want to bother.

USAGE
-----

Look over the `frack.yaml` for the example of format of the config file for how to
specify the list of the workouts.

You will then launch the UI :

```
$ go build
$ ./frack
```

You can then access to :

http://localhost:8080

to see the generated UI.

You can update the `frack.yaml` as much as you like it get read on every rest call


Requirements
------------

- go get -u gopkg.in/yaml.v2
