#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv)
{
    int LOWER = 0;
    int UPPER = 300;
    int STEP = 20;

    switch(argc){
        case 2:
            LOWER = atoi(argv[1]);
            UPPER = LOWER;
            break;
        case 4:
            LOWER = atoi(argv[1]);
            UPPER = atoi(argv[2]);
            STEP = atoi(argv[3]);
            break;
        default:
            printf("You need to send either the fahrenheit degrees to convert or the start, end and increment.\n");
            printf("How to execute: ./fahrenheit_celsius <num_farenheit_degrees>\n");
            printf("OR\n");
            printf("./fahrenheit_celsius <start> <end> <increment>\n");
            return 1;
    }

    int fahr;
    for (fahr = LOWER; fahr <= UPPER; fahr = fahr + STEP){
	    printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
    }

    return 0;
}
