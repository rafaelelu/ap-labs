#include <stdio.h>
#include "logger.h"
#include <pthread.h>
#include <stdlib.h>
#include <errno.h>

struct dPArgs {
    long *vec1;
    long *vec2;
    int i;
    int j;
    int k;
    int lock;
};

#define NUM_THREADS 2000
#define RESULT_SIZE 4000001

int NUM_BUFFERS;
long **buffers;
pthread_mutex_t *mutexes;
long result[RESULT_SIZE];

pthread_attr_t attr;
long *matA;
long *matB;

// Reads matrix file and returns a long type array with content of matrix.
long *readMatrix(char *filename) {
    int i = 1;
    size_t len = 0;
    ssize_t nread;
    char *line = NULL;
    long *arr = calloc(RESULT_SIZE, sizeof(long));

    FILE *fp = fopen(filename, "rb");
    if (!fp) {
        panicf("File couldn't be opened\n");
    }

    while ((nread = getline(&line, &len, fp)) != -1) {
        arr[i++] = strtol(line, NULL, 10);
    }
    free(line);
    return arr;
}

// Returns the specified array that represents a column in the matrix array.
long *getColumn(int col, long *matrix) {
    long *colArr = calloc(NUM_THREADS + 1, sizeof(long));

    for (int i = 1; i < NUM_THREADS + 1; i++) {
        colArr[i] = matrix[((i * NUM_THREADS) - NUM_THREADS) + col];
    }
    return colArr;
}

// Returns the specified array that represents a row in the matrix array.
long *getRow(int row, long *matrix) {
    long *rowArr = calloc(NUM_THREADS + 1, sizeof(long));

    for (int i = 1; i < NUM_THREADS + 1; i++) {
        rowArr[i] = matrix[((row * NUM_THREADS) - NUM_THREADS) + i];
    }
    return rowArr;
}

// Search for an available buffer, if so it returns the available lock id which is the same buffer id, otherwise returns -1
int getLock() {
    int r = -1;
    int i = 0;
    while (r != 0) {
        r = pthread_mutex_trylock(&mutexes[i % NUM_BUFFERS]);
        if (r == 0) {
            return i % NUM_BUFFERS;
        }
        i++;
    }
    return -1;
}

// Releases a buffer and unlock the mutex. Returns 0 for successful unlock, otherwise -1.
int releaseLock(int lock) {
    if (pthread_mutex_unlock(&mutexes[lock]) == 0) {
        return 0;
    }
    return -1;
}

// Given 2 arrays of 2000 lenght as one struct argument, it calculates the dot product.
void *dotProduct(void *args)
{
    struct dPArgs *dArgs = (struct dPArgs *)args;
    long temp = 0;
    int lock = getLock();
    dArgs->vec1 = getRow(dArgs->j, matA);
    dArgs->vec2 = getColumn(dArgs->k, matB);
    for (int i = 1; i < NUM_THREADS + 1; i++) {
        temp += dArgs->vec1[i] * dArgs->vec2[i];
    }

    free(dArgs->vec1);
    free(dArgs->vec2);
    buffers[dArgs->lock][0] = temp;
    result[dArgs->i] = buffers[dArgs->lock][0];
    releaseLock(lock);
    return NULL;
}

long *multiply(long *matA, long *matB) {
    pthread_attr_init(&attr);
    pthread_t threads[NUM_THREADS];
    struct dPArgs *dArgsArr = calloc(RESULT_SIZE, sizeof(struct dPArgs));

    int i = 1;
    for (int j = 1; j < NUM_THREADS + 1; j++) {
        for (int k = 1; k < NUM_THREADS + 1; k++) {
            dArgsArr[i].i = i;
            dArgsArr[i].j = j;
            dArgsArr[i].k = k;
            int error = pthread_create(&threads[i % NUM_THREADS], &attr, dotProduct, (void *)&dArgsArr[i]);
            if (error != 0) {
                panicf("Error when creating a thread: %d", error);
                exit(-1);
            }
            pthread_detach(threads[i % NUM_THREADS]);
            i++;
        }
    }
}

// Saves result matrix into a new result.dat file, return 0 for a successful operation, otherwise it will return -1
int saveResultMatrix(long *result) {
    FILE *fp = fopen("result.dat", "w");
    if (!fp) {
        panicf("File couldn't be opened\n");
        return -1;
    }

    for (int i = 1; i < RESULT_SIZE; i++) {
        char *temp = calloc(8, sizeof(char));
        sprintf(temp, "%ld", result[i]);
        fputs(temp, fp);
        fputc('\n', fp);
        free(temp);
    }
    return 0;
}

int main(int argc, char *argv[]) {
    if (argc <= 2) {
        errorf("How to execute: ./multiplier -n [number of buffers]\n");
        return -1;
    } else {
    infof("Calculating the result of the multiplication...\n");
    }
    NUM_BUFFERS = atoi(argv[2]);
    buffers = calloc(NUM_BUFFERS, sizeof(long *));

    for (int i = 0; i < NUM_BUFFERS; i++) {
        buffers[i] = calloc(1, sizeof(long));
    }
    mutexes = calloc(NUM_BUFFERS, sizeof(pthread_mutex_t));

    for (int i = 0; i < NUM_BUFFERS; i++) {
        pthread_mutex_init(&mutexes[i], NULL);
    }
    matA = readMatrix("matA.dat");
    matB = readMatrix("matB.dat");
    multiply(matA, matB);
    infof("The matrix multiplication was finished successfully\n");
    saveResultMatrix(result);
    free(buffers);
    free(matA);
    free(matB);

    for (int i = 0; i < NUM_BUFFERS; i++) {
        pthread_mutex_destroy(&mutexes[i]);
    }
    pthread_exit(NULL);
}
