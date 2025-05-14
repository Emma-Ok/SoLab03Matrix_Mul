// matrix_mul.c
// Laboratorio 3 - Sistemas Operativos
// Práctica 3: Multiplicación de matrices con procesos
// Autores: Emmanuel Bustamante Valbuena & Sebastian Amaya Pérez

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/shm.h>
#include <sys/wait.h>
#include <sys/types.h>
#include <unistd.h>
#include <time.h>

#define MIN(a,b) ((a)<(b)?(a):(b))

double* leer_matriz(const char* nombre, int* filas, int* cols) {
    FILE* f = fopen(nombre, "r");
    if (!f) {
        perror("Error abriendo archivo");
        exit(EXIT_FAILURE);
    }
    // Contar filas y columnas
    int capacity = 1024;
    double* data = malloc(sizeof(double) * capacity);
    *filas = 0;
    *cols = 0;
    char line[8192];
    while (fgets(line, sizeof(line), f)) {
        char* ptr = line;
        int cnt = 0;
        while (*ptr) {
            double val;
            int read;
            if (sscanf(ptr, "%lf%n", &val, &read) == 1) {
                if (*filas * (*cols == 0 ? 1 : *cols) + cnt >= capacity) {
                    capacity *= 2;
                    data = realloc(data, sizeof(double) * capacity);
                }
                data[*filas * (*cols == 0 ? cnt+1 : *cols) + cnt] = val;
                cnt++;
                ptr += read;
            } else break;
            while (*ptr==' '||*ptr=='\t') ptr++;
        }
        if (cnt == 0) continue;
        if (*cols == 0) *cols = cnt;
        else if (cnt != *cols) {
            fprintf(stderr, "Columnas inconsistentes en %s, linea %d\n", nombre, *filas+1);
            exit(EXIT_FAILURE);
        }
        (*filas)++;
    }
    fclose(f);
    return data;
}

void escribir_matriz(const char* nombre, double* C, int filas, int cols) {
    FILE* f = fopen(nombre, "w");
    if (!f) {
        perror("Error creando archivo");
        exit(EXIT_FAILURE);
    }
    for (int i = 0; i < filas; i++) {
        for (int j = 0; j < cols; j++) {
            fprintf(f, "%f", C[i*cols + j]);
            if (j < cols-1) fputc(' ', f);
        }
        fputc('\n', f);
    }
    fclose(f);
}

double tiempo_seg(struct timespec start, struct timespec end) {
    return (end.tv_sec - start.tv_sec) + (end.tv_nsec - start.tv_nsec)/1e9;
}

int main(int argc, char* argv[]) {
    int K = 4;
    if (argc == 2) {
        K = atoi(argv[1]);
        if (K < 1) K = 1;
    }
    int N, M, M2, P;
    double* A = leer_matriz("matriz_a.txt", &N, &M);
    double* B = leer_matriz("matriz_b.txt", &M2, &P);
    if (M != M2) {
        fprintf(stderr, "Dimensiones incompatibles: A cols %d != B filas %d\n", M, M2);
        exit(EXIT_FAILURE);
    }
    // Crear memoria compartida
    int shA = shmget(IPC_PRIVATE, sizeof(double) * N * M, IPC_CREAT | 0666);
int shB = shmget(IPC_PRIVATE, sizeof(double) * M * P, IPC_CREAT | 0666);
int shC = shmget(IPC_PRIVATE, sizeof(double) * N * P, IPC_CREAT | 0666);

double* shmpA = shmat(shA, NULL, 0);
double* shmpB = shmat(shB, NULL, 0);
double* shmpC = shmat(shC, NULL, 0);

memcpy(shmpA, A, sizeof(double) * N * M);
memcpy(shmpB, B, sizeof(double) * M * P);
memset(shmpC, 0, sizeof(double) * N * P);


    struct timespec t0, t1;
    // Secuencial
    clock_gettime(CLOCK_MONOTONIC, &t0);
    for (int i = 0; i < N; i++)
        for (int j = 0; j < P; j++) {
            double sum=0;
            for (int k = 0; k < M; k++)
                sum += shmpA[i * M + k] * shmpB[k * P + j];
            shmpC[i*P + j] = sum;
        }
    clock_gettime(CLOCK_MONOTONIC, &t1);
    double t_seq = tiempo_seg(t0, t1);

    // Paralelo
    memset(shmpC, 0, sizeof(double)NP);
    int filas_pp = (N + K - 1)/K;
    clock_gettime(CLOCK_MONOTONIC, &t0);
    for (int p = 0; p < K; p++) {
        pid_t pid = fork();
        if (pid == 0) {
            int start = p * filas_pp;
            int end = MIN(N, start + filas_pp);
            for (int i = start; i < end; i++)
                for (int j = 0; j < P; j++) {
                    double sum=0;
                    for (int k = 0; k < M; k++)
                        sum += shmpA[iM + k] * shmpB[kP + j];
                    shmpC[i*P + j] = sum;
                }
            shmdt(shmpA);
            shmdt(shmpB);
            shmdt(shmpC);
            exit(EXIT_SUCCESS);
        }
    }
    for (int i = 0; i < K; i++) wait(NULL);
    clock_gettime(CLOCK_MONOTONIC, &t1);
    double t_par = tiempo_seg(t0, t1);

    // Guardar resultado paralelo
    escribir_matriz("matriz_c.txt", shmpC, N, P);

    // Imprimir tiempos
    printf("Sequential time: %.6f seconds\n", t_seq);
    printf("Parallel time (%d processes): %.6f seconds\n", K, t_par);
    printf("Speedup: %.2fx\n", t_seq / t_par);

    // Limpieza
    shmdt(shmpA); shmdt(shmpB); shmdt(shmpC);
    shmctl(shA, IPC_RMID, NULL);
    shmctl(shB, IPC_RMID, NULL);
    shmctl(shC, IPC_RMID, NULL);
    free(A); free(B);
    return 0;
}

/*
Compilar:
    gcc -o matrix_mul matrix_mul.c -lrt
Ejecutar:
    ./matrix_mul <num_procesos>
*/