This tool requires the root password to be set as an environment value:
```
$ export PASSWORD=<root password>
```

The tool is hard-coded to look for a configuration in `./config/desired.yml`

Usage:
./simple-sync

The tool will read the config and, in this order:
- stop services it manages (the ones defines in `./config/.known.yml`)
- copy/override files and apply metadata
- install packages
- restart the packages marked as services