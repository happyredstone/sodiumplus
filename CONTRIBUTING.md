# Contributing

Sodium Plus is open source and under the [MIT License](./LICENSE), and contributions are very encouraged! Pull requests and issues will be looked at as soon as possible.

## Packwiz

We use packwiz to manage mods and assets, so please use it for that too! Their docs are here: https://packwiz.infra.link

## Mod List

To automatically generate a mod list from all the `.pw.toml` files in the `mods/` and `resourcepacks/` folders, just run the python script `modlist.py` and it will generate a `MODLIST.md`.

## Changelogs

To create changelogs, you can run the python script (`changelog.py`) to auto-pull all the latest changelogs from Modrinth, but don't modify the CHANGELOG.md from that. That should be manually created. The script will generate a CHANGELIST.md that you can base that off of, but will not be in the repository for others to see.
