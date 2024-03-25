import sys
import time
import curses
from threading import Thread, Lock
from ytm_tui.src.util import debounce
from ytm_tui.src.Logging import logging
from ytm_tui.src.youtubeApi import YoutubeAPI
from ytm_tui.src.TracksMenu import TracksMenu
from ytm_tui.src.LibraryMenu import LibraryMenu
from ytm_tui.src.PlaylistMenu import PlaylistMenu
from ytm_tui.src.DeviceMenu import DeviceMenu
from ytm_tui.src.SearchInput import SearchInput
from ytm_tui.src.NowPlaying import NowPlaying

starttime = time.time()
lock = Lock()

class MainForm:
    def __init__(self, stdscr):
        self.stdscr = stdscr
        self.api = YoutubeAPI()
        self.pause_updates = False
        self.device_id = None
        self.tracklist_uri = None
        self.status = self.api.get_playing()
        self.active = True

        self.app_name = "SpoTUI"

        # Events
        self.events = {
            155: self.handle_exit,
            27: self.handle_esc,
            ord("q"): self.handle_exit,
            9: self.select_next_component,
            curses.KEY_RESIZE: self.handle_resize,
            ord("d"): self.show_device_menu,
            ord("/"): self.show_search_bar,
            ord("i"): self.show_inner_search_bar,
            ord(" "): self.toggle_playback,
            ord("p"): self.previous_track,
            ord("n"): self.next_track,
            ord("s"): self.toggle_shuffle,
            ord("r"): self.cycle_repeat,
            ord(","): self.lower_volume,
            ord("."): self.raise_volume,
            curses.KEY_RIGHT: self.seek_forward,
            curses.KEY_LEFT: self.seek_backward,
        }

        # window size
        scry, scrx = self.stdscr.getmaxyx()

        # UI components
        self.components = [
            TracksMenu(stdscr, self.api, self.play_track, self.status),
            # LibraryMenu(stdscr, self.api, self.change_tracklist),
            PlaylistMenu(stdscr, self.api, self.change_tracklist),
            NowPlaying(stdscr),
        ]
        self.search_component = SearchInput(self.stdscr, self.api, "Search",
                                            self.search)
        self.inner_search_component = SearchInput(self.stdscr, self.api,
                                                  "Inner Search", self.inner_search)
        self.device_menu_component = DeviceMenu(self.stdscr, self.api,
                                                self.select_device,
                                                self.hide_popup)

        # Active component
        self.active_component = 0
        self.components[0].activate()

        # Popups
        self.popup = None

        # Searches
        self.inside_search = False

        # Stack
        self.previous_tracklist = None
        self.tracklist_stack = []

        # Set initial tracklist
        # if self.status and 'context' in self.status and type(self.status["context"]) is dict and 'uri' in self.status["context"]:
        #     self.change_tracklist(
        #         self.api.get_playlist_tracks(self.status["context"]["uri"]), "Previous Session")
        # else:
        #     self.change_tracklist(self.api.get_top_tracks(), "Top Tracks")

        # Set initial device ID
        # devices = None

        # Initial render
        self.render()

        # Poll playing status every second in a new thread
        status_loop = Thread(target=self.status_loop)
        status_loop.daemon = True
        status_loop.start()

        # Start the main event loop (used for responding to key presses and keeping the main process running)
        while self.active:
            try:
                if not self.pause_updates:
                    # capture and handle key press
                    key = self.stdscr.getch()
                    if key in self.events.keys():
                        # run the event handler for the key
                        self.events[key]()
                        self.render()
                    elif self.popup:
                        # or pass it to the active popup
                        self.popup.receive_input(key)
                    else:
                        # or pass the input to the active component
                        self.components[self.active_component].receive_input(
                            key)
                    # re-render
                    self.render()
                    # self.render()
            except KeyboardInterrupt:
                sys.exit(0)

    def status_loop(self):
        while self.active:
            try:
                if not self.pause_updates:                
                    self.status = self.api.get_playing()
                    self.components[0].refresh_now_playing(self.status) #tracksmenu.py

                with lock:
                    if not self.pause_updates:
                        self.render()
                time.sleep(1 - ((time.time() - starttime) % 1))
            except Exception as e:
                pass

    def render(self):
        self.stdscr.erase()
        for component in self.components:
            # render each component
            component.render(self.status)
        if self.popup:
            self.popup.render()
        self.stdscr.refresh()

    # events
    def change_tracklist(self, tracks, title, tracklist_uri=None):
        self.components[0].update_tracks(tracks, title)
        self.tracklist_uri = tracklist_uri
        self.tracklist_stack.append(self.previous_tracklist)
        self.previous_tracklist = {"tracks": tracks, "title": title, "tracklist_uri": tracklist_uri}
        self.activate_tracklist()

    def select_next_component(self):
        if self.popup:
            return
        # visually de-activate the current component
        self.components[self.active_component].deactivate()
        # incremement the active component (or go back to start)
        self.active_component = (
            self.active_component +
            1 if self.active_component < len(self.components) - 1 else 0)
        # skip read-only components
        if self.components[self.active_component].interactive:
            self.components[self.active_component].activate()
        else:
            self.select_next_component()

    def play_track(self, track):
        self.api.start_playback(track)

    @debounce(0.5)
    def toggle_playback(self):
        if not self.status:
            return
        self.api.toggle_playback()


    def lower_volume(self):
        if self.status and self.status["is_playing"]:
            self.api.change_volume(-10)

    def raise_volume(self):
        if self.status and self.status["is_playing"]:
            self.api.change_volume(+10)
    @debounce(0.5)
    def previous_track(self):
        if self.status and self.status["is_playing"]:
            self.api.previous_track()

    # @debounce(2)
    def next_track(self):
        if self.status and self.status["is_playing"]:
            self.api.next_track()
    @debounce(1.5)
    def toggle_shuffle(self):
        if self.api.loaded_tracks_ids: # if a playlist is loaded
            self.api.toggle_shuffle()

    @debounce(1.5)
    def cycle_repeat(self):
        if self.status["is_playing"]:
            self.api.repeat()
    
    # @debounce(0)
    def seek_backward(self):
        if self.status and self.status["is_playing"]:
            progress = self.status["progress_ms"]
            self.api.seek_track(-10)

    # @debounce(0)
    def seek_forward(self):
        if self.status and self.status["is_playing"]:
            self.api.seek_track(10)

    def search(self, query):
        self.hide_popup()
        query = query.strip()
        if query and len(query) > 1:
            results = self.api.search(query)
            self.inside_search = True
            self.change_tracklist(results, "Searching: " + query)
            self.render()

    @debounce(1.5)
    def inner_search(self, query):
        self.hide_popup()
        query = query.strip()
        tracks = self.components[0].tracks

        if query and len(query) > 1:
            results = []
            for track in tracks:
                if query.upper() in track['name'].upper() or query.upper() in track['artist'].upper():
                    results.append(track)
            self.inside_search = True
            self.change_tracklist(results, f'Searching: {query}')
            self.render()

    def activate_tracklist(self):
        self.components[self.active_component].deactivate()
        self.active_component = 0
        self.components[self.active_component].activate()

    @debounce(2)
    def show_device_menu(self):
        self.components[self.active_component].deactivate()
        self.popup = self.device_menu_component
        self.popup.restart()
        self.popup.activate()
        self.render()

    def show_search_bar(self):
        if self.popup:
            return
        self.pause_updates = True
        self.popup = self.search_component
        self.components[self.active_component].deactivate()
        self.popup.activate()
        self.render()


    def show_inner_search_bar(self):
        if self.popup:
            return
        self.pause_updates = True
        self.popup = self.inner_search_component
        self.components[self.active_component].deactivate()
        self.popup.activate()
        self.render()

    def select_device(self, device_id):
        self.device_id = device_id

    def hide_popup(self):
        if self.popup:
            self.popup.deactivate()
        self.popup = None
        self.components[self.active_component].activate()
        self.pause_updates = False
        self.stdscr.clear()
        self.render()

    def handle_resize(self):
        for component in self.components:
            # render each component
            component.restart()
        self.stdscr.clear()

    def handle_esc(self):
        if self.popup:
            self.hide_popup()
        elif self.inside_search:
            self.inside_search = False
            self.change_tracklist(**self.tracklist_stack.pop())
            self.previous_tracklist = None
            self.render()

    def handle_exit(self):
        sys.exit(0)

    def __filter_tracklist(self, track):
        return track["type"] == 'track'

    def __map_tracklist(self, track):
        return track["uri"]
