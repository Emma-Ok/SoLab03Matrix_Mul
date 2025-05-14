# Matrix Multiplication con Procesos en Go y C
Realizado por Emmanuel Bustamante Valbuena & Sebastian alaya

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

¡Bienvenido(a) a contribuir issues o pull requests para mejoras!
