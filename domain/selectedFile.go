package domain

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SelectedFile struct {
	Name string
}

func NewSelectedFile(fileName string) (*SelectedFile, error) {
	if strings.TrimSpace(fileName) == "" {
		return nil, errors.New("required is fileName")
	}
	selectedFile := new(SelectedFile)
	selectedFile.Name = fileName
	return selectedFile, nil
}

// selecteファイルデータ取得
func (s SelectedFile) GetData() ([]int, error) {
	// selectファイル読み込み
	selectFile, err := os.Open(s.Name)
	if err != nil {
		return nil, err
	}
	defer selectFile.Close()

	scanner := bufio.NewScanner(selectFile)
	numbers := []int{}
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, num)
	}
	return numbers, nil
}

// phase番号取得
func (s SelectedFile) GetPhase() (int, error) {
	no, err := strconv.Atoi(s.Name[6:7])
	if err != nil {
		return 0, err
	}
	phase := no + 1

	return phase, nil
}

func (s SelectedFile) CreateNextFiile() (*os.File, error) {
	phase, err := s.GetPhase()
	if err != nil {
		return nil, err
	}
	nextSelectFile, err := os.Create(fmt.Sprintf("select%d.txt", phase))
	if err != nil {
		return nil, err
	}

	return nextSelectFile, nil
}
