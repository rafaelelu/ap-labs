#include <stdio.h>


int main(int argc, char **argv)
{
  if (argc < 2) {
    printf("You need to send the number of the grades to convert\n");
    printf("How to execute: ./fc <number>");

    return 1;
  }
    fahr = atoi(argv[1]);
  int fahr = atoi(argv[1]);
  printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
  return 0;
}
