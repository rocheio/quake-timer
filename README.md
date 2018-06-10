# Quake Timer

A Windows background application to manage cooldown timers for Quake Champions.

## Getting Started

Download the `.zip` from [the Releases page](https://github.com/rocheio/quake-timer/releases).

Unzip it, run the `.exe`, and click through the Windows security warnings (fixing the app so it doesn't have these warnings will be done in the future).

Turn your volume up so you can hear the alerts.

Press the registered hotkey, even when terminal is not focused.

- Mega Health: `Alt + 1` (30 seconds)
- Heavy Armor: `Alt + 2` (30 seconds)
- Quad Damage: `Alt + 3` (120 seconds)
- Protection: `Alt + 4` (120 seconds)

The item's name will be announced 10 seconds before the timer expires.
The longer cooldowns will also alert 30 seconds before expiration.

To close the app, hit `Alt + Ctrl + X` or send a keyboard interrupt in the console (`Ctrl + C`).

## Development

Build and run the app on Windows:

```sh
go build ./cmd/quaketimer
./quaketimer.exe
```

Other operating systems are not yet supported.

## Credits

- Bootstrapping of basic global hotkey listener: https://stackoverflow.com/a/38954281
- Audio files: http://www.fromtexttospeech.com/
