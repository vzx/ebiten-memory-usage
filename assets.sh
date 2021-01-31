#!/bin/bash

set -eu

cd assets

curl -o deepfield.png https://upload.wikimedia.org/wikipedia/commons/2/22/Hubble_Extreme_Deep_Field_%28full_resolution%29.png
curl -o mplus.tar.xz "https://gemmei.ftp.acc.umu.se/mirror/osdn.net/mplus-fonts/62344/mplus-TESTFLIGHT-063a.tar.xz"
tar xJf mplus.tar.xz

go generate

rm -f mplus.tar.xz deepfield.png
rm -rf mplus-TESTFLIGHT-063a/
