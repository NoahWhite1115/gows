package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const mapWidth = 24
const mapHeight = 24

type ColorRGB struct {
	R uint8
	G uint8
	B uint8
}

//using a constant map for rn
var worldMap [][]int = [][]int{
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 7, 7, 7, 7, 7, 7, 7, 7},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 7, 7, 0, 7, 7, 7, 7, 7},
	{4, 0, 5, 0, 0, 0, 0, 5, 0, 5, 0, 5, 0, 5, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 6, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 8, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 0, 5, 5, 5, 5, 7, 7, 7, 7, 7, 7, 7, 1},
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 6, 0, 6, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2},
	{4, 0, 0, 5, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3},
}

var window *sdl.Window
var renderer *sdl.Renderer
var texture *sdl.Texture
var surface *sdl.Surface

var width int32 = 640
var height int32 = 480
var texWidth = 64
var texHeight = 64
var texNumber = 8

var pixels []uint32 = make([]uint32, 4*width*height)
var keys []uint8
var textures [][]uint32 = make([][]uint32, texNumber)

func main() {

	//I tried moving all of this to it's own function, but it led to segfault/segvar errors

	text := "test"
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(text, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, 0)

	if err != nil {
		panic(err)
	}

	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING, width, height)

	if err != nil {
		panic(err)
	}

	surface, err = window.GetSurface()
	if err != nil {
		panic(err)
	}

	//init textures
	for i := 0; i < texNumber; i++ {
		textures[i] = make([]uint32, texHeight*texWidth)
	}
	err = loadTextures()

	var posX, posY, dirX, dirY, planeX, planeY, time, oldTime float64

	posX = 21
	posY = 12
	dirX = -1
	dirY = 0
	planeX = 0
	planeY = 0.66

	time = 0
	oldTime = 0

	for !Done() {

		//floor renderer
		for y := height/2 + 1; y < height; y++ {
			rayDirX0 := dirX - planeX
			rayDirY0 := dirY - planeY
			rayDirX1 := dirX + planeX
			rayDirY1 := dirY + planeY

			p := y - height/2

			posZ := 0.5 * float64(height)

			rowDistance := posZ / float64(p)
			floorStepX := rowDistance * (rayDirX1 - rayDirX0) / float64(width)
			floorStepY := rowDistance * (rayDirY1 - rayDirY0) / float64(width)

			floorX := posX + rowDistance*rayDirX0
			floorY := posY + rowDistance*rayDirY0

			var x int32
			for x = 0; x < width; x++ {
				cellX := int32(floorX)
				cellY := int32(floorY)

				tx := int(float64(texWidth)*(floorX-float64(cellX))) & (texWidth - 1)
				ty := int(float64(texHeight)*(floorY-float64(cellY))) & (texHeight - 1)

				floorX += floorStepX
				floorY += floorStepY

				checkerBoardPattern := int(cellX+cellY) & 1

				var floorTexture int

				if checkerBoardPattern == 0 {
					floorTexture = 3
				} else {
					floorTexture = 4
				}

				ceilingTexture := 6

				color := textures[floorTexture][texHeight*ty+tx]
				color = (color >> 1) & 8355711
				pixels[(x + width*y*4)] = color

				color = textures[ceilingTexture][texHeight*ty+tx]
				color = (color >> 1) & 8355711
				pixels[(x + width*(height-y-1)*4)] = color
			}

		}

		var x int32
		for x = 0; x < width; x++ {

			//ray directions
			cameraX := float64(2*x)/float64(width) - 1
			rayDirX := dirX + planeX*cameraX
			rayDirY := dirY + planeY*cameraX

			//map positions
			mapX := int(posX)
			mapY := int(posY)

			//ray length to next side
			var sideDistX, sideDistY float64

			//
			deltaDistX := math.Abs(1 / rayDirX)
			deltaDistY := math.Abs(1 / rayDirY)
			var perpWallDist float64

			var stepX, stepY int

			var side int
			hit := 0 //probably should be a bool

			if rayDirX < 0 {
				stepX = -1
				sideDistX = (posX - float64(mapX)) * deltaDistX
			} else {
				stepX = 1
				sideDistX = (float64(mapX) + 1 - posX) * deltaDistX
			}

			if rayDirY < 0 {
				stepY = -1
				sideDistY = (posY - float64(mapY)) * deltaDistY
			} else {
				stepY = 1
				sideDistY = (float64(mapY) + 1 - posY) * deltaDistY
			}

			for hit == 0 {
				if sideDistX < sideDistY {
					sideDistX += deltaDistX
					mapX += stepX
					side = 0
				} else {
					sideDistY += deltaDistY
					mapY += stepY
					side = 1
				}

				if worldMap[mapX][mapY] > 0 {
					hit = 1
				}
			}

			//Calculate distance projected on camera direction (Euclidean distance will give fisheye effect!)
			if side == 0 {
				perpWallDist = (float64(mapX) - posX + (1-float64(stepX))/2) / rayDirX
			} else {
				perpWallDist = (float64(mapY) - posY + (1-float64(stepY))/2) / rayDirY
			}

			lineHeight := int32(float64(height) / perpWallDist)

			drawStart := -lineHeight/2 + height/2
			if drawStart < 0 {
				drawStart = 0
			}

			drawEnd := lineHeight/2 + height/2
			if drawEnd >= height {
				drawEnd = height - 1
			}

			//new drawing algo
			texNum := worldMap[mapX][mapY] - 1

			var wallX float64

			if side == 0 {
				wallX = posY + perpWallDist*rayDirY
			} else {
				wallX = posX + perpWallDist*rayDirX
			}
			wallX -= math.Floor(wallX)

			texX := int(wallX * float64(texWidth))
			if side == 0 && rayDirX > 0 {
				texX = texWidth - texX - 1
			}
			if side == 1 && rayDirY < 0 {
				texX = texWidth - texX - 1
			}

			step := float64(texHeight) / float64(lineHeight)
			texPos := (float64(drawStart) - float64(height)/2 + float64(lineHeight)/2) * step

			for y := drawStart; y < drawEnd; y++ {
				texY := int(texPos) & (texHeight - 1)
				texPos += step
				color := textures[texNum][texHeight*texY+texX]

				pixels[(x + width*y*4)] = color
			}

			/*old one

			var color ColorRGB
			//fmt.Println(worldMap[mapX][mapY])

			switch {
			case worldMap[mapX][mapY] == 1:
				color = ColorRGB{R: 255, G: 0, B: 0}
			case worldMap[mapX][mapY] == 2:
				color = ColorRGB{R: 0, G: 255, B: 0}
			case worldMap[mapX][mapY] == 3:
				color = ColorRGB{R: 0, G: 0, B: 255}
			case worldMap[mapX][mapY] == 4:
				color = ColorRGB{R: 255, G: 255, B: 255}
			default:
				color = ColorRGB{R: 0, G: 255, B: 0}
			}

			if side == 1 {
				color.R = color.R / 2
				color.G = color.G / 2
				color.B = color.B / 2
			}

			VerLine(x, drawStart, drawEnd, color)
			*/
		}

		oldTime = time
		time = float64(sdl.GetTicks())
		frametime := (time - oldTime) / 1000.0
		//fmt.Println(1.0 / frametime)

		rotSpeed := frametime * 3
		moveSpeed := frametime * 5
		readKeys()

		if keyDown(sdl.SCANCODE_W) {
			if worldMap[int(posX+dirX*moveSpeed)][int(posY)] == 0 {
				posX += dirX * moveSpeed
			}
			if worldMap[int(posX)][int(posY+dirY*moveSpeed)] == 0 {
				posY += dirY * moveSpeed
			}
		}

		if keyDown(sdl.SCANCODE_S) {
			if worldMap[int(posX-dirX*moveSpeed)][int(posY)] == 0 {
				posX -= dirX * moveSpeed
			}
			if worldMap[int(posX)][int(posY-dirY*moveSpeed)] == 0 {
				posY -= dirY * moveSpeed
			}
		}

		if keyDown(sdl.SCANCODE_D) {
			oldDirX := dirX
			dirX = dirX*math.Cos(-rotSpeed) - dirY*math.Sin(-rotSpeed)
			dirY = oldDirX*math.Sin(-rotSpeed) + dirY*math.Cos(-rotSpeed)
			oldPlaneX := planeX
			planeX = planeX*math.Cos(-rotSpeed) - planeY*math.Sin(-rotSpeed)
			planeY = oldPlaneX*math.Sin(-rotSpeed) + planeY*math.Cos(-rotSpeed)
		}

		if keyDown(sdl.SCANCODE_A) {
			oldDirX := dirX
			dirX = dirX*math.Cos(rotSpeed) - dirY*math.Sin(rotSpeed)
			dirY = oldDirX*math.Sin(rotSpeed) + dirY*math.Cos(rotSpeed)
			oldPlaneX := planeX
			planeX = planeX*math.Cos(rotSpeed) - planeY*math.Sin(rotSpeed)
			planeY = oldPlaneX*math.Sin(rotSpeed) + planeY*math.Cos(rotSpeed)
		}

		Redraw()
		Cls()
	}

}

func Redraw() {
	texture.UpdateRGBA(nil, pixels, int(width)*4)
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func VerLine(x, y1, y2 int32, color ColorRGB) {
	if y2 < y1 {
		//yay for tuple assignment
		y1, y2 = y2, y1
	}

	if y2 < 0 || y1 >= height || x < 0 || x >= width {
		return
	}

	if y1 < 0 {
		y1 = 0
	}

	SDLColor := sdl.MapRGBA(surface.Format, color.R, color.G, color.B, 255)

	for y := y1; y < y2 && y < height; y++ {
		pixels[(x + width*y*4)] = SDLColor
	}
}

func Done() bool {
	readKeys()
	if keys[sdl.SCANCODE_ESCAPE] != 0 {
		return true
	}

	return false
}

func Cls() {
	for i := range pixels {
		pixels[i] = 0
	}
}

func readKeys() {
	sdl.PumpEvents()
	keys = sdl.GetKeyboardState()
}

func keyDown(key int) bool {
	if keys == nil {
		return false
	}
	return (keys[key] != 0)
}

func genTextures() {
	for x := 0; x < texWidth; x++ {
		for y := 0; y < texHeight; y++ {
			ycolor := y * 256 / texHeight
			xycolor := y*128/texHeight + x*128/texWidth

			textures[0][texWidth*y+x] = 128 + 256*128 + 65536*0
			textures[1][texWidth*y+x] = uint32(xycolor) + 256*uint32(xycolor) + 65536*uint32(xycolor)
			textures[2][texWidth*y+x] = 65536 * 192 * (uint32(x) % 16)
			textures[3][texWidth*y+x] = 65536 * uint32(ycolor)
			textures[4][texWidth*y+x] = 128 + 256*128 + 65536*128
		}
	}
}

func loadTextures() error {
	err := loadTexture(0, "pics/eagle.png")
	if err != nil {
		return err
	}

	err = loadTexture(1, "pics/redbrick.png")
	if err != nil {
		return err
	}

	loadTexture(2, "pics/purplestone.png")
	if err != nil {
		return err
	}

	err = loadTexture(3, "pics/greystone.png")
	if err != nil {
		return err
	}

	err = loadTexture(4, "pics/bluestone.png")
	if err != nil {
		return err
	}

	err = loadTexture(5, "pics/mossy.png")
	if err != nil {
		return err
	}

	err = loadTexture(6, "pics/wood.png")
	if err != nil {
		return err
	}

	err = loadTexture(7, "pics/colorstone.png")
	if err != nil {
		return err
	}

	return nil
}

func loadTexture(texIndex int, fileName string) error {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: File could not be opened")
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		return err
	}

	for y := 0; y < texHeight; y++ {
		for x := 0; x < texWidth; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			val := sdl.MapRGBA(surface.Format, uint8(b), uint8(g), uint8(r), uint8(a))
			textures[texIndex][texWidth*y+x] = val
		}
	}
	return nil
}
