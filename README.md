[![Github All Releases](https://img.shields.io/github/downloads/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/total.svg)](https://github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj)](https://goreportcard.com/report/github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj)
![Build@release](https://github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/actions/workflows/release.yml/badge.svg)
![Tests@main](https://github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/actions/workflows/test.yml/badge.svg)

### Hi there 👋

If this repo becomes taken down, then it is a sign that GitHub supports war atrocities commited by putin!

IMPORTANT!!! Please take your security seriously and use VPN while running this software from your own devices(PC/Laptop/Android/etc)!

NOTE: Since this app is not signed by any certificate for obvious reasons, Windows Defender will complain that this app is from an unknown source. It is safe to run this app anyway **unless it was downloaded from another location**. In which case the app may be altered, so be vigilant and compare the hashes of your binaries!

![GUI](Capture.JPG)

This repository was created in order to help defend against Russia propaganda during the war activities of Russian armies in Ukraine 2022.

There are two variants of the program:
- GUI
- CLI

For GUI simply download the latest program for your operating system from the Releases section.
It is possible to run CLI program using Docker.


# Usage
## Desktop App
- `RusskijKorablIdiNaxuj.exe.zip` for Windows platforms 
- `RusskijKorablIdiNaxuj.apk` for Android platforms 
- `RusskijKorablIdiNaxuj.tar.xz` for Linux platforms 
- Please build manually for Mac-OS. 

After executing the binary there will be a window similar to the one in the screenshot above. If you have a less formidable PC, then reduce "Workers" number to something like 500.
Then you can click on the little triangle on the right and start the process.

## Command line

Just clone this repo, install Go and from clonned directory do (preferrably on a VPS that is close to Russia):
```
$ go run ./cmd/RusskijKorablIdiNaxuj-cli -i targets/targets.txt
```

Also, you can install the GUI executable using this command(although you will still need Fyne dependencies for your system):
```
go install github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj@latest
```

Or the following command for the cli version:

```
go install github.com/RusskijKorablIdiNaxuj/RusskijKorablIdiNaxuj/cmd/RusskijKorablIdiNaxuj-cli@latest
```

NOTE: It may consume a lot of RAM as it tries to leave connections open for a as long time as possible. Minimum 4GB is needed.

## Docker

Docker image is not on Dockerhub yet. Not sure if it is useful.

# Important resources

- International Ukraine DDOS Alliance Telegram: https://t.me/ukraineddos
- IT ARMY of Ukraine Telegram: https://t.me/itarmyofukraine2022

# Contribute targets/Code

Either create an issue or a PR.
