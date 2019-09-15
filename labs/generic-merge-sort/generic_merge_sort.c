/*

In this pr_ogram, there are two arrays, one for numbers and one for strings. To sort the numeric one, run the code with
"./generic_merge_sort -n". To sort the string array run the code with "./generic_merge_sort".

*/

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

void mergeSort(void *ptr[], int left, int right, int (*comp)(void *, void *));
int numcmp(char *s1, char *s2);

int main(int argc, char *argv[]){
    char *numArr[] = {"9","1","8","2","7","3","6","4","5"};
    char *strArr[] = {"Diego", "Alberto", "Juan", "David", "Mario", "Brian", "Pablo", "Rafael", "Santiago"};

    int elements = 9;

    if(argc > 1 && strcmp(argv[1], "-n") == 0){
        mergeSort(numArr, 0, elements-1, numcmp);

        for(int i = 0; i < elements; i++){
            printf("%s, ", numArr[i]);
        }
        printf("\n");


    } else {
        mergeSort(strArr, 0, elements-1, strcmp);
        for(int i = 0; i < elements; i++){
            printf("%s, ", strArr[i]);
        }
        printf("\n");
    }
    return 0;
}

void mergeSort(void *ptr[], int left, int right, int (*comp)(void *, void *)) {
    if (left < right) { 
        int middle = left+(right-left)/2; 
        mergeSort(ptr, left, middle, comp); 
        mergeSort(ptr, middle+1, right, comp); 
        
        int i, j, k; 
        int n1 = middle - left + 1; 
        int n2 = right - middle; 

        void *leftArr[n1], *rightArr[n2]; 
        for (i = 0; i < n1; i++) 
            leftArr[i] = ptr[left + i]; 
        for (j = 0; j < n2; j++) 
            rightArr[j] = ptr[middle + 1+ j]; 
        i = 0;
        j = 0;
        k = left;
        while (i < n1 && j < n2) { 
            if(comp(leftArr[i],rightArr[j]) < 1 || comp(leftArr[i], rightArr[j]) == 0) { 
                ptr[k] = leftArr[i]; 
                i++; 
            } else { 
                ptr[k] = rightArr[j]; 
                j++; 
            } 
            k++; 
        } 
        while (i < n1) { 
            ptr[k] = leftArr[i]; 
            i++; 
            k++; 
        } 
        while (j < n2) { 
            ptr[k] = rightArr[j]; 
            j++; 
            k++; 
        } 
    }
}

int numcmp(char *s1, char *s2) {
    double v1, v2;

    v1 = atof(s1);
    v2 = atof(s2);
    if (v1 < v2){
        return -1;
    } else if (v1 > v2) {
        return 1;
    } else {
        return 0;
    }
}
