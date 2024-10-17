# jp

`jp` is a command line tool that reformats JSON to make it easier to read:

    $ cat names.json
    {"names":["Alice","Bob"]}
    $ jp names.json
    {
      "names": [
        "Alice",
        "Bob"
      ]
    }

It's fast, doesn't mess with the data, and handles invalid JSON (within
reason). For more information see the [project
homepage][jp].

## Installing

Using [Homebrew](http://brew.sh/):

    brew install paulhammond/tap/jp

If you don't use Homebrew you can download a [precompiled binary][releases] and
copy the `jp` file inside to somewhere in your path.

[jp]: https://paulhammond.org/jp
[releases]: https://github.com/paulhammond/jp/releases
