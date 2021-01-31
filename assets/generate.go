package assets

//go:generate file2byteslice -package=assets -input=./deepfield.png -output=./deepfield.go -var=Deepfield_png
//go:generate file2byteslice -package=assets -input=./mplus-TESTFLIGHT-063a/mplus-1p-bold.ttf -output=./mplus.go -var=Mplus_1p_bold_ttf

import (
	_ "github.com/hajimehoshi/file2byteslice"
)
