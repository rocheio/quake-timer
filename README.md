# Quake Timer

A Windows background application to manage cooldown timers for Quake Champions.

Build and run the app on Windows:

```sh
go build ./cmd/quaketimer
./quaketimer.exe
```

Turn your volume up so you can hear the alerts.

Press the registered hotkey, even when terminal is not focused.

- Mega Health: `Alt + 1` (30 seconds)
- Heavy Armor: `Alt + 2` (30 seconds)
- Quad Damage: `Alt + 3` (120 seconds)
- Protection: `Alt + 4` (120 seconds)

You should hear the item's name 5 seconds before the cooldown for each item expires.

To close the app, hit `Alt + Ctrl + X` or send a keyboard interrupt in the console (`Ctrl + C`)

## Credits

- Bootstrapping of basic global hotkey listener: https://stackoverflow.com/a/38954281
- Audio files: http://www.fromtexttospeech.com/
