This tool requires the root password to be set as an environment value:
```
$ export PASSWORD=<root password>
```

The tool is hard-coded to look for a configuration in `./config/desired.yml`

Usage:
./simple-sync

Configuration:
This is how the `./config/desired.yml` looks like:

```yaml
---
servers:
- ip: 50.0.0.1
- ip: 50.0.0.2
apps:
- name: hello-world
  packages:
  - name: apache2
    is-service: true #if true, app starts a service with the same name
  - name: php5
  - name: libapache2-mod-php5
  files: # files are copied overriding existing files
  - path: /etc/apache2/mods-available/dir.conf
    mode: 644 # this value is passed as is to chmod
    owner: root # this value is passed as is to chown
    group: root # this value is passed as is to chown
    content: |
      <IfModule mod_dir.c>
        DirectoryIndex index.php index.html index.cgi index.pl index.xhtml index.htm
      </IfModule>
  - path: /var/www/html/index.php
    mode: 644
    owner: root
    group: root
    content: |
      <?php

      header("Content-Type: text/plain");

      echo "Hello, world!\n";
```

A **service** is a configured package with `is-service: true`

The tool will read the config and, in this order:
- stop services it manages (the ones defined in `./config/.known.yml`) if any
- copy/override files and apply metadata
- install packages
- restart the packages marked as services