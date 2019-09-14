int mystrlen(char *str);
char *mystradd(char *origin, char *addition);
int mystrfind(char *origin, char *substr);


int mystrlen(char *str){
    int i = 0;
    while(str[i] != '\0'){
        i++;
    }
    return i;
}

char *mystradd(char *origin, char *addition){

    int originLen = mystrlen(origin);
    int i;
    for(i = 0; addition[i] != '\0'; i++){
        origin[originLen + i] = addition[i];
    }
    origin[originLen + i] = '\0';
    return origin;
}

int mystrfind(char *origin, char *substr){
    int originLen = mystrlen(origin);
    int substrLen = mystrlen(substr);
    int j = 0;
    for(int i = 0; i <= originLen; i++){
        if(origin[i] == substr[0] && j < 1){
            j++;
        } else if(origin[i] == substr[j]){
            j++;
        } else {
            j = 0;
        }
        if(j >= substrLen){
            return 1;
        }
    }
    return 0;
}
