# Specs

**Core**:

```
core:list
core:version
core:download
core:download -v 3.6
core:download --version 3.5.RC1
core:install
core:config
```

**Plugins:**

```
plugin:list
plugin:install PLUGIN_NAME
plugin:install PLUGIN_NAME --git GIT_URL
plugin:install PLUGIN_NAME --url FILE_URL
plugin:delete PLUGIN_NAME
```

**Themes:**

```
theme:list
theme:download THEME_NAME
```

**Database:**

```
db:import DB_NAME DB_FILE
db:export DB_NAME DB_FILE
```