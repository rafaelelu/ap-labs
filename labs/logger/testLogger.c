# include "logger.c"

int main() {

    infof("info: %d\n", 7);
    warnf("warning %c\n", 'a');
    errorf("error %f\n", 7.0f);
    panicf("panic\n");

    return 0;
}
