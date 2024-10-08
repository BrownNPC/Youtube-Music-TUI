# Linux Build
> ##### get gcc and mpv packages for your distro
### Ubuntu | Debian  | Mint | Pop!_OS | WSL
```
sudo apt update && sudo apt install -y libmpv-dev gcc golang
```
### Fedora | Nobara
```
sudo dnf install -y mpv-devel gcc golang
```
### Arch Linux
```
sudo pacman -Syu mpv gcc go
```
### OpenSUSE
```
sudo zypper install -y mpv-devel gcc go
```
> #####  Run the command sequence
```bash
git clone https://github.com/BrownNPC/Youtube-Music-TUI
cd ./Youtube-Music-TUI/ytt
CGO_ENABLED=1 go build -o ../dist-linux/ytt . 
```
---
#### [⬅ Back to homepage](https://github.com/BrownNPC/Youtube-Music-TUI/?tab=readme-ov-file#compiling-)
---