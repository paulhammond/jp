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

Using [Homebrew](http://brew.sh/):

```
brew install https://gist.github.com/paulhammond/9441506/raw/jp.rb
```

If you don't use Homebrew you can download a
[precompiled binary](https://github.com/paulhammond/jp/releases) and copy the
`jp` file inside to somewhere in your path.

Or if you have a working go installation run
`go get github.com/paulhammond/jp`.
