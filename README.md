# Quake Timer

A Windows background application to manage cooldown timers for Quake Champions.

## Getting Started

Download the latest `.zip` from [the Releases page](https://github.com/rocheio/quake-timer/releases).

Unzip it, run the `.exe`, and click through the Windows security warnings (fixing the app so it doesn't have these warnings will be done in the future).

In Quake Champions, set your Video Display Mode to `Borderless` or `Windowed`.

Ensure your volume is loud enough to hear the alerts.

Press the registered hotkey at the appropriate time during gameplay:

- Mega Health: `Alt + 1` (30 second respawn)
    - Press immediately after someone picks it up
- Heavy Armor: `Alt + 2` (30 second respawn)
    - Press immediately after someone picks it up
- Quad Damage: `Alt + 3` (120 second respawn)
    - Press immediately after someone picks up the Protection
- Protection: `Alt + 4` (120 second respawn)
    - Press immediately after someone picks up the Quad Damage

The item's name will be announced 10 seconds and 20 seconds before the timer expires.

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
