# Youtube-Music TUI
Listen to Youtube Music from the Terminal.
	
  
![](assets/ytt.webp)


Getting Started

-----------



Install with ```pip install ytm-tui```

Make sure [MPV](https://mpv.io/installation/) is installed. Along with [yt-dlp](https://pypi.org/project/yt-dlp/)

Run `ytt` to generate a sample config file at `~/.config/ytt/ytt.toml` or create one manually with the following:


```
playlists = [
	# most viewed songs on yt
	PL15B1E77BB5708555", 
	#lofi hip hop
	"PLofht4PTcKYnaH8w5olJCI-wUVxuoMHqM",
	#synthwave radio
	"PLUNz3rL3KK9W21UspvmRt3bwsKZFX73DE",
]
[other]
use_nerd_fonts = no
config_version = 1

# optional
[theme]
# POSSIBLE COLORS:
# COLOR_BLUE COLOR_GREEN COLOR_RED COLOR_YELLOW
# COLOR_BLACK COLOR_CYAN COLOR_MAGENTA COLOR_WHITE

progress_bar='COLOR_GREEN'
inactive_menu='COLOR_WHITE'
active_menu='COLOR_YELLOW'
search_box = 'COLOR_MAGENTA'
highlight_box='COLOR_WHITE'
highlight_text='COLOR_BLACK

```
Controls

-------
**Navigation**

`tab` Switch section

`k`/`↑` Up

`j`/`↓` Down

`g` Scroll to top

`G` Scroll to bottom

`Enter` Select

`/` Search 

`Esc`/`q` Quit/Back

**Playback**

`space` Play/Pause

`n` Next track

`p` Previous track

`→` Seek 10s forwards

`←` Seek 10s backwards

`,` Volume down

`.` Volume up

`s` Toggle shuffle

`r` Toggle repeat

------
