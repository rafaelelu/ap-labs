#include <stdio.h>

void reverse(char *a, int min, int max);
void swap(char *a, char *b);
void printArray(char *a, int length);
void newLine();

int main(){
    char input[500];
    char c;
    printf("Enter a word: \n");
    for(int i = 0, c = getchar(); c != EOF; i++){
        input[i] = c;
        c =  getchar();
        if(c == '\n'){
            char word[i+1];
            for(int j = 0; j < i+1; j++){
                word[j] = input[j];
            }
            int length = sizeof(word)/sizeof(char);
            reverse(word, 0, length-1);
            printf("Reversed word: \n");
            printArray(word, length);
            newLine();
            printf("Enter a new word: \n");
            i = -1;
            continue;
        }
    }
    return 0;
}

void reverse(char *a, int min, int max) {
    if(max-min < 1){return;}
    else {
    swap(&a[min], &a[max]);
    reverse(a, min+1, max-1);
    }
}

void swap(char *a, char *b){
    char tmp = *a;
    *a = *b;
    *b = tmp;
}

void printArray(char *a, int length){
    for(int i = 0; i < length; i++){
    printf("%c", a[i]);
    }
}

void newLine(){
    printf("\n");
}
