# jp

`jp` is a command line tool that reformats JSON to make it easier to read:

    $ cat names.json
    {"names":["Paul","Amy"]}
    $ jp names.json
    {
      "names": [
        "Paul",
        "Amy"
      ]
    }

It is insanely fast, doesn't mess with the data, and handles invalid JSON
(within reason). For more information see the [project
homepage](http://www.paulhammond.org/jp/).

## Installing

Download a [precompiled binary](http://www.paulhammond.org/jp/) and
copy the `jp` file inside to somewhere in your path.

Or if you have a working go installation run
`go get github.com/paulhammond/jp`.