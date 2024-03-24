import curses
from ytm_tui.src.config import get_config
from ytm_tui.src.util import theme_load

def init_colors():

    theme = get_config()['theme']

    # Default theme
    
    highlight_text = curses.COLOR_BLACK
    highlight_box = curses.COLOR_WHITE
    inactive_menu = curses.COLOR_WHITE
    active_menu = curses.COLOR_YELLOW
    search_box = curses.COLOR_MAGENTA
    progress_bar = curses.COLOR_GREEN

    # config theme
    if theme.get('highlight_text'):
        highlight_text = theme_load(theme['highlight_text'])
    if theme.get('highlight_box'):
        highlight_box = theme_load(theme['highlight_box'])
    if theme.get('inactive_menu'):
        inactive_menu = theme_load(theme['inactive_menu'])
    if theme.get('active_menu'):
        active_menu = theme_load(theme['active_menu'])
    if theme.get('search_box'):
        search_box = theme_load(theme['search_box'])
    if theme.get('progress_bar'):
        progress_bar = theme_load(theme['progress_bar'])


    # default
    curses.init_pair(1, curses.COLOR_WHITE, 0)
    # White text (inactive color)
    curses.init_pair(4, inactive_menu, curses.COLOR_BLACK)
    # Yellow text (active_color)
    curses.init_pair(5, active_menu, curses.COLOR_BLACK)
    # Magenta text (search_box_color)
    curses.init_pair(10, search_box, curses.COLOR_BLACK)
    # Green text (Progress bar color)
    curses.init_pair(11, progress_bar, curses.COLOR_BLACK)
    # Cyan text
    curses.init_pair(12, curses.COLOR_CYAN, curses.COLOR_BLACK)
    # Selected item (highlight text color, highlight box color)
    curses.init_pair(6, highlight_text, highlight_box)
    # Highlighted (no bg)
    curses.init_pair(7, curses.COLOR_MAGENTA, curses.COLOR_BLACK)
    # Highlighted (bg)
    curses.init_pair(8, curses.COLOR_WHITE, curses.COLOR_MAGENTA)
