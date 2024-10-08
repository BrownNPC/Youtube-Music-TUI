# Windows Build

## From Windows (Native Compilation)

[Setup a development envioronment using MSYS2 and Golang](https://gist.github.com/glycerine/355121fc4cc525b81d057d3882673531)
> ###### make sure to also read the comments on the above guide

Run the MinGW64 shell from the start menu
```bash
pacman -S mingw64/mingw-w64-x86_64-mpv
git clone https://github.com/BrownNPC/Youtube-Music-TUI
cd Youtube-Music-TUI/ytt
go mod tidy
CGO_ENABLED=1 go build .
```
---
#### [⬅ Back to homepage](https://github.com/BrownNPC/Youtube-Music-TUI/?tab=readme-ov-file#features-)
---


## From Linux or WSL2 (Cross-Compile to Windows)

Install mingw-gcc and mpv dev packages for your distro

### Ubuntu | Debian  | Mint | Pop!_OS | WSL
```
sudo apt update
sudo apt install mingw-w64 libmpv-dev
```
### Fedora | Nobara
```
sudo dnf install mingw64-gcc mpv-devel
```
### Arch Linux
```
sudo pacman -S mingw-w64-gcc mpv
```
### OpenSUSE
```
sudo zypper install mingw64-gcc libmpv-devel
```




---
#### [⬅ Back to homepage](https://github.com/BrownNPC/Youtube-Music-TUI/?tab=readme-ov-file#features-)
---