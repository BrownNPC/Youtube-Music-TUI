import os
import toml
import errno

def get_config():
    user_config_dir = os.path.expanduser("~")
    
    filename = os.path.join(user_config_dir, ".config", "ytt", "ytt.toml")
    if not os.path.exists(os.path.dirname(filename)):
        try:
            os.makedirs(os.path.dirname(filename))
        except OSError as exc:  # prevent race condition
            if exc.errno != errno.EEXIST:
                raise

    # if config file doesn't exist
    if not os.path.isfile(filename):
        print(f'Created config file in: {filename}')
        _create_default_config(filename)

    config = toml.load(filename)

    return config


def _create_default_config(filename):
    toml_string = """
#PASTE YOUTUBE PLAYLIST ID'S HERE, PLEASE DONT FORGET THE COMMAS
playlists = [
    # most viewed songs on yt
    "PL15B1E77BB5708555", 

    #lofi hip hop
    "PLofht4PTcKYnaH8w5olJCI-wUVxuoMHqM",

    #synthwave radio
    "PLUNz3rL3KK9W21UspvmRt3bwsKZFX73DE",
]


[theme]
# POSSIBLE COLORS:
# COLOR_BLUE     COLOR_GREEN    COLOR_RED      COLOR_YELLOW
# COLOR_BLACK    COLOR_CYAN     COLOR_MAGENTA  COLOR_WHITE

progress_bar='COLOR_GREEN'
inactive_menu='COLOR_WHITE'
active_menu='COLOR_YELLOW'

search_box = 'COLOR_MAGENTA'
highlight_box='COLOR_WHITE'
highlight_text='COLOR_BLACK'

[other]
use_nerd_fonts = 'no'

# SAMPLE THEMES:

#Futuristic:
# progress_bar = 'COLOR_CYAN'
# inactive_menu = 'COLOR_BLACK'
# active_menu = 'COLOR_MAGENTA'

# search_box = 'COLOR_GREEN'
# highlight_box = 'COLOR_BLUE'
# highlight_text = 'COLOR_WHITE'

#Catppuccin (kinda..): (use a pastel color-scheme terminal)
# progress_bar = 'COLOR_YELLOW'
# inactive_menu = 'COLOR_BLUE'
# active_menu = 'COLOR_MAGENTA'

# search_box = 'COLOR_GREEN'
# highlight_box = 'COLOR_CYAN'
# highlight_text = 'COLOR_WHITE'



"""
    with open(filename, "w") as config_file:
        config_file.write(toml_string)
