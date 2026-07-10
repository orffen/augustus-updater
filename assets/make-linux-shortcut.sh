#!/bin/sh
set -e

desktop_file="${HOME}/.local/share/applications/augustus-updater.desktop"
icon_url="https://raw.githubusercontent.com/orffen/augustus-updater/main/assets/augustus-updater.png"
script_dir=$(cd "$(dirname "$0")" && pwd)
updater_exe="augustus-updater-linux-amd64"

cd "$script_dir"

if [ ! -f "$updater_exe" ]; then
    echo "Please put this file in the same folder as $updater_exe"
    exit 1
fi

icon_file="${icon_url##*/}"
if [ ! -f "$icon_file" ]; then
    curl -sLO "$icon_url"
    echo "Downloaded icon from $icon_url"
    chmod 644 "$icon_file"
fi

mkdir -p "$(dirname "$desktop_file")"

cat << EOF > "$desktop_file"
[Desktop Entry]
Type=Application
Name=Augustus Updater
Comment=The Augustus Updater is a simple utility to automatically check for and download the latest Augustus Unstable build for your platform before running it.
Exec=${script_dir}/${updater_exe}
Path=$script_dir
Icon=${script_dir}/${icon_file}
Terminal=true
Categories=Game;
EOF

chmod +x "$desktop_file"

echo "Created $desktop_file"
