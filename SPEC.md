# Specs

List of commands and options

## Core

Core coomands:

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

Flags:

- `-p`, `--path PATH` - Specify path to wordpress core dir

### core:download

Downloads and extracts wordpress core. By default will download latest available core.
Also possible to specify beta releases. In case when upgrade is required,
provide `-f` flag and existing core will be replaced with a new one. Configuration
is required after core is replaced.

Flags:

- `-v`, `--version VERSION` - Specify version to download
- `-p`, `--path PATH` - Specify path to extract
- `-f`, `--force` - Override existing core

### core:config

Configure wordpress core. Not implemented yet.

Flags:

- `--dbname` - Name of database
- `--dbhost` - Database connection host (Default: localhost)
- `--dbuser` - Database connection user
- `--dbpass` - Database connection password

## Plugins

```
plugin:list
plugin:install PLUGIN_NAME
plugin:install PLUGIN_NAME --git GIT_URL
plugin:install PLUGIN_NAME --url FILE_URL
plugin:delete PLUGIN_NAME
```

## Themes

```
theme:list
theme:download THEME_NAME
```

## Database

```
db:import DB_NAME DB_FILE
db:export DB_NAME DB_FILE
```