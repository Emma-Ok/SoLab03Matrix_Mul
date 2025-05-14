<!-- 🎉 Bienvenido al repositorio de Matrix Multiplication con Procesos en Go y C 🚀 -->

# ⚙️ Matrix Multiplication con Procesos en Go y C

**🔧 Autores:** Emmanuel Bustamante Valbuena & Sebastian Amaya Pérez

Este repositorio 🚧 compara dos implementaciones para multiplicar matrices grandes usando procesos del sistema operativo:

* 🟢 **C (Memoria Compartida)**
* 🐹 **Go (Pipes)**

---

## 📊 1. Resultados de Ejecución en C (Memoria Compartida)

| 🔢 Procesos (K) | ⏱️ Secuencial (s) | ⏱️ Paralelo (s) | 📈 Speedup |
| --------------: | ----------------: | --------------: | ---------: |
|               1 |          0.000013 |        0.000606 |      0.02× |
|               2 |          0.000013 |        0.000838 |      0.02× |
|               4 |          0.000013 |        0.001196 |      0.01× |
|               6 |          0.000138 |        0.001717 |      0.08× |
|               8 |          0.000013 |        0.001839 |      0.01× |

> 💡 **Nota:** El tiempo paralelo incluye overhead de `fork()` y sincronización.

---

## 📊 2. Resultados de Ejecución en Go (Pipes)

| 🔢 Procesos (K) | ⏱️ Secuencial (s) | ⏱️ Paralelo (s) | 📈 Speedup |
| --------------: | ----------------: | --------------: | ---------: |
|               1 |    0.000002795 \* |     0.000001476 |      0.53× |
|               2 |       0.000002882 |     0.000004880 |      0.59× |
|               4 |       0.000001476 |     0.000007147 |      0.21× |
|               6 |       0.000002202 |     0.000009401 |      0.23× |
|               8 |       0.000001498 |     0.000009037 |      0.17× |

\*⭐ *K=1 es promedio de `./matrix_mul` y `./matrix_mul 1`.*

---

## 🔍 3. Análisis Comparativo

1. **Overhead de IPC**

   * 🛠️ C: `shmget`/`shmat` añade coste fijo, pero evita copias extra.
   * 🐹 Go: Pipes requieren serialización/parsing de texto → overhead elevado.

2. **Speedup vs Tamaño**

   * Ambos muestran **speedup < 1** para matrices pequeñas (\~µs).
   * Go 📈 (0.53× en K=1) supera a C (0.02×) en matrices mínimas.

3. **Escalabilidad**

   * Más procesos (K↑) **no mejora**; overhead crece más rápido que computo.

4. **Recomendaciones**

   * 🖥️ Probar con matrices ≥500×500 para amortizar IPC.
   * 🐹 Go: considerar compartir memoria (`golang.org/x/sys/unix`).
   * 🛠️ C: comparar pipes vs shm según entorno.

---

## 🎯 4. Conclusión

> Las implementaciones actuales son **prueba de concepto**. Para medir paralelismo real:
>
> * ➡️ Generar matrices grandes.
> * 📈 Graficar speedup vs K.
> * 🔄 Evaluar distintos métodos de IPC.

---

## 🚀 5. ¡Manos a la Obra! Uso

### 🛠️ En C (Memoria Compartida)

```bash
gcc -o matrix_mul matrix_mul.c -lrt
./matrix_mul <K>
```

### 🐹 En Go (Pipes)

```bash
go build -o matrix_mul matrix_mul.go
./matrix_mul <K>
```

💾 **Archivos necesarios:** `matriz_a.txt`, `matriz_b.txt`.
📂 **Salida:** `matriz_c.txt`.

---


![Performance](https://img.shields.io/badge/Benchmarking-Experimental-orange) ![License](https://img.shields.io/badge/License-MIT-green)
