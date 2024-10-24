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

It only adds and removes whitespace, which means that your data won’t get
silently altered. For example, `"\u2603"` won’t get converted to `"☃"`, and
`1.1e1` won’t turn into `11`. The ordering remains the same, and invalid JSON
can be reformatted (within reason). This stuff shouldn’t matter, but people make
mistakes even with a well defined format like JSON, and accurate tools are
important when you’re trying to work out what’s gone wrong.

For more information see the [project homepage][jp].

[jp]: https://paulhammond.org/jp

## Installing

Using [Homebrew](http://brew.sh/):

    brew install paulhammond/tap/jp

If you don't use Homebrew you can download a [precompiled binary][releases] and
copy the `jp` file inside to somewhere in your path.

[releases]: https://github.com/paulhammond/jp/releases

## Using

To prettify a JSON file:

    jp data.json

To prettify from stdin, use - as the filename:

    curl -sL https://phmmnd.me/names.json | jp -

To compact a JSON file:

    jp --compact data.json

To get help:

    jp --help
