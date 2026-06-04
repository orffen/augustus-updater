# Augustus Updater <img src="assets/augustus-updater.png" alt="Augustus Updater" width="48" height="48">

[![Build Augustus Updater](https://github.com/orffen/augustus-updater/actions/workflows/build.yml/badge.svg)](https://github.com/orffen/augustus-updater/actions/workflows/build.yml)

The Augustus Updater is a simple utility to automatically check for and download the latest [Augustus](https://github.com/Keriew/augustus) [Unstable build](https://augustus.josecadete.net/) for your platform before running it.

It currently supports Windows (64-bit), macOS (Universal) and Linux (AppImage), and will automatically execute Augustus Unstable after updating.

## Installation & Usage

| Platform | Latest release |
| --- | --- |
| Windows 64-bit | [![Download for Windows](https://img.shields.io/badge/Download-Windows-blue?logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiIgZmlsbD0id2hpdGUiIHZpZXdCb3g9IjAgMCAxNiAxNiI+PHBhdGggZD0iTTcuNDYyIDBIMHY3LjE5aDcuNDYyek0xNiAwSDguNTM4djcuMTlIMTZ6TTcuNDYyIDguMjExSDBWMTZoNy40NjJ6bTguNTM4IDBIOC41MzhWMTZIMTZ6Ii8+PC9zdmc+)](https://github.com/orffen/augustus-updater/releases/latest/download/augustus-updater.exe) |
| macOS | [![Download for macOS](https://img.shields.io/badge/Download-macOS-blue?logo=Apple&logoColor=white)](https://github.com/orffen/augustus-updater/releases/latest/download/augustus-updater-mac) |
| Linux amd64 | [![Download for Linux amd64](https://img.shields.io/badge/Download-Linux_amd64-blue?logo=Linux&logoColor=white)](https://github.com/orffen/augustus-updater/releases/latest/download/augustus-updater-linux-amd64) |

On Windows and Linux, it will download to the current directory. You can run it from the same directory you want it to install Augustus Unstable into, or `cd` into the directory and run it from another location.

On macOS, it will always download to `~/Applications/Augustus Unstable.app`.

Note that on macOS, you may need to modify its permissions to allow it to execute. You can do this by running:
```bash
xattr -d com.apple.quarantine augustus-updater-mac && chmod +x augustus-updater-mac
```

## License & Attribution

This project is licensed under the GNU Affero General Public License v3 (AGPLv3). See the `LICENSE.txt` file for the full license text.

- **Application Icon:** Based on `augustus_512.png` by the [Augustus](https://github.com/Keriew/augustus) developers (licensed under AGPLv3), modified by Steve Simenic (2026).
