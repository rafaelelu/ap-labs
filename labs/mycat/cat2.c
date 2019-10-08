#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>

int main(int argc, char *argv[])
{
    int fd = open(argv[1], O_RDONLY);
    if (fd == -1)
    {
        printf("Failed to open the file.\n");
        return 1;
    }
    int size = lseek(fd, sizeof(char), SEEK_END);
    close(fd);
    fd = open(argv[1], O_RDONLY);
    if (fd == -1)
    {
        printf("Failed to open the file.\n");
        return 1;
    }
    char buf[size];
    read(fd, buf, size);
    close(fd);
    buf[size - 1] = '\0';

    write(1, buf, strlen(buf));

    return 0;
}
