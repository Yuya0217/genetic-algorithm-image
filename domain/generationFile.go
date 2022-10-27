package domain

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type GenerationFile struct {
	Name string
}

type string_slice_t []string

var (
	R      = []uint8{30, 255, 255, 46, 245}
	G      = []uint8{144, 255, 0, 139, 222}
	B      = []uint8{255, 255, 0, 87, 179}
	width  = 16
	height = 16
)

func NewGeneration(fileName string) (*GenerationFile, error) {
	if strings.TrimSpace(fileName) == "" {
		return nil, errors.New("required is fileName")
	}
	generationFile := new(GenerationFile)
	generationFile.Name = fileName
	return generationFile, nil
}

func (ss string_slice_t) uint() []uint8 {
	f := make([]uint8, len(ss))
	for i, v := range ss {
		a, _ := strconv.Atoi(v)
		f[i] = uint8(a)
	}
	return f
}

func setField() []int {
	rand.Seed(time.Now().UnixNano())
	field := make([]int, width*height)
	for i := 0; i < len(field); i++ {

		field[i] = rand.Intn(5)
	}
	return field
}

// genファイル読み込み
func (gen GenerationFile) GetData() []string {
	bytes, err := ioutil.ReadFile(gen.Name)
	if err != nil {
		panic(err)
	}

	genFileData := strings.Replace(string(bytes), "[", "", -1)
	genFileData = strings.Replace(genFileData, "]", "", -1)
	genFileDatas := strings.Split(genFileData, "\n")
	return genFileDatas
}

// genファイル名から世代番号を取得
func (gen GenerationFile) GetFileNo() int {
	generationNo, _ := strconv.Atoi(gen.Name[3:4])
	return generationNo
}

// genデータ加工
func (gen GenerationFile) GenDataFormat(genFileData string) []uint8 {
	return string_slice_t(strings.Split(genFileData, " ")).uint()
}

// genデータ初期作成
func (gen GenerationFile) InitGen() {
	// ファイルの生成
	fp, err := os.Create(gen.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()

	for i := 0; i < 128; i++ {
		r, g, b := []uint8{}, []uint8{}, []uint8{}
		field := setField()
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				index := field[y*height+x]
				r = append(r, R[index])
				g = append(g, G[index])
				b = append(b, B[index])
			}
		}
		fp.WriteString(fmt.Sprintf("%d\n%d\n%d\n", r, g, b))
	}
}

// 次の世代のデータを作成
func (gen GenerationFile) NextCreate(numbers []int) error {
	fp, err := os.Create(fmt.Sprintf("./gen%d.txt", gen.GetFileNo()+1))
	if err != nil {
		return err
	}
	defer fp.Close()

	generationFileDatas := gen.GetData()

	var cnt = 0
	for i := 0; i < 128; i++ {
		cnt = 0
		var r, g, b = []uint8{}, []uint8{}, []uint8{}
		// 最終的に選んだ4つの親を次の世代のデータに含める
		a := i % 31
		if a == 0 && i != 0 {
			index := (i / 31) - 1
			rGenIndex := numbers[index] * 3
			rGenerationFileData := generationFileDatas[rGenIndex]
			rGenerationFileDataUint := gen.GenDataFormat(rGenerationFileData)
			r = rGenerationFileDataUint

			gGenIndex := numbers[index]*3 + 1
			gRedGenerationFileData := generationFileDatas[gGenIndex]
			gRedGenerationFileDataUint := gen.GenDataFormat(gRedGenerationFileData)
			g = gRedGenerationFileDataUint

			bGenIndex := numbers[index]*3 + 2
			bRedGenerationFileData := generationFileDatas[bGenIndex]
			bRedGenerationFileDataUint := gen.GenDataFormat(bRedGenerationFileData)
			b = bRedGenerationFileDataUint
			_, err = fp.WriteString(fmt.Sprintf("%d\n%d\n%d\n", r, g, b))
			if err != nil {
				return err
			}
			continue
		}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// 突然変異 (8/128の確率)
				rand.Seed(time.Now().UnixNano())
				randomNumber := rand.Intn(127)
				isMutation := randomNumber < 8
				if isMutation {
					field := setField()
					index := field[y*height+x]
					r = append(r, R[index])
					g = append(g, G[index])
					b = append(b, B[index])
					cnt += 1
					continue
				}
				// 4つからランダムに選んで交配
				rand.Seed(time.Now().UnixNano())
				index := rand.Intn(4)
				rGenIndex := numbers[index] * 3
				rGenerationFileData := generationFileDatas[rGenIndex]
				rGenerationFileDataUint := gen.GenDataFormat(rGenerationFileData)
				r = append(r, rGenerationFileDataUint[cnt])

				gGenIndex := numbers[index]*3 + 1
				gRedGenerationFileData := generationFileDatas[gGenIndex]
				gRedGenerationFileDataUint := gen.GenDataFormat(gRedGenerationFileData)
				g = append(g, gRedGenerationFileDataUint[cnt])

				bGenIndex := numbers[index]*3 + 2
				bRedGenerationFileData := generationFileDatas[bGenIndex]
				bRedGenerationFileDataUint := gen.GenDataFormat(bRedGenerationFileData)
				b = append(b, bRedGenerationFileDataUint[cnt])
				cnt += 1
			}
		}
		_, err = fp.WriteString(fmt.Sprintf("%d\n%d\n%d\n", r, g, b))
		if err != nil {
			return err
		}
	}
	return nil
}

// 前の世代データを退避
func (gen GenerationFile) MovingPrevious() error {
	mkDirName := fmt.Sprintf("gen%d", gen.GetFileNo())
	if err := os.Mkdir(mkDirName, 0777); err != nil {
		return err
	}
	if err := os.Rename(gen.Name, fmt.Sprintf("./%s/%s", mkDirName, gen.Name)); err != nil {
		return err
	}
	if err := os.Rename("./select1.txt", fmt.Sprintf("./%s/%s", mkDirName, "select1.txt")); err != nil {
		return err
	}
	if err := os.Rename("./select2.txt", fmt.Sprintf("./%s/%s", mkDirName, "select2.txt")); err != nil {
		return err
	}
	if err := os.Rename("./select3.txt", fmt.Sprintf("./%s/%s", mkDirName, "select3.txt")); err != nil {
		return err
	}
	if err := os.Rename("./select4.txt", fmt.Sprintf("./%s/%s", mkDirName, "select4.txt")); err != nil {
		return err
	}
	if err := os.Rename("./select5.txt", fmt.Sprintf("./%s/%s", mkDirName, "select5.txt")); err != nil {
		return err
	}
	return nil
}

// 代表画像の生成
func (gen GenerationFile) CreateRepresentImage(numbers []int) error {
	// genファイル読み込み
	generationFileDatas := gen.GetData()

	number := numbers[0]
	r := gen.GenDataFormat(generationFileDatas[number*3])
	g := gen.GenDataFormat(generationFileDatas[number*3+1])
	b := gen.GenDataFormat(generationFileDatas[number*3+2])

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	cnt := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{r[cnt], g[cnt], b[cnt], 255})
			cnt += 1
		}
	}

	outImageFile, err := os.Create(fmt.Sprintf("./%d.png", gen.GetFileNo()))
	if err != nil {
		return err
	}
	defer outImageFile.Close()

	err = png.Encode(outImageFile, img)
	if err != nil {
		return err
	}
	return nil
}

// 比較画像の生成
func (gen GenerationFile) CreateComparisonImage(number1 int, number2 int) error {
	genFileDatas := gen.GetData()
	r1 := gen.GenDataFormat(genFileDatas[number1*3])
	g1 := gen.GenDataFormat(genFileDatas[number1*3+1])
	b1 := gen.GenDataFormat(genFileDatas[number1*3+2])
	r2 := gen.GenDataFormat(genFileDatas[number2*3])
	g2 := gen.GenDataFormat(genFileDatas[number2*3+1])
	b2 := gen.GenDataFormat(genFileDatas[number2*3+2])

	img1 := image.NewRGBA(image.Rect(0, 0, width, height))
	img2 := image.NewRGBA(image.Rect(0, 0, width, height))

	cnt := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img1.Set(x, y, color.RGBA{r1[cnt], g1[cnt], b1[cnt], 255})
			img2.Set(x, y, color.RGBA{r2[cnt], g2[cnt], b2[cnt], 255})
			cnt += 1
		}
	}

	outImg := image.NewRGBA(image.Rect(0, 0, width*2, height))
	draw.Draw(outImg, image.Rect(0, 0, width, height), img1, image.Point{0, 0}, draw.Over)
	draw.Draw(outImg, image.Rect(width+1, 0, width*2, height), img2, image.Point{0, 0}, draw.Over)

	outImageFile, err := os.Create("./image.png")
	if err != nil {
		return err
	}
	defer outImageFile.Close()

	err = png.Encode(outImageFile, outImg)
	if err != nil {
		return err
	}
	return nil
}
