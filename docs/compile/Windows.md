# Windows Build

## From Windows (Native Compilation)

[Setup a development envioronment using MSYS2 and Golang](https://gist.github.com/glycerine/355121fc4cc525b81d057d3882673531)
> ###### make sure to also read the comments on the above guide

Run the MinGW64 shell from the start menu
```bash
git clone https://github.com/BrownNPC/Youtube-Music-TUI
cd ./Youtube-Music-TUI/
python buildWindows.py
```
---
## From Linux or WSL2 (Cross-Compile to Windows)

> #####  Install mingw-gcc and 7zip packages for your distro

### Ubuntu | Debian  | Mint | Pop!_OS | WSL
```
sudo apt update
sudo apt install mingw-w64 p7zip-full
```
### Fedora | Nobara
```
sudo dnf install mingw64-gcc  p7zip
```
### Arch Linux
```
sudo pacman -S mingw-w64-gcc p7zip
```
### OpenSUSE
```
sudo zypper install mingw64-gcc p7zip
```



> #####  Run the command sequence
```bash
git clone https://github.com/BrownNPC/Youtube-Music-TUI
cd ./Youtube-Music-TUI/
python buildWindows.py
```
---
#### [⬅ Back to homepage](https://github.com/BrownNPC/Youtube-Music-TUI/?tab=readme-ov-file#compiling)
---