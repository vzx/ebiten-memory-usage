package main

//go:generate curl -o deepfield.png https://upload.wikimedia.org/wikipedia/commons/2/22/Hubble_Extreme_Deep_Field_%28full_resolution%29.png
//go:generate file2byteslice -package=main -input=./deepfield.png -output=./deepfield.go -var=Deepfield_png

import (
	_ "github.com/hajimehoshi/file2byteslice"
)
