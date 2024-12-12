# Contributing

Sodium Plus is open source and under the [MIT License](./LICENSE), and contributions are very encouraged! Pull requests and issues will be looked at and worked on as soon we are able.

## Packwiz

We use packwiz to manage mods and assets, so please use it for that too! Their docs are here: https://packwiz.infra.link

## Packaging

To package, use our custom builder tool, which can be built just like any other go program. All formats can be packaged using the `bundle` subcommand.

## Mod List

To automatically generate a mod list from all the `.pw.toml` files in the `mods/` and `resourcepacks/` folders, use our builder tool with the arguments `list --markdown --output MODLIST.md`, and it will generate a mod list for you in the `MODLIST.md` file.

## Changelogs

We adhere to the Keep a Changelog format. See [the changelog](./CHANGELOG.md) for more details.
