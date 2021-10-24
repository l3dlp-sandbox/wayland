package main

import "os"
import "image/png"
import "image/color"
import "strings"
import "fmt"

type Font struct {
	cellx   int
	celly   int
	mapping map[string][][3]byte
}

var hexfont = " #    # ##  ##   #  ###  ## ###  #   #   #  ##   #  ##  ### ### " +
	"# #  ##   #   # #   #   #     # # # # # # # # # # # # # #   #   " +
	"# # # #   #  #  ### ##  ##   #   #   ## # # ##  #   # # ### ### " +
	"# #   #  #    #  #    # # #  #  # #   # ### # # # # # # #   #   " +
	" #    # ### ##   #  ##   #   #   #  ##  # # ##   #  ##  ### #   " +
	"                                                                "

func hexfontGet(hex, x, y byte) bool {
	switch hex {
	case '0':
		hex = 0
	case '1':
		hex = 1
	case '2':
		hex = 2
	case '3':
		hex = 3
	case '4':
		hex = 4
	case '5':
		hex = 5
	case '6':
		hex = 6
	case '7':
		hex = 7
	case '8':
		hex = 8
	case '9':
		hex = 9
	case 'A', 'a':
		hex = 10
	case 'B', 'b':
		hex = 11
	case 'C', 'c':
		hex = 12
	case 'D', 'd':
		hex = 13
	case 'E', 'e':
		hex = 14
	case 'F', 'f':
		hex = 15
	}
	if hex >= 16 {
		panic("")
	}
	if x >= 4 {
		panic("")
	}
	if y >= 6 {
		panic("")
	}

	return hexfont[4*int(hex)+int(y)*64+int(x)] == '#'
}

func (f *Font) GetRGBTexture(code string) [][3]byte {

	var a, ok = f.mapping[code]
	if !ok {
		if f.cellx < 12 || f.celly < 24 {
			return nil
		}

		faketexture := make([][3]byte, f.cellx*f.celly)
		fakestring := fmt.Sprintf("%+q", code)
		if len(fakestring) > 3 && (fakestring[0:3] == "\"\\u" || fakestring[0:3] == "\"\\U") {
			fakestring = fakestring[3:]
		}
		if len(fakestring) > 1 && fakestring[0:1] == "\"" {
			fakestring = fakestring[1:]
		}
		if len(fakestring) >= 1 && fakestring[len(fakestring)-1] == '"' {
			fakestring = fakestring[0 : len(fakestring)-1]
		}
		println(fakestring)
		var i = 0
		for ybox := byte(0); ybox < 4; ybox++ {
			for xbox := byte(0); xbox < 3; xbox++ {

				for y := byte(0); y < 6; y++ {
					for x := byte(0); x < 4; x++ {
						pos := int(ybox)*f.cellx*6 + int(xbox)*4 + int(y)*f.cellx + int(x)
						if len(fakestring) > i {
							if hexfontGet(fakestring[i], x, y) {
								faketexture[pos][0] = 255
								faketexture[pos][1] = 255
								faketexture[pos][2] = 255
							}
						}
					}
				}
				i++
			}
		}
		// memoization
		f.mapping[code] = faketexture
		return faketexture
	}
	return a
}

func (f *Font) Alias(alias, key string) error {
	if f.mapping == nil {
		println("no mapping")
		return fmt.Errorf("no mapping")
	}
	if f.mapping[key] == nil {
		println("key missing")
		return fmt.Errorf("key missing")
	}
	f.mapping[alias] = f.mapping[key]
	return nil
}

func (f *Font) Load(name, descriptor string) error {
	file, err := os.Open(name)
	if err != nil {
		print("Font not found: ")
		println(name)
		return err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		print("Cannot decode png: ")
		println(name)
		return err
	}
	b := img.Bounds()

	var width = b.Max.X - b.Min.X
	var height = b.Max.Y - b.Min.Y

	var buffer = strings.Split(strings.ReplaceAll(descriptor, "\r\n", "\n"), "\n")
	var buf0 = strings.Split(buffer[0], "\t")

	var cellx = width / len(buf0)
	var celly = height / len(buffer)

	if f.mapping == nil {
		f.cellx = cellx
		f.celly = celly
	} else if f.cellx != cellx || f.celly != celly {
		return fmt.Errorf("only same cell sized fonts can be merged")
	}

	var mapping = make(map[string][2]int)
	var mapping2 = make(map[[2]int][][3]byte)

	for y, v := range buffer {
		var buf = strings.Split(strings.Trim(v, "\t"), "\t")
		for x, cell := range buf {

			mapping[cell] = [2]int{x, y}
		}
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		var iy = (y - b.Min.Y) / f.celly
		for x := b.Min.X; x < b.Max.X; x++ {
			var ix = (x - b.Min.X) / f.cellx
			var i = [2]int{ix, iy}

			var sli = mapping2[i]

			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)

			sli = append(sli, [3]byte{c.R, c.G, c.B})

			mapping2[i] = sli
		}
	}
	if f.mapping == nil {
		f.mapping = make(map[string][][3]byte)
	}
	for k, v := range mapping {
		f.mapping[k] = mapping2[v]
	}
	return nil
}
