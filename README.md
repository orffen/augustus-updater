# Augustus Updater

The Augustus Updater is a simple utility to automatically check for and download the latest [Augustus Unstable](https://josecadete.net/) build for your platform before running Augustus.

It currently supports Windows (64bit), MacOS (Universal) and Linux (AppImage), and will automatically execute Augustus Unstable after updating.

## Installation & Usage

On Windows and Linux, it will download to the current directory. You can run it from the same directory you want it to install Augustus Unstable into, or `cd` into the directory and run it from another location.

On MacOS, it will always download to `~/Applications/Augustus Unstable.app`.

Note that on MacOS, you may need to modify its permissions to allow it to execute, as Augustus Updater is not signed. You can do this by running `xattr -d com.apple.quarantine augustus-updater-mac`.

## License

This project is licensed under the GNU Affero General Public License v3 (AGPLv3). See the `LICENSE.txt` file for the full text layout.
