# Matrix Multiplication con Procesos en Go y C
Realizado por Emmanuel Bustamante Valbuena & Sebastian alaya

# Multiplicación de Matrices con Procesos en C

Este repositorio contiene la implementación en C para multiplicar dos matrices grandes utilizando:

* **Modo Secuencial** (1 proceso)
* **Modo Paralelo** (K procesos OS) con comunicación vía **memoria compartida (shm)**

---

## Resultados de Ejecución en C

| Procesos (K) | Tiempo Secuencial (s) | Tiempo Paralelo (s) | Speedup |
| -----------: | --------------------: | ------------------: | ------: |
|            1 |              0.000013 |            0.000606 |   0.02× |
|            2 |              0.000013 |            0.000838 |   0.02× |
|            4 |              0.000013 |            0.001196 |   0.01× |
|            6 |              0.000138 |            0.001717 |   0.08× |
|            8 |              0.000013 |            0.001839 |   0.01× |

> **Nota:** El tiempo secuencial se midió en una sola ejecución. El tiempo paralelo incluye el overhead de `fork()` y sincronización.

---

## Conclusión

* Los resultados muestran speedups muy inferiores a 1 para todas las configuraciones, lo cual indica que el overhead de creación de procesos y el acceso a memoria compartida domina en matrices de pequeño tamaño.
* Para evidenciar mejoras reales de paralelismo a nivel de proceso, es indispensable usar matrices de mayor dimensión (por ejemplo, ≥500×500), donde el coste computacional supere el overhead de IPC.
* **Recomendación:** Generar datos sintéticos grandes, volver a medir y graficar el comportamiento de speedup vs. número de procesos.

---

## Uso

1. Preparar archivos de entrada `matriz_a.txt` y `matriz_b.txt`.
2. Compilar:

   ```bash
   gcc -o matrix_mul matrix_mul.c -lrt
   ```
3. Ejecutar en modo secuencial (por defecto K=4):

   ```bash
   ./matrix_mul
   ```
4. Ejecutar con K procesos:

   ```bash
   ./matrix_mul <K>
   ```
5. La matriz resultado se guarda en `matriz_c.txt`.


----

Este repositorio contiene una implementación en Go para multiplicar dos matrices grandes utilizando:

- **Modo Secuencial** (1 proceso)
- **Modo Paralelo** (K procesos OS) con comunicación vía **pipes**

---

## Resultados de Ejecución

| Procesos (K) | Tiempo Secuencial (s) | Tiempo Paralelo (s) | Speedup |
|-------------:|----------------------:|--------------------:|--------:|
|            1 |          0.000002795* |           0.000001476 |   0.53× |
|            2 |          0.000002882 |           0.000004880 |   0.59× |
|            4 |          0.000001476 |           0.000007147 |   0.21× |
|            6 |          0.000002202 |           0.000009401 |   0.23× |
|            8 |          0.000001498 |           0.000009037 |   0.17× |

*Nota: el valor de 1 proceso en paralelo se corresponde a `./matrix_mul` sin argumentos y a `./matrix_mul 1`, se promedia para comparación.

---

## Conclusión

- El **speedup** obtenido es menor que 1 en todos los casos, indicando que el overhead de creación de procesos y comunicación supera el beneficio computacional para matrices pequeñas.  
- Para matrices de mayor tamaño, se espera que la ventaja de paralelizar a nivel de procesos sea más evidente.  
- **Paso siguiente**: probar con matrices de dimensiones ≥500×500 y graficar velocidad vs. número de procesos para analizar la escalabilidad real.

---

## Uso

1. Genera archivos de entrada `matriz_a.txt` y `matriz_b.txt`.  
2. Compila:
   ```bash
   go build -o matrix_mul matrix_mul.go
   ```
3. Ejecuta en modo secuencial (por defecto K=4):
   ```bash
   ./matrix_mul
   ```
4. Ejecuta con K procesos:
   ```bash
   ./matrix_mul <K>
   ```
5. La matriz resultado se guarda en `matriz_c.txt`.

---