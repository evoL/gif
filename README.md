# gif
[![Build Status](https://travis-ci.org/evoL/gif.svg?branch=master)](https://travis-ci.org/evoL/gif)

A command-line image library manager for nerds.

Its main use case is to manage a library with tags and URLs for fast sharing. Written in [Go](http://golang.org/).

[![Screencast](https://asciinema.org/a/25986.png)](https://asciinema.org/a/25986)

## Installation

### `go get`

If you have [Go](http://golang.org/) installed, you can use `go get` to install `gif` along with the sources:
```
go get github.com/evoL/gif
```

The code is being developed on Go 1.6.2.

### Binaries

Every stable release is available for download on [GitHub](https://github.com/evoL/gif/releases). Builds are prepared for following systems:

- Linux i386 / amd64
- OS X (Darwin) 64-bit
- Windows 32-bit / 64-bit

To install, [download the latest release](https://github.com/evoL/gif/releases) for your system and put it somewhere in your PATH.

## Usage

`gif` is composed of multiple commands, similar to `git`. If you run it without passing any commands, a help message is displayed.

You can prepend the command with the `--config`/`-c` option to specify a configuration file. By default `$HOME/.gifconfig` is used. Read more about configuration in the [`config` section](#config).

### `add`
**Adds an image.**

This command adds an image to the database. You can pass a local file or an URL as an argument. If you pass in a URL, the image will be downloaded, stored on your computer along with its URL. If you pass a local file, only it will be stored.

If you add a file using a URL and you already have it in your database, the URL will be updated. This is useful for updating broken URLs and adding them altogether for images that don't have them.

### `config`
**Prints the current configuration.**

`gif` is configured with a JSON file, kept by default in `$HOME/.gifconfig`. If you don't have one, that's OK — default values will be used. You can specify only part of the values to configure.

### `export`
**Exports the database.**

`gif` can export the database in two formats:

- JSON — Contains only tags and URLs. Only remote images are exported that way.
- bundle — Contains tags, URLs and image files.

The default is to export as JSON to stdout. You can specify a file to export to using the `--output` or `-o` option. If the file name ends with `.tar.gz` or `.gifb` (gif bundle) a bundle will be exported.

The `--bundle` option allows to export a bundle. Since by default it exports to stdout, you'll probably wat to use it in conjunction with the `--output`/`-o` option.

### `import`
**Imports multiple images into the database.**

To import images to the database, you can specify the following things as an argument:

- local file name — Should point to a JSON file or bundle. Allows for importing images and preserving tags and URLs.
- URL — Just as above, but the JSON/bundle will be downloaded.
- local directory — Will add every image in the directory to the database as local images. File names will be used as tags to allow recognition. You can use the `--recursive`/`-r` option to check the directory for images recursively.

### `list`
**Lists stored images.**

This command causes `gif` to list the content of its database. First, the number of images is displayed. Then, the following information about each image is displayed:

- first 8 characters of the ID (IDs are 40 characters long)
- whether it's a local or remote image
- size
- addition date
- comma-separated tags

The output can be filtered using what you pass as the argument. You can pass:

- an ID prefix, e.g. `0ab`
- a tag name, e.g. `awesome`

You can disable searching by ID with the `--tag`/`-t` option. Beside that, there are following options available:

- `--untagged` — Lists only images that have no tag.
- `--local` — Lists only local images.

### `path`
**Lists paths to images.**

This is the perfect option to get to the image file directly.

By default, one random image is chosen from the matching set. You can change it using the `--order`/`--sort`/`-s` option. Valid values are: `random`, `newest`, `oldest`.

To find a matching image you can use the same options as in [`list`](#list). In addition to that, you can get all matching images instead of just one using the `--all`/`-a` option.

### `remove`
**Removes images.**

To find a matching image you can use the same options as in [`list`](#list).

If there is more than one image that match your criteria, you can choose them from a list. You can also use `--all` to remove all matching images. This also works without specifying any filter.

If you want to bypass removal confirmation, use the `--really` option.

### `tag`
**Enables to change tags for images.**

To find a matching image you can use the same options as in [`list`](#list).

### `tags`
**Lists tags available in the database along with their image count.**

This command allows you to see what tags you used. First, the number of tags will be displayed. If there are any tags, the column headers are displayed along with the tags.

You can enter a tag prefix as an argument to filter the list. Example:
```
$ gif tags wtf
2 tags
TAG                IMAGE COUNT
wtf                3
wtf are you doing  1
```

### `url`
**Lists URLs of images.**

Works the same as [`path`](#path), just lists URLs for remote images.

### `upload`
**Uploads images to a server and saves the URLs for later use.**

To find matching images you can use the same options as in [`list`](#list).

Currently only uploading to [imgur](https://imgur.com/) is supported. To enable it, you have to [register an application](https://api.imgur.com/oauth2/addclient) (only anonymous usage is currently supported in `gif`) and set the Client ID in your [configuration file](#config). The reason for this is that the API supports up to 1,250 uploads a day per client (checked Sep 5 2015), so someone with a large library could use the limit.

If you want to just upload a couple of images though and don't want to register an application, contact me and I'll send you a Client ID for personal use.

Here's a minimal configuration file that sets up uploading:

```
{
  "Upload": {
    "Provider": "imgur",
    "Credentials": {
      "ClientId": "your client id"
    }
  }
}
```

## Examples

### Add a reaction GIF to the database

```
gif add http://www.reactiongifs.com/wp-content/uploads/2012/11/excited.gif
```

### Copy a random facepalm image URL to clipboard on OS X

```
gif url facepalm | pbcopy
```

### Preview all images on OS X

```
gif path --all | xargs qlmanage -p
```

### Backup your database to a bundle

```
gif export --bundle -o gif-backup.tar.gz
# or just
gif export -o gif-backup.tar.gz
```

### Import the database from a backup

```
gif import gif-backup.tar.gz
```

### Import a directory

```
gif import path/to/directory
```

## Footnote

Contributions are welcome! If you have something to add or you found a bug, file an issue or send a pull request.

License: [MIT](https://github.com/evoL/gif/blob/master/LICENSE)

Copyright © 2015–2016 Rafał Hirsz.
