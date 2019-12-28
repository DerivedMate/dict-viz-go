package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"os"
	"strconv"

	"github.com/fogleman/gg"
)

func Length(arr *[][]string, i int) int {
	return len((*arr)[i][0])
}

type Entry struct {
	A       string
	Cosines []float64
}

type Msg struct {
	a []string
	i int
}

func DoMapping(arr *[][]string, in chan Msg, out chan Entry) {
	for m := range in {
		a, i := m.a, m.i

		// 	aCategoriesInt := MapToInt(a[1:])
		e := Entry{
			A:       a[0],
			Cosines: []float64{},
		}
		for j := i + 1; j < len(*arr); j++ {
			b := (*arr)[j]

			// bCategoriesInt := MapToInt(b[1:])
			cosLetters := CosStrings(a[0], b[0])
			// cosCategories := CosVecs(aCategoriesInt, bCategoriesInt)
			p := 100000000.0

			cos := cosLetters //* cosCategories
			cos = math.Round(cos*p) / p

			e.Cosines = append(e.Cosines, cos)
		}
		fmt.Printf("%v : %v\n", i, a[0])
		out <- e
	}
	close(out)
}

func saveEntry(w *bufio.Writer, entry *Entry, m, n int) {
	j, err := json.Marshal(entry.Cosines)
	var stringified string
	if m == n {
		stringified = fmt.Sprintf("\"%v\": %v\n", entry.A, string(j))
	} else {
		stringified = fmt.Sprintf("\"%v\": %v,\n", entry.A, string(j))
	}
	w.WriteString(stringified)
	HandleError(err)
}

func GenRecords() {
	f, err := os.Open("out-clean.csv")
	HandleError(err)
	reader := bufio.NewReader(f)
	rCSV := csv.NewReader(reader)
	entries, err := rCSV.ReadAll()
	HandleError(err)
	offset := 15000
	count := 500
	if len(os.Args) >= 3 {
		s, err := strconv.Atoi(os.Args[2])
		HandleError(err)
		offset = s
	}
	if len(os.Args) >= 4 {
		s, err := strconv.Atoi(os.Args[3])
		HandleError(err)
		count = s
	}
	entries = entries[offset+1 : offset+count+2]
	n := len(entries)

	in := make(chan Msg)
	out := make(chan Entry)
	outFile, err := os.Create(fmt.Sprintf("vertices_%v+%v.json", offset, count))
	HandleError(err)
	w := bufio.NewWriter(outFile)
	w.WriteString("{\n")

	m := 0
	go func() {
		for i, a := range entries {
			in <- Msg{a, i}
		}
		close(in)
	}()
	go DoMapping(&entries, in, out)

	for e := range out {
		m++
		saveEntry(w, &e, m, n)
	}
	w.WriteString("}\n")
	w.Flush()
}

func Draw() {
	offset, count := os.Args[2], os.Args[3]
	f, err := os.Open(fmt.Sprintf("vertices_%v+%v.json", offset, count))
	defer f.Close()
	HandleError(err)

	// Let's be lazy for now
	var entries map[string][]float64
	j := json.NewDecoder(f)
	err = j.Decode(&entries)
	HandleError(err)

	s := 6500.0
	w := s
	h := s
	c := gg.NewContext(int(w), int(h))
	c.DrawRectangle(0.0, 0.0, w, h)
	c.SetHexColor("#1b1b1b")
	c.Fill()

	n := len(entries)
	dth := math.Pi * 2.0 / float64(n)
	r := w / 3.5
	// rSq := r * r
	c.Translate(w/2, h/2)
	c.SetLineWidth(0.07)
	i := 0
	// r0, r1 := 165.0, 220.0
	// r0, r1 := 140.0, 200.0
	// r0, r1 := 0.0, 300.0
	// s0, s1 := 0.9, 1.0
	a0, a1 := 0.0, 180.0
	// dr := r1 - r0
	// kr := 0.8

	for k, v := range entries {
		k0 := float64(i)
		th := dth * k0
		x0 := r * math.Cos(th)
		y0 := r * math.Sin(th)
		c.Push()
		c.RotateAbout(th, x0, y0)
		c.SetColor(color.White)
		c.DrawString(fmt.Sprintf(" %v", k), x0, y0)
		c.Pop()
		for j, cos := range v {
			k1 := float64(i + j + 1)
			x1 := r * math.Cos(dth*k1)
			y1 := r * math.Sin(dth*k1)

			// red, g, b, _ := colorful.Hsl(cos*dr+r0, 1, 0.5).RGBA()
			co := uint8(255.0 * cos)
			red, g, b := co, co, co
			c.SetColor(color.RGBA{R: red, G: g, B: b, A: uint8(cos*(a1-a0) + a0)})

			phi := (k1 + k0) * dth / 2.0
			x2, y2 := r*math.Cos(phi), r*math.Sin(phi)
			c.CubicTo(
				x0, y0,
				x2, y2,
				x1, y1,
			)
			// c.DrawLine(x0, y0, x1, y1)
			c.Stroke()
		}
		i++
		fmt.Printf("%v/%v\n", i, n)
	}

	err = c.SavePNG(fmt.Sprintf("out%v+%v.png", offset, count))
	HandleError(err)
}

func main() {
	VectorizationMain()
	/*
		@TODO:
		1. Try sorted runs: by length, alphabetically
	*/
	/*
		f := false
		if len(os.Args) > 1 {
			switch os.Args[1] {
			case "gen", "Gen", "0":
				GenRecords()
			case "draw", "Draw", "1":
				Draw()
			default:
				fmt.Printf("Unknown command \"%v\"", os.Args[1])
				f = true
			}
		}
		if f || len(os.Args) == 1 {
			fmt.Println("Usage: go run . <mode> <offset> <count>\n  mode: gen/draw\n  offset: the number of words to skip\n  count: the number of words to take")
		}
	*/
}

/*

var HasComa *regexp.Regexp = regexp.MustCompile(",")

func main_() {
	inF, err := os.Open("out2.csv")
	HandleError(err)
	outF, err := os.Create("out-clean.csv")
	HandleError(err)

	defer inF.Close()
	defer outF.Close()

	in := bufio.NewReader(inF)
	out := csv.NewWriter(outF)

	correctN := 0
	first := true

	for {
		b, _, err := in.ReadLine()
		line := string(b)
		l := strings.Split(line, ",")
		processed := []string{""}

		if first {
			correctN = len(l)
			first = false
		}

		for i, e := range l {
			if i <= len(l)-correctN {
				if processed[0] == "" {
					processed[0] = e
				} else {
					processed[0] = fmt.Sprintf("%v,%v", processed[0], e)
				}
			} else {
				processed = append(processed, e)
			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		out.Write(processed)
	}

	out.Flush()
}*/
