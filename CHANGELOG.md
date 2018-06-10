# Changelog

## 0.2.0 - 2018-06-10
- Add `Manager.AddKey` method for more intuitive key creation
- Add `Hotkey.Name` attribute for better reference in output
- Use `time.After` in `Listen` to sleep safely
- Add `KeyPress` logic to parse actual time of press from `WindowsMessage`
- Add `Cooldown` class for more consistent definition

## 0.1.0 - 2018-06-08
- Working prototype with five-second alerts for each major map spawn
