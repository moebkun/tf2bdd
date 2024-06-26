# TF2BDd

Very simple service to allow tools like those listed below to download and integrate player list contributions from outside sources. This
is designed to work over discord as a bot, allowing multiple people to contribute their lists and have them merged into
a single master list. The results served over an HTTP endpoint `/v1/steamids`. Data is backed by a very simple sqlite database.

- [tf2_bot_detector](https://github.com/PazerOP/tf2_bot_detector)
- [bd](https://github.com/leighmacdonald/bd)
- [MAC](https://github.com/MegaAntiCheat)

If you have other examples of software that is able to update lists over http like these, please open a PR to add them to the list.

Example results from the [@trusted](https://trusted.roto.lol/v1/steamids) list.

## Commands

Bot command list:

- `!add <steamid/profile> [attributes]` Add the user to the master ban list. eg: `suspicious/cheater/bot`. If none are defined, it will use cheater by default.
- `!del <steamid/profile>` Remove the player from the master list
- `!check <steamid/profile>` Checks if the user exists in the database
- `!count` Shows the current count of players tracked
- `!import <attached_playerlist_files>` Imports the steam ids from a players custom ban list, multiple can be attached
- `!steamid <steamid/vanity_name/profile_link>` Accepts any steamid format including bare vanity name and profile link. Will print out all forms.

Discord [slash commands](https://support.discord.com/hc/en-us/articles/1500000368501-Slash-Commands-FAQ) are not 
currently supported as this was written before that was an option, however if there is enough
demand, or somebody creates a PR for it, I will add them.

## Building From Source

    $ git clone git@github.com:leighmacdonald/tf2bdd.git
    $ cd tf2bdd
    $ go build

## Configuration

There is an example config located inside the releases `tf2bdd_example.yml`. Rename it to `tf2bdd.yml` and edit the 
values as documented inside of it.

Make sure you enable "Message Content Intent" on your discord config under the Bot settings via discord website. If your
bot does not respond to your commands, this is probably why.

## Running Binary

You can either use the binary you build from source, or download the latest release from the [releases](https://github.com/leighmacdonald/tf2bdd/releases)
page.

    $ cp tf2bdd_example.yml tf2bdd.yml     # Copy example config
    $ vim tf2bdd.yml                       # Edit/Add your config options
    $ ./tf2bdd

You will probably want to create something like a systemd service to automate this.

## Running Docker

Running over docker is generally the recommended approach, along with a reverse proxy such 
as [caddy](https://caddyserver.com/) that can provide automatic TLS certs for HTTPS support. 
If you are using an internal docker network (recommended), ensure you also add the `--network your_network` flag
to the run command. Take note that the container binds only to localhost in the example command shown `-p 127.0.0.1:8899:8899`, so you
will not be able to access it remotely unless you remove the `127.0.0.1` or add a reverse proxy in front of it. When
using a reverse proxy, ensure that you set the `external_url` config option to the url that people can read your server
at.

You can also use the `latest` image tag if you do not care about pinning to a specific version: `ghcr.io/leighmacdonald/tf2bdd:latest`.
It will always be using the latest release tag.

    docker run --name tf2bdd -d -it --restart unless-stopped \
        --pull always \
        -p 127.0.0.1:8899:8899 \
        --mount type=bind,source="$(pwd)"/db.sqlite,target=/app/db.sqlite \
        --mount type=bind,source="$(pwd)"/tf2bdd.yml,target=/app/tf2bdd.yml \
        ghcr.io/leighmacdonald/tf2bdd:v1.0.2

Make sure that when running under docker, you do set `listen_host: ""` in your config file so that its actually
exposed over the docker ip.