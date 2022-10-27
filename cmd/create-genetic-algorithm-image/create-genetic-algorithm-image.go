package main

import (
	"fmt"
	"genetic-algorithm-image/domain"
	"log"
	"path/filepath"
	"sort"
)

func main() {
	// initGen, err := domain.NewGeneration("gen0.txt")
	// if err != nil {
	// 	log.Fatal("Failed to NewGeneration: ", err)
	// }
	// initGen.InitGen()
	for i := 0; i < 5; i++ {
		uiSelect()
	}
	crossing()
}

func inputStdin() string {
	var inputValue string
	fmt.Println(">")
	fmt.Scan(&inputValue)
	return inputValue
}

// 対象のファイル名取得
func getTargetFileName(pattern string) string {
	fileNames, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	sort.Strings(fileNames)
	fileName := fileNames[len(fileNames)-1]
	return fileName
}

func uiSelect() {
	// 対象のgenファイル名取得
	genFileName := getTargetFileName("./gen*.txt")

	generation, err := domain.NewGeneration(genFileName)
	if err != nil {
		log.Fatal("Failed to NewGeneration: ", err)
	}

	// 対象のselectファイル名を取得
	selectFileName := getTargetFileName("./select*.txt")

	selectedFile, err := domain.NewSelectedFile(selectFileName)
	if err != nil {
		log.Fatal("Failed to NewSelectedNumbers: ", err)
	}

	// selectファイル名からphase取得
	phase, err := selectedFile.GetPhase()
	if err != nil {
		log.Fatal("Failed to GetPhase: ", err)
	}

	// selecteファイル読み込み
	selectedNumbers, err := selectedFile.GetData()
	if err != nil {
		log.Fatal("Failed to selectedNumbers GetData: ", err)
	}

	// 次のphaseのselectファイルを生成
	nextSelectFile, err := selectedFile.CreateNextFiile()
	if err != nil {
		log.Fatal("Failed to CreateNextFiile: ", err)
	}
	defer nextSelectFile.Close()

	for i := 0; i < len(selectedNumbers); i++ {
		if i%2 != 0 {
			continue
		}

		// 比較画像生成
		err := generation.CreateComparisonImage(selectedNumbers[i], selectedNumbers[i+1])
		if err != nil {
			log.Fatal("Failed to CreateComparisonImage: ", err)
		}

		// 画像選択（j:左、k:右）
		for {
			fmt.Printf("{%d}-{%d}", phase, i)
			selected := inputStdin()
			if selected == "j" {
				nextSelectFile.Write([]byte(fmt.Sprintf("%d\n", selectedNumbers[i])))
				break
			}
			if selected == "k" {
				nextSelectFile.Write([]byte(fmt.Sprintf("%d\n", selectedNumbers[i+1])))
				break
			}
		}
	}
}

func crossing() {
	// 対象のselectファイル名を取得
	selectFileName := getTargetFileName("./select*.txt")

	selectedFile, err := domain.NewSelectedFile(selectFileName)
	if err != nil {
		log.Fatal("Failed to NewSelectedFile: ", err)
	}

	// selecteファイル読み込み
	selectedNumbers, err := selectedFile.GetData()
	if err != nil {
		log.Fatal("Failed to selectedFile GetData: ", err)
	}

	// 対象のgenファイル名取得
	generationFileName := getTargetFileName("./gen*.txt")

	generationFile, err := domain.NewGeneration(generationFileName)
	if err != nil {
		log.Fatal("Failed to NewGeneration: ", err)
	}

	// フェーズごとの代表画像の生成
	err = generationFile.CreateRepresentImage(selectedNumbers)
	if err != nil {
		log.Fatal("Failed to CreateRepresentImage: ", err)
	}

	// 次の世代のgenファイルを生成
	err = generationFile.NextCreate(selectedNumbers)
	if err != nil {
		log.Fatal("Failed to NextCreate: ", err)
	}
	// 前世代のファイルを移動
	err = generationFile.MovingPrevious()
	if err != nil {
		log.Fatal("Failed to MovingPrevious: ", err)
	}
}
