// matrix_mul.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// leerMatriz lee una matriz desde un archivo de texto.
func leerMatriz(ruta string) ([]float64, int, int) {
	f, err := os.Open(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error abriendo %s: %v\n", ruta, err)
		os.Exit(1)
	}
	defer f.Close()

	var data []float64
	filas, cols := 0, 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		campos := strings.Fields(line)
		if cols == 0 {
			cols = len(campos)
		} else if len(campos) != cols {
			fmt.Fprintf(os.Stderr, "Columnas inconsistentes en %s\n", ruta)
			os.Exit(1)
		}
		for _, s := range campos {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Valor inválido en %s: %v\n", ruta, err)
				os.Exit(1)
			}
			data = append(data, v)
		}
		filas++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo %s: %v\n", ruta, err)
		os.Exit(1)
	}
	return data, filas, cols
}

// escribirMatriz guarda la matriz en un archivo de texto.
func escribirMatriz(ruta string, C []float64, filas, cols int) {
	f, err := os.Create(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creando %s: %v\n", ruta, err)
		os.Exit(1)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for i := 0; i < filas; i++ {
		for j := 0; j < cols; j++ {
			fmt.Fprintf(w, "%f", C[i*cols+j])
			if j < cols-1 {
				w.WriteByte(' ')
			}
		}
		w.WriteByte('\n')
	}
	w.Flush()
}

// multSecuencial: multiplicación clásica.
func multSecuencial(A, B, C []float64, N, M, P int) {
	for i := 0; i < N; i++ {
		for j := 0; j < P; j++ {
			sum := 0.0
			for k := 0; k < M; k++ {
				sum += A[i*M+k] * B[k*P+j]
			}
			C[i*P+j] = sum
		}
	}
}

func main() {
	// Modo hijo: args[1]=="child", args[2]=start, args[3]=end
	if len(os.Args) > 1 && os.Args[1] == "child" {
		start, _ := strconv.Atoi(os.Args[2])
		end, _ := strconv.Atoi(os.Args[3])
		// Leer matrices (solo tamaños necesarios)
		A, _, M := leerMatriz("matriz_a.txt")
		B, M2, P := leerMatriz("matriz_b.txt")
		if M != M2 {
			fmt.Fprintf(os.Stderr, "Dimensiones no compatibles\n")
			os.Exit(1)
		}
		// Calcular subconjunto de filas
		Csub := make([]float64, (end-start)*P)
		for i := start; i < end; i++ {
			for j := 0; j < P; j++ {
				sum := 0.0
				for k := 0; k < M; k++ {
					sum += A[i*M+k] * B[k*P+j]
				}
				Csub[(i-start)*P+j] = sum
			}
		}
		// Enviar resultados parciales por stdout
		w := bufio.NewWriter(os.Stdout)
		for idx, val := range Csub {
			fmt.Fprintf(w, "%f", val)
			if (idx+1)%P == 0 {
				w.WriteByte('\n')
			} else {
				w.WriteByte(' ')
			}
		}
		w.Flush()
		return
	}

	// Modo padre
	// Número de procesos
	K := 4
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil && v > 0 {
			K = v
		}
	}
	// Leer matrices
	A, N, M := leerMatriz("matriz_a.txt")
	B, M2, P := leerMatriz("matriz_b.txt")
	if M != M2 {
		fmt.Fprintf(os.Stderr, "Error: columnas A (%d) != filas B (%d)\n", M, M2)
		os.Exit(1)
	}

	// Buffers de resultado
	Cseq := make([]float64, N*P)
	Cpar := make([]float64, N*P)

	// Secuencial
	start := time.Now()
	multSecuencial(A, B, Cseq, N, M, P)
	tSec := time.Since(start).Seconds()

	// Paralelo con procesos e IPC por pipes
	size := (N + K - 1) / K
	for g := 0; g < K; g++ {
		st := g * size
		en := st + size
		if en > N {
			en = N
		}
		cmd := exec.Command(os.Args[0], "child", strconv.Itoa(st), strconv.Itoa(en))
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Pipe error: %v\n", err)
			os.Exit(1)
		}
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Start error: %v\n", err)
			os.Exit(1)
		}
		// Leer resultados parciales
		s := bufio.NewScanner(pipe)
		row := st
		for s.Scan() {
			campos := strings.Fields(s.Text())
			for j, str := range campos {
				v, _ := strconv.ParseFloat(str, 64)
				Cpar[row*P+j] = v
			}
			row++
		}
		cmd.Wait()
	}
	// Guardar matriz resultado
	escribirMatriz("matriz_c.txt", Cpar, N, P)

	// Imprimir tiempos
	fmt.Printf("Sequential time: %.9f seconds\n", tSec)
	// Medir tiempo paralelo
	start = time.Now()
	for g := 0; g < K; g++ {
		st := g * size
		en := st + size
		if en > N {
			en = N
		}
		cmd := exec.Command(os.Args[0], "child", strconv.Itoa(st), strconv.Itoa(en))
		pipe, _ := cmd.StdoutPipe()
		cmd.Start()
		s := bufio.NewScanner(pipe)
		for s.Scan() { /* descartar */
		}
		cmd.Wait()
	}
	tPar := time.Since(start).Seconds()
	fmt.Printf("Parallel time (%d processes): %.9f seconds\n", K, tPar)
	fmt.Printf("Speedup: %.2fx\n", tSec/tPar)
}
