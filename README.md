## Jabba

This is a super-simple configuration management tool. It's like salt, chef, ansible, puppet etc only way way way smaller and simpler. It doesn't do nearly as many things as those tools. That's a feature.

This will:
* Never work on Windows
* Probably never work on BSD
* Probably never work on anything other than Debian-derived Linux distributions
* Never have a DSL

This is designed for a world in which all of your complex application configuration is handled by Docker containers, and all you need out of your config management tool is to run a couple commands, set up your users and put the odd config file in place. Maybe install the odd database here and there. In Salt, installing a database would look something like:

```
setup_repo:
  pkgrepo.managed:
    - humanname: RethinkDB Repo
    - name: deb http://download.rethinkdb.com/apt {{ grains['oscodename'] }} main
    - file: /etc/apt/sources.list.d/rethinkdb.list
    - key_url: http://download.rethinkdb.com/apt/pubkey.gpg
    - require_in:
      - pkg: rethinkdb

rethinkdb:
  pkg:
    - installed

/etc/rethinkdb/instances.d/default.conf:
  file.managed:
    - source: salt://files/etc/rethinkdb/instances.d/default.conf
    - user: root
    - group: root
    - mode: 644
```

In this it's:
```
jabba.RunOrDie("source", "/etc/lsb-release", "&&", "echo", ""deb", "http://download.rethinkdb.com/apt", "$DISTRIB_CODENAME", "main"", "|", "sudo", "tee", "/etc/apt/sources.list.d/rethinkdb.list")
jabba.RunOrDie("wget", "-qO-", "https://download.rethinkdb.com/apt/pubkey.gpg", "|", "sudo", "apt-key", "add", "-")
jabba.RunOrDie("sudo", "apt-get", "update")
jabba.RunOrDie("sudo", "apt-get", "install", "rethinkdb")
jabba.WriteFile(rethinkConfigFile)
```

The Salt one is smart because it's portable between all kinds of different operating system flavours. The second one is simpler because it's just a damn copy+paste out of the Rethink install docs rather than spelunking through Salt documentation for an hour looking for the right incantations.

All configurations are defined procedurally in code rather than declaratively. There are a few helper functions to assist with that in the `blobby` package. These will be factored out into their own project as soon as the API is considered stable.

All files, data etc is configured in Go templates. This is nice because:
1. Go templates are pretty powerful
1. It results in them getting compiled into the built binary

The nice thing about #2 is that once this is built into an executable, it can be shipped on its own without the need for any additional source or configuration files. Makes it much easier to SCP about the place.
