import curses
import sys
import argparse
from contextlib import contextmanager, redirect_stderr, redirect_stdout
from os import devnull
from ytm_tui.src.ui import App
from ytm_tui.__version__ import __version__

parser = argparse.ArgumentParser(
                    prog = 'YT-Music Tui',
                    description = 'Youtube music in the terminal')
parser.add_argument('-v', '--version',
                    action='store_true')
args = parser.parse_args()

if args.version:
    print(__version__)
    sys.exit(0)

@contextmanager
def suppress_stdout_stderr():
    """A context manager that redirects stdout and stderr to devnull"""
    with open(devnull, 'w') as fnull:
        with redirect_stderr(fnull) as err, redirect_stdout(fnull) as out:
            yield (err, out)


def main():
    with suppress_stdout_stderr():
        curses.wrapper(App)


if __name__ == "__main__":
    main()
