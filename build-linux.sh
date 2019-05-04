#!/bin/sh

APP=outofmyway-timetracker
APPDIR=build/${APP}_0.1.5

mkdir -p $APPDIR/usr/bin
mkdir -p $APPDIR/usr/share/applications
mkdir -p $APPDIR/usr/share/icons/hicolor/1024x1024/apps
mkdir -p $APPDIR/usr/share/icons/hicolor/256x256/apps
mkdir -p $APPDIR/DEBIAN

go build -o $APPDIR/usr/bin/$APP

cp icons/icon.png $APPDIR/usr/share/icons/hicolor/1024x1024/apps/${APP}.png
cp icons/icon.png $APPDIR/usr/share/icons/hicolor/256x256/apps/${APP}.png

cat > $APPDIR/usr/share/applications/${APP}.desktop << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=$APP
Exec=$APP
Icon=$APP
Terminal=false
StartupWMClass=Lorca
EOF

cat > $APPDIR/DEBIAN/control << EOF
Package: ${APP}
Version: 1.0-0
Section: base
Priority: optional
Architecture: amd64
Maintainer: David McPike <david@linkconsultinginc.com>
Description: Out Of My Way Time Tracker
EOF

dpkg-deb --build $APPDIR
