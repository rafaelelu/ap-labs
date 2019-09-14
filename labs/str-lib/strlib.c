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
    int lengthOfNewString;
    int lengthOfOrigin = mystrlen(origin);
    int lengthOfAddition = mystrlen(addition);
    lengthOfNewString =  lengthOfOrigin + lengthOfAddition + 1;
    char newString[lengthOfNewString];
    int i = 0;
    for(int j = 0; j < lengthOfOrigin; j++){
        newString[i] = origin[j];
        i++;
    }
    for(int k = 0; k < lengthOfAddition; k++){
        newString[i] = addition[k];
        i++;
    }
    newString[i] = '\0';
    return newString;
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
