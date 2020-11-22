# cornstone

> The single utility for all your modded Minecraft needs

## Introduction

cornstone is command-line utility for Minecraft which allows you to easily download, configure and update modpacks for both clients and servers. Specifically designed for 1-click deployment, distributing your modpack has never been simpler.

## Features

- Operate on Curse/Corn manifests
- Install [MultiMC](https://multimc.org/) with bundled portable Java
- Install Curse/Corn modpack from URL or file
- Install Curse/Corn server modpack from URL or file

## Supported Platforms

- Windows x64
- Linux x64
  - Make sure you satisfy the MultiMC [dependencies](https://multimc.org/)
- macOS x64

## Usage

For help on any command, simply invoke it with `--help`, like so:

```bash
cornstone --help
cornstone manifest --help
```

## Corn manifest

cornstone uses a special modpack manifest called the Corn manifest.
The Corn manifest, saved under the file name `manifest.json`, is a superset of the well-known Curse manifest used by all major launchers like the Twitch Launcher. Being a superset means that the Corn manifest is **completely backwards-compatible** with the Curse manifest. You can use Curse modpacks with cornstone, and you can use Corn modpacks in Twitch Launcher. The only difference is the addition of some very helpful fields, prefixed with an underscore, which are only understood by cornstone. The differences are illustrated below:

> **NOTE:** The `_metadata` field is auto-generated and only used for the developer's convenience. You should not edit it. Use the `format` command to generate or update it.

```diff
{
  ...
  "files": [
    {
      "projectID": 244049,
      "fileID": 3097937,
      "required": true,
+     "_metadata": {
+       "projectName": "Woot",
+       "fileName": "woot-1.16.3-1.0.0.1.jar",
+       "summary": "Build your own mob factory",
+       "websiteUrl": "https://www.curseforge.com/minecraft/mc-mods/woot",
+       "hash": "2440493097937"
+     },
+     "_serverIgnored": false
    }
  ],
+ "_externalFiles": [
+   {
+     "name": "Spark",
+     "url": "https://ci.lucko.me/job/spark/151/artifact/spark-forge/build/libs/spark-forge.+ar",
+     "installPath": "mods/spark-forge.jar",
+     "required": true,
+     "extract": {
+       "enable": false,
+       "unwrap": false
+     },
+     "serverIgnored": false
+   }
+ ]
}
```
