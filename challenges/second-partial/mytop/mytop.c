#include <stdio.h>
#include <string.h>
#include <signal.h>
#include <stdlib.h>
#include <ctype.h>
#include <unistd.h>
#include <dirent.h>
#include <time.h>
#include <stdbool.h>

void clear();

struct proc
{
    char pid[10],
        ppid[10],
        name[50],
        state[50],
        memory[50],
        threads[50];
    int openFiles;
};
struct proc procs[1000];

struct pid
{
    int value;
};



int main(int argc, char **argv)
{
    while (1)
    {
        int nProc = 0;
        DIR *dir = opendir("/proc/");
        struct dirent *drent;

        if (dir != NULL)
        {
            while ((drent = readdir(dir)))
            {
                if (isdigit(drent->d_name[0]))
                {
                    nProc += 1;
                }
            }
            (void)closedir(dir);
        }
        else
        {
            return -1;
        }

        struct dirent *currDir;
        struct pid pids[nProc];
        int pos = 0;
        int files = 0;
        int k = 0;
        char path[30];
        char fdpath[30];

        DIR *d = opendir("/proc/");
        FILE *fp;

        strcpy(path, "/proc/");
        strcpy(fdpath, "/proc/");

        if (d != NULL)
        {
            while ((currDir = readdir(d)))
            {
                if (isdigit(currDir->d_name[0]))
                {
                    files = 0;

                    pids[pos].value = atoi(currDir->d_name);
                    pos += 1;

                    strcat(path, currDir->d_name);
                    strcat(path, "/status");

                    strcat(fdpath, currDir->d_name);
                    strcat(fdpath, "/fd");

                    fp = fopen(path, "r");
                    char buff[255];

                    while (fgets(buff, 255, (FILE *)fp))
                    {
                        if (buff[0] == 'N')
                        {
                            if (buff[1] == 'a' && buff[2] == 'm' && buff[3] == 'e')
                            {
                                int i;
                                char nombre[20];
                                int posN = 0;
                                bool done = false;
                                for (i = 5; i < 255; i++)
                                {
                                    if (buff[i] == '\n')
                                    {
                                        break;
                                    }
                                    if (buff[i] == ' ')
                                    {
                                        if (done)
                                        {
                                            nombre[posN] = '\0';
                                            break;
                                        }
                                        continue;
                                    }

                                    nombre[posN] = buff[i];
                                    posN += 1;
                                    done = true;
                                }
                                strcpy(procs[k].name, nombre);
                                memset(nombre, 0, 20);
                            }
                        }
                        else if (buff[0] == 'S')
                        {
                            if (buff[1] == 't' && buff[2] == 'a' && buff[3] == 't' && buff[4] == 'e')
                            {
                                int i;
                                char state[20];
                                int posN = 0;
                                bool done = false;

                                for (i = 6; i < 255; i++)
                                {
                                    if (buff[i] == '\n')
                                    {
                                        break;
                                    }
                                    if (buff[i] == ' ')
                                    {
                                        if (done)
                                        {
                                            state[posN] = '\0';
                                            break;
                                        }
                                        continue;
                                    }
                                    state[posN] = buff[i];
                                    posN += 1;
                                    done = true;
                                }
                                switch (state[1])
                                {
                                case 'N':
                                    strcpy(procs[k].state, "New");
                                    break;
                                case 'S':
                                    strcpy(procs[k].state, "Sleeping");
                                    break;
                                case 'I':
                                    strcpy(procs[k].state, "Idle");
                                    break;
                                case 'T':
                                    strcpy(procs[k].state, "Terminated");
                                    break;
                                case 'R':
                                    strcpy(procs[k].state, "Running");
                                    break;
                                default:
                                    strcpy(procs[k].state, state);
                                }
                                memset(state, 0, 20);
                            }
                        }
                        else if (buff[0] == 'P')
                        {
                            if (buff[1] == 'i' && buff[2] == 'd')
                            {
                                int i;
                                char dataPID[20];
                                int posN = 0;
                                bool done = false;
                                for (i = 4; i < 255; i++)
                                {
                                    if (buff[i] == '\n')
                                    {
                                        break;
                                    }
                                    if (buff[i] == ' ')
                                    {
                                        if (done)
                                        {
                                            dataPID[posN] = '\0';
                                            break;
                                        }
                                        continue;
                                    }
                                    dataPID[posN] = buff[i];
                                    posN += 1;
                                    done = true;
                                }
                                strcpy(procs[k].pid, dataPID);
                                memset(dataPID, 0, 20);
                            }
                            else if (buff[1] == 'P' && buff[2] == 'i' && buff[3] == 'd')
                            {
                                int i;
                                char dataPPID[20];
                                int posN = 0;
                                bool done = false;
                                for (i = 5; i < 255; i++)
                                {
                                    if (buff[i] == '\n')
                                    {
                                        break;
                                    }
                                    if (buff[i] == ' ')
                                    {
                                        if (done)
                                        {
                                            dataPPID[posN] = '\0';
                                            break;
                                        }
                                        continue;
                                    }
                                    dataPPID[posN] = buff[i];
                                    posN += 1;
                                    done = true;
                                }
                                strcpy(procs[k].ppid, dataPPID);
                                memset(dataPPID, 0, 20);
                            }
                        }
                        else if (buff[0] == 'V')
                        {
                            if (buff[1] == 'm' && buff[2] == 'R' && buff[3] == 'S' && buff[4] == 'S')
                            {

                                int i;
                                char memory[20];
                                int posN = 0;
                                for (i = 6; i < 255; i++)
                                {
                                    if (buff[i] == '\n')
                                    {
                                        break;
                                    }
                                    if (buff[i] == ' ')
                                    {
                                        continue;
                                    }

                                    memory[posN] = buff[i];
                                    posN += 1;
                                }
                                if (memory[0] == '\0')
                                {
                                    strcpy(procs[k].memory, " ");
                                }
                                else
                                {
                                    strcpy(procs[k].memory, memory);
                                    memset(memory, 0, 20);
                                }
                            }
                        }
                        else if (buff[0] == 'T')
                        {
                            if (buff[1] == 'h' && buff[2] == 'r' && buff[3] == 'e' && buff[4] == 'a' && buff[5] == 'd' && buff[6] == 's')
                            {
                                int i;
                                char threads[20];
                                int posN = 0;
                                bool done = false;
                                for (i = 8; i < 255; i++)
                                {
                                    if (buff[i] == '\n')
                                    {
                                        break;
                                    }
                                    if (buff[i] == ' ')
                                    {
                                        if (done == 1)
                                        {
                                            threads[posN] = '\0';
                                            break;
                                        }
                                        continue;
                                    }
                                    threads[posN] = buff[i];
                                    posN += 1;
                                    done = true;
                                }
                                if (threads[0] == '\0')
                                {
                                    strcpy(procs[k].memory, " ");
                                }
                                else
                                {
                                    strcpy(procs[k].threads, threads);
                                    memset(threads, 0, 20);
                                }
                            }
                        }
                    }
                    fclose(fp);

                    DIR *d2 = opendir(fdpath);
                    int l = 0;
                    struct dirent *ep2;

                    if (d != NULL)
                    {
                        while ((ep2 = readdir(d2)))
                        {
                            l += 1;
                        }
                    }
                    (void)closedir(d2);
                    procs[k].openFiles = l - 2;
                    k++;
                    strcpy(path, "/proc/");
                    strcpy(fdpath, "/proc/");
                }
            }
            (void)closedir(d);
        }
        else
        {
            printf("ERROR!\n");
        }
        int j = 0;
        printf("\n");
        printf("|------------|--------|-----------------------------------------------------------|-----------------|------------|------------|------------|\n");
        printf("|    PID     | Parent |                            Name                           |      State      |    Memory  |   #Threads | Open Files |\n");
        printf("|------------|--------|-----------------------------------------------------------|-----------------|------------|------------|------------|\n");

        while (procs[j].name[0] != '\0')
        {
            printf("| %-5s | %-6s | %-50s | %-15s | %-7lf M | %-6s | %-10d |\n", procs[j].pid, procs[j].ppid, procs[j].name, procs[j].state, atof(procs[j].memory) / 1000, procs[j].threads, procs[j].openFiles);
            j++;
        }
        printf("|------------|--------|-----------------------------------------------------------|-----------------|------------|------------|------------|\n");

        sleep(3);
        clear();
    }
    return 0;
}

void clear()
{
    printf("\e[1;1H\e[2J");
}
