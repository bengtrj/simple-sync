Hello!

Let me start by saying that I had a lot of fun coding this tool.
It was challenging, interesting and I learned a lot. Thanks!

I'm happy to explain how I decided to tackle this problem and what where my priorities and trade-offs.

Let's see if I can give you all the right signals you're looking for in a good fit for this position!

**Simple configuration management tool:**
- Problem states that the tool should be rudimentary
- Problem suggests spending ~4 hours including the other exercise
- asks for a way to install/remove packages
- asks for a way of defining files and file metadata
- asks for idempotency

To me, this meant keeping it simple, focusing on getting configuration/packages/files in sync in a repeatable way, on a list of remote hosts.
If you run the tool repeatedly with the same configuration, it doesn't change the effective state of the system. 

To make things simple, the scope is limited to 2 hosts with the same `root` password.

I opted to go with a simple strategy:

```
->| RUN | -> | new state ? | - y -> | For each *server ip* configured:  |
                    |               | copy/override files + metadata    |
                    n               | install packages                  |
                    |               | starts services                   |
                    |             
                    ↓
    | For each *server ip* configured:                             |
    | stops known services                                         |
    | copy/override files + metadata                               |
    | deletes unnecessary files                                    |
    | installs missing packages                                    |
    | removes unnecessary packages                                 |
    | starts services                                              |
```

**Where:**
- *desired state* is configured with a simple `yaml` file
- *known state* of the system is abstracted as a file `.known` in the same directory as the apps.
- *new state* is the absence of a file `.known`
- if all the hosts are updated successfully, the `.known` state file is updated
- the tool assumes healthy servers, with adequate available disk space and access to installing packages using `apt-get`

**Limitations:**
- very simple error handling
- only install packages via `apt-get`
- assumes that deleting/overriding/uninstalling files/packages is ok
- assumes the files are text files that you can provide the content as configuration

**My choices:**
- Convention over configuration
  - the config file is hard-coded as `./config/desired.yml`
  - the tool automatically watches the config file for changes, responding accordingly
- Language of choice - Go:
  - it's easy to distribute in a self-contained binary
  - I found nice `ssh` and `file watcher` libraries
  - I'm familiar with Go

**How features are implemented:**

The tool iterates over the servers and for each one applies the configuration. Here's an example configuration that works for the proposed exercise:
```yaml
---
servers:
- ip: 50.0.0.1
- ip: 50.0.0.2
apps:
- name: hello-world
  packages:
  - name: apache2
    is-service: true
  - name: php5
  - name: libapache2-mod-php5
  files:
  - path: /etc/apache2/mods-available/dir.conf
    mode: 644
    owner: root
    group: root
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

 The tool uses the concept of an App.
 The `App` is formed by zero or more packages and zero or more files
 - To install a package, add one under `apps` as above's **php5** example
 - To copy a file, add it under `apps` as above **/var/www/html/index.php** example
 - If a package should be managed as a service, add the `is-service: true` option as above's **apache2** example. This will make the tool stop/start a service with the same name as the package, as required.
