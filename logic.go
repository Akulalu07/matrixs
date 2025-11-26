package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func NewNum(a, b int) Num {
	if b == 0 {
		panic("denominator cannot be zero")
	}
	return normalize(Num{a, b})
}

func (n Num) Print() {
	if n.b == 1 {
		fmt.Print(n.a)
	} else {
		fmt.Printf("%d/%d", n.a, n.b)
	}
}

func (n Num) String() string {
	if n.b == 1 {
		return fmt.Sprintf("%d", n.a)
	} else {
		return fmt.Sprintf("%d/%d", n.a, n.b)
	}
}

func (n Num) Println() {
	if n.b == 1 {
		fmt.Println(n.a)
	} else {
		fmt.Printf("%d/%d\n", n.a, n.b)
	}
}

func Add(a, b Num) Num {
	return normalize(Num{
		a: a.a*b.b + b.a*a.b,
		b: a.b * b.b,
	})
}

func Mul(a, b Num) Num {
	return normalize(Num{
		a: a.a * b.a,
		b: a.b * b.b,
	})
}

func normalize(n Num) Num {
	g := gcd(abs(n.a), abs(n.b))
	n.a /= g
	n.b /= g
	if n.b < 0 {
		n.a = -n.a
		n.b = -n.b
	}
	return n
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func New() Matrix {
	return Matrix{
		lenm:   0,
		lenk:   0,
		matrix: make([][]Num, 0),
	}
}

func FromIntMatrix(IntMat [][]int) Matrix {
	lenm := len(IntMat)
	lenk := len(IntMat[0])
	M := make([][]Num, lenm)

	for i := 0; i < lenm; i++ {
		row := make([]Num, lenk)
		for j := 0; j < lenk; j++ {
			row[j] = NewNum(IntMat[i][j], 1)
		}
		M[i] = row
	}
	return Matrix{lenm, lenk, M}
}

func (m *Matrix) Append(row []Num) error {
	if m.lenm == 0 {
		m.lenk = len(row)
	} else if len(row) != m.lenk {
		return ErrIrregularMatrix
	}
	m.matrix = append(m.matrix, row)
	m.lenm++
	return nil
}

func (m Matrix) T() Matrix {
	if m.lenm == 0 || m.lenk == 0 {
		return New()
	}

	res := make([][]Num, m.lenk)
	for i := 0; i < m.lenk; i++ {
		res[i] = make([]Num, m.lenm)
		for j := 0; j < m.lenm; j++ {
			res[i][j] = m.matrix[j][i]
		}
	}
	return Matrix{m.lenk, m.lenm, res}
}

func (m *Matrix) Print() {
	for i := 0; i < m.lenm; i++ {
		for j := 0; j < m.lenk; j++ {
			n := m.matrix[i][j]
			n.Print()
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
func (m *Matrix) String() string {
	resp := ""
	for i := 0; i < m.lenm; i++ {
		for j := 0; j < m.lenk; j++ {
			if j != 0 {
				resp += ", "
			}
			n := m.matrix[i][j]
			resp += n.String()
		}
		resp += ";\n"
	}
	return resp
}

func (m *Matrix) SwapRows(i, j int) {
	m.matrix[i], m.matrix[j] = m.matrix[j], m.matrix[i]
}

func (m *Matrix) ScaleRow(i int, k Num) {
	for col := 0; col < m.lenk; col++ {
		m.matrix[i][col] = Mul(m.matrix[i][col], k)
	}
}

func (m *Matrix) MulMatrixNum(a Num) {
	n := m.lenm
	k := m.lenk
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			m.matrix[i][j] = Mul(m.matrix[i][j], a)
		}
	}
}
func AddMatrix(a Matrix, b Matrix) (Matrix, error) {
	if a.lenm != b.lenm && a.lenk != b.lenk {
		return a, errors.New("not")
	}
	m := a.Copy()
	for i := 0; i < a.lenm; i++ {
		for j := 0; j < a.lenk; j++ {
			m.matrix[i][j] = Add(a.matrix[i][j], b.matrix[i][j])
		}
	}
	return m, nil
}

func (m *Matrix) AddRow(dst, src int, k Num) {
	for col := 0; col < m.lenk; col++ {
		toAdd := Mul(m.matrix[src][col], k)
		m.matrix[dst][col] = Add(m.matrix[dst][col], toAdd)
	}
}

func (m *Matrix) SwapCols(i, j int) {
	for row := 0; row < m.lenm; row++ {
		m.matrix[row][i], m.matrix[row][j] = m.matrix[row][j], m.matrix[row][i]
	}
}

func (m Matrix) Copy() Matrix {
	cp := make([][]Num, m.lenm)
	for i := range m.matrix {
		cp[i] = make([]Num, m.lenk)
		copy(cp[i], m.matrix[i])
	}
	return Matrix{m.lenm, m.lenk, cp}
}
func (m Matrix) Det() (Num, error) {
	if m.lenm != m.lenk {
		return NewNum(0, 1), ErrSquare
	}

	n := m.lenm
	a := m.Copy()

	det := NewNum(1, 1)

	for col := 0; col < n; col++ {
		pivot := col
		for pivot < n && a.matrix[pivot][col].a == 0 {
			pivot++
		}

		if pivot == n {
			return NewNum(0, 1), nil
		}

		if pivot != col {
			a.SwapRows(pivot, col)
			det = Mul(det, NewNum(-1, 1))
		}

		pivotVal := a.matrix[col][col]
		det = Mul(det, pivotVal)

		inv := NewNum(1, 1)
		inv.a = pivotVal.b
		inv.b = pivotVal.a
		inv = normalize(inv)

		a.ScaleRow(col, inv)

		for row := col + 1; row < n; row++ {
			if a.matrix[row][col].a != 0 {
				k := NewNum(-a.matrix[row][col].a, a.matrix[row][col].b)
				a.AddRow(row, col, k)
			}
		}
	}

	return det, nil
}

func MulMatrix(a, b Matrix) (Matrix, error) {
	if a.lenk != b.lenm {
		return Matrix{}, ErrIrregularMatrix
	}

	res := make([][]Num, a.lenm)
	for i := 0; i < a.lenm; i++ {
		res[i] = make([]Num, b.lenk)
		for j := 0; j < b.lenk; j++ {
			sum := NewNum(0, 1)
			for p := 0; p < a.lenk; p++ {
				sum = Add(sum, Mul(a.matrix[i][p], b.matrix[p][j]))
			}
			res[i][j] = sum
		}
	}

	return Matrix{a.lenm, b.lenk, res}, nil
}

func (m Matrix) Inverse() (Matrix, error) {
	if m.lenm != m.lenk {
		return Matrix{}, ErrIrregularMatrix
	}

	n := m.lenm
	A := m.Copy()
	I := NewIdentity(n)

	for col := 0; col < n; col++ {

		pivot := col
		for pivot < n && A.matrix[pivot][col].a == 0 {
			pivot++
		}

		if pivot == n {
			return Matrix{}, errors.New("matrix is singular")
		}

		if pivot != col {
			A.SwapRows(pivot, col)
			I.SwapRows(pivot, col)
		}

		pivotVal := A.matrix[col][col]

		inv := NewNum(pivotVal.b, pivotVal.a)
		inv = normalize(inv)

		A.ScaleRow(col, inv)
		I.ScaleRow(col, inv)

		for row := 0; row < n; row++ {
			if row != col {
				if A.matrix[row][col].a != 0 {
					k := NewNum(-A.matrix[row][col].a, A.matrix[row][col].b)

					A.AddRow(row, col, k)
					I.AddRow(row, col, k)
				}
			}
		}
	}

	return I, nil
}

func NewIdentity(n int) Matrix {
	m := New()
	m.lenk = n
	m.lenm = n
	m.matrix = make([][]Num, n)
	for i := 0; i < n; i++ {
		row := make([]Num, n)
		for j := 0; j < n; j++ {
			if i == j {
				row[j] = NewNum(1, 1)
			} else {
				row[j] = NewNum(0, 1)
			}
		}
		m.matrix[i] = row
	}
	return m
}

func (m Matrix) REF() Matrix {
	A := m.Copy()
	rows := A.lenm
	cols := A.lenk

	r := 0
	for c := 0; c < cols && r < rows; c++ {
		pivot := r
		for pivot < rows && A.matrix[pivot][c].a == 0 {
			pivot++
		}

		if pivot == rows {
			continue
		}

		if pivot != r {
			A.SwapRows(pivot, r)
		}

		pivotVal := A.matrix[r][c]
		inv := NewNum(pivotVal.b, pivotVal.a)
		inv = normalize(inv)
		A.ScaleRow(r, inv)

		for row := r + 1; row < rows; row++ {
			if A.matrix[row][c].a != 0 {
				k := NewNum(-A.matrix[row][c].a, A.matrix[row][c].b)
				A.AddRow(row, r, k)
			}
		}
		r++
	}

	return A
}

func (m Matrix) RREF() Matrix {
	A := m.REF()

	rows := A.lenm
	cols := A.lenk

	for i := rows - 1; i >= 0; i-- {

		pivotCol := -1
		for j := 0; j < cols; j++ {
			if A.matrix[i][j].a != 0 {
				pivotCol = j
				break
			}
		}
		if pivotCol == -1 {
			continue
		}

		for up := i - 1; up >= 0; up-- {
			if A.matrix[up][pivotCol].a != 0 {
				k := NewNum(-A.matrix[up][pivotCol].a, A.matrix[up][pivotCol].b)
				A.AddRow(up, i, k)
			}
		}
	}

	return A
}

func (m Matrix) Rank() int {
	ref := m.REF()

	rank := 0
	for i := 0; i < ref.lenm; i++ {
		zero := true
		for j := 0; j < ref.lenk; j++ {
			if ref.matrix[i][j].a != 0 {
				zero = false
				break
			}
		}
		if !zero {
			rank++
		}
	}
	return rank
}

func (s *Script) WriteToFile() error {
	outDir := "outputs"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create outputs dir: %w", err)
	}
	base := filepath.Base(s.path)
	scriptPath := filepath.Join(outDir, base)
	pdfName := strings.TrimSuffix(base, filepath.Ext(base)) + ".pdf"
	pdfPath := filepath.Join(outDir, pdfName)
	if err := os.Remove(pdfPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove old pdf: %w", err)
	}
	f, err := os.OpenFile(scriptPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open typ file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(s.s + "\n"); err != nil {
		return fmt.Errorf("write typ file: %w", err)
	}

	cmd := exec.Command("typst", "compile", scriptPath, pdfPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("typst compile failed: %w", err)
	}
	return nil
}

func NewScript(nameFile string) Script {
	return Script{
		"", nameFile + ".typ",
	}
}

func (s *Script) Append(some string) {
	s.s += some
}
func (s *Script) Appendln(some string) {
	s.s += some + "\n"
}
func (s *Script) AddMatrix(m Matrix, name string) {
	s.Append(fmt.Sprintf("\n%s = mat( %s )\n", name, m.String()))
}

func (s *Script) AddMatrixln(m Matrix, name string) {
	s.Append(fmt.Sprintf("\n%s = mat( %s )\n\\\n", name, m.String()))
}

func (s *Script) AddEqualMatrix(m Matrix) {
	s.Append(fmt.Sprintf("~\nmat( %s )\n", m.String()))
}
