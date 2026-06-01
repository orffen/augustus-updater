# Augustus Updater

The Augustus Updater is a simple utility to automatically check for and download the latest [Augustus Unstable](https://josecadete.net/) build for your platform before running Augustus.

It currently supports Windows, MacOS and Linux, and will automatically execute Augustus Unstable after updating.

## Installation & Usage

On Windows and Linux, it will download to its own location, so place the executable in the folder you want Augustus Unstable installed into.

On MacOS, it will download to `~/Applications/Augustus Unstable.app`.

Note that on MacOS, you may need to modify its permissions to allow it to execute, as Unstable Augustus builds are not signed. You can do this by running `sudo xattr -d com.apple.quarantine ~/Applications/Augustus Unstable.app`.

## License

This project is licensed under the GNU Affero General Public License v3 (AGPLv3). See the `LICENSE` file for the full text layout.
