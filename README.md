# wptool

Experimental Wordpress CLI written in Go.

## Development

Go v1.1.1 or newer is required. 
You can use something like `gvm` to manage different Go languages.

Make sure you have `GOPATH` set. 

Clone repository and install dependencies:

```
git clone https://github.com/doejo/wptool.git
cd wptool
go get
```

Build:

```
go build wptool.go
```

## Usage

```
wptool COMMAND [ARGS]
```

List of available commands:

- `core:list`
- `core:version`
- `core:download`
- `core:config`
- `core:install`

### core:list

Returns a list of all available wordpress core versions. List of versions is 
maintained as a plain text file and generated from SVN tags of wordpress
core repository. Updated manually.

### core:version

Prints installed wordpress core version. Does not invoke any PHP code at all, 
just reads contents of `wp-includes/version.php` and parses wordpress version. 
If path is not specified, it will use current path.

Arguments:

- `-p`, `--path PATH` - Specify path to wordpress core dir

### core:download

Downloads and extracts wordpress core. By default will download latest available
core. Also possible to specify beta releases. In case when upgrade is required,
provide `-f` flag and existing core will be replaced with a new one.
Configuration is required after core is replaced.

Arguments:

- `-v`, `--version VERSION` - Specify version to download
- `-p`, `--path PATH` - Specify path to extract
- `-f`, `--force` - Override existing core

### core:config

Configure wordpress core. Generates a new `wp-config.php` file under wordpress
core dir. 

Arguments:

- `-p`, `--path PATH` - Path to wordpress core
- `-t`, `--template PATH` - Path to config template
- `-f`, `--force` - Force config override
- `--dbname` - Name of database
- `--dbhost` - Database connection host (Default: "localhost")
- `--dbuser` - Database connection user
- `--dbpass` - Database connection password
- `--dbcharset` - Database charset (Default: "utf8")
- `--dbcollate` - Database collate
- `--dbprefix` - Tables prefix (Default: "wp_")

## License

MIT License

Copyright (c) 2013 Dan Sosedoff

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
