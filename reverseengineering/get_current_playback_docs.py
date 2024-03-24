# NO PLAYER LOGGED IN = NONE

# IS PAUSED

status = {
    'is_playing': False,
    "shuffle_state": False,
    'repeat_state': 'off',
    'progress_ms': 5000,
    
    # video id
    'currently_playing': None,
}

is_paused={
    "device": {
        "id": "6848ed3a99ff09722e4280af1cb375db02bdf63b",
        "is_active": True,
        "is_private_session": False,
        "is_restricted": False,
        "name": "Web Player (Chrome)",
        "supports_volume": True,
        "type": "Computer",
        "volume_percent": 15,
    },
    "shuffle_state": False,
    "smart_shuffle": False,
    "repeat_state": "off",
    "timestamp": 1710892583456,
    "context": {
        "external_urls": {
            "spotify": "https://open.spotify.com/playlist/37i9dQZF1DZ06evO0TOYhj"
        },
        "href": "https://api.spotify.com/v1/playlists/37i9dQZF1DZ06evO0TOYhj",
        "type": "playlist",
        "uri": "spotify:playlist:37i9dQZF1DZ06evO0TOYhj",
    },
    "progress_ms": 54408,
    "item": { # IMPORTANT TO PORT
        "album": {
            "album_type": "album",
            "artists": [
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/1Cd373x8qzC7SNUg5IToqp"
                    },
                    "href": "https://api.spotify.com/v1/artists/1Cd373x8qzC7SNUg5IToqp",
                    "id": "1Cd373x8qzC7SNUg5IToqp",
                    "name": "BoyWithUke",
                    "type": "artist",
                    "uri": "spotify:artist:1Cd373x8qzC7SNUg5IToqp",
                }
            ],
            "available_markets": [
                "ET",
                "XK",
            ],
            "external_urls": {
                "spotify": "https://open.spotify.com/album/1I79ZTFJ5FVLwMYRWvhk73"
            },
            "href": "https://api.spotify.com/v1/albums/1I79ZTFJ5FVLwMYRWvhk73",
            "id": "1I79ZTFJ5FVLwMYRWvhk73",
            "images": [
                {
                    "height": 640,
                    "url": "https://i.scdn.co/image/ab67616d0000b273ee07023115f822012390d2a0",
                    "width": 640,
                },
                {
                    "height": 300,
                    "url": "https://i.scdn.co/image/ab67616d00001e02ee07023115f822012390d2a0",
                    "width": 300,
                },
                {
                    "height": 64,
                    "url": "https://i.scdn.co/image/ab67616d00004851ee07023115f822012390d2a0",
                    "width": 64,
                },
            ],
            "name": "Serotonin Dreams",
            "release_date": "2022-05-06",
            "release_date_precision": "day",
            "total_tracks": 11,
            "type": "album",
            "uri": "spotify:album:1I79ZTFJ5FVLwMYRWvhk73", #USED BY TRACKSMENU.PY IN refresh_now_playing
        },
        "artists": [
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/1Cd373x8qzC7SNUg5IToqp"
                },
                "href": "https://api.spotify.com/v1/artists/1Cd373x8qzC7SNUg5IToqp",
                "id": "1Cd373x8qzC7SNUg5IToqp",
                "name": "BoyWithUke",
                "type": "artist",
                "uri": "spotify:artist:1Cd373x8qzC7SNUg5IToqp",
            }
        ],
        "available_markets": [
            "AR",
            "AU",
        ],
        "disc_number": 1,
        "duration_ms": 171946,
        "explicit": True,
        "external_ids": {"isrc": "USUG12201585"},
        "external_urls": {
            "spotify": "https://open.spotify.com/track/6dGqGkYDoRrKh5UiIcTT22"
        },
        "href": "https://api.spotify.com/v1/tracks/6dGqGkYDoRrKh5UiIcTT22",
        "id": "6dGqGkYDoRrKh5UiIcTT22",
        "is_local": False,
        "name": "Understand",
        "popularity": 71,
        "preview_url": None,
        "track_number": 9,
        "type": "track",
        "uri": "spotify:track:6dGqGkYDoRrKh5UiIcTT22",
    },
    "currently_playing_type": "track", 
    "actions": {"disallows": {"pausing": True}},
    "is_playing": False, # IMPORTANT TO PORT
}

# IS PLAYING
is_playing = {
    "device": {
        "id": "6848ed3a99ff09722e4280af1cb375db02bdf63b",
        "is_active": True,
        "is_private_session": False,
        "is_restricted": False,
        "name": "Web Player (Chrome)",
        "supports_volume": True,
        "type": "Computer",
        "volume_percent": 15,
    },
    "shuffle_state": False,
    "smart_shuffle": False,
    "repeat_state": "off",
    "timestamp": 1710893194780,
    "context": {
        "external_urls": {
            "spotify": "https://open.spotify.com/playlist/37i9dQZF1DZ06evO0TOYhj"
        },
        "href": "https://api.spotify.com/v1/playlists/37i9dQZF1DZ06evO0TOYhj",
        "type": "playlist",
        "uri": "spotify:playlist:37i9dQZF1DZ06evO0TOYhj",
    },
    "progress_ms": 57554,
    "item": {
        "album": {
            "album_type": "album",
            "artists": [
                {
                    "external_urls": {
                        "spotify": "https://open.spotify.com/artist/1Cd373x8qzC7SNUg5IToqp"
                    },
                    "href": "https://api.spotify.com/v1/artists/1Cd373x8qzC7SNUg5IToqp",
                    "id": "1Cd373x8qzC7SNUg5IToqp",
                    "name": "BoyWithUke",
                    "type": "artist",
                    "uri": "spotify:artist:1Cd373x8qzC7SNUg5IToqp",
                }
            ],
            "available_markets": [
                "ET",
                "XK",
            ],
            "external_urls": {
                "spotify": "https://open.spotify.com/album/1I79ZTFJ5FVLwMYRWvhk73"
            },
            "href": "https://api.spotify.com/v1/albums/1I79ZTFJ5FVLwMYRWvhk73",
            "id": "1I79ZTFJ5FVLwMYRWvhk73",
            "images": [
                {
                    "height": 640,
                    "url": "https://i.scdn.co/image/ab67616d0000b273ee07023115f822012390d2a0",
                    "width": 640,
                },
                {
                    "height": 300,
                    "url": "https://i.scdn.co/image/ab67616d00001e02ee07023115f822012390d2a0",
                    "width": 300,
                },
                {
                    "height": 64,
                    "url": "https://i.scdn.co/image/ab67616d00004851ee07023115f822012390d2a0",
                    "width": 64,
                },
            ],
            "name": "Serotonin Dreams",
            "release_date": "2022-05-06",
            "release_date_precision": "day",
            "total_tracks": 11,
            "type": "album",
            "uri": "spotify:album:1I79ZTFJ5FVLwMYRWvhk73",
        },
        "artists": [
            {
                "external_urls": {
                    "spotify": "https://open.spotify.com/artist/1Cd373x8qzC7SNUg5IToqp"
                },
                "href": "https://api.spotify.com/v1/artists/1Cd373x8qzC7SNUg5IToqp",
                "id": "1Cd373x8qzC7SNUg5IToqp",
                "name": "BoyWithUke",
                "type": "artist",
                "uri": "spotify:artist:1Cd373x8qzC7SNUg5IToqp",
            }
        ],
        "available_markets": [
            "ET",
            "XK",
        ],
        "disc_number": 1,
        "duration_ms": 171946,
        "explicit": True,
        "external_ids": {"isrc": "USUG12201585"},
        "external_urls": {
            "spotify": "https://open.spotify.com/track/6dGqGkYDoRrKh5UiIcTT22"
        },
        "href": "https://api.spotify.com/v1/tracks/6dGqGkYDoRrKh5UiIcTT22",
        "id": "6dGqGkYDoRrKh5UiIcTT22",
        "is_local": False,
        "name": "Understand",
        "popularity": 71,
        "preview_url": None,
        "track_number": 9,
        "type": "track",
        "uri": "spotify:track:6dGqGkYDoRrKh5UiIcTT22",
    },
    "currently_playing_type": "track",
    "actions": {"disallows": {"resuming": True}},
    "is_playing": True,
}
