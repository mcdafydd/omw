#!/bin/sh

APP="Omw.app"
mkdir -p build/$APP/Contents/{MacOS,Resources}

go build -o build/$APP/Contents/MacOS/OutOfMyWay
cat > build/$APP/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>OutOfMyWay</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>com.mcdafydd.outofmyway</string>
</dict>
</plist>
EOF
cp icons/icon.icns build/$APP/Contents/Resources/icon.icns
find $APP
