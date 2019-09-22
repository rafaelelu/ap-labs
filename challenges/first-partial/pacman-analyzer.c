#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <stdbool.h>

#define FETCHING_DATE 10
#define FINDING_ACTION 11
#define FETCHING_ACTION 12
#define FETCHING_NAME 13
#define FINDING_NEW_LINE 14

#define INSTALLED 0
#define REMOVED 1
#define UPGRADED 2

#define REPORT_FILE "packages_report.txt"

struct Package
{
    char name[50];
    char idate[17];
    char lupdate[17];
    int updates;
    char rdate[17];
    int status;
};

struct Hashtable
{
    int size;
    int nelements;
    struct Package array[1000];
};

void analizeLog(char *logFile, char *report);
bool isAction(char c1, char c2);
void addToHashtable(struct Hashtable *ht, struct Package *p);
int getHashCode(char s[]);
bool findInHashtable(struct Hashtable *ht, char key[]);
struct Package *get(struct Hashtable *ht, char key[]);
void printHashtable(struct Hashtable *ht);
void htToString(char string[], struct Hashtable *ht);
void printPackage(struct Package *p);
void pToString(char string[], struct Package *ht);
void makeReport(char *reportS, int iPackages, int rPackages, int uPackages, int cInstalled, struct Hashtable *ht);

int main(int argc, char **argv)
{

    if (argc < 2)
    {
        printf("Usage:./pacman-analizer.o pacman.txt\n");
        return 1;
    }

    analizeLog(argv[1], REPORT_FILE);

    return 0;
}

void analizeLog(char *logFile, char *report)
{
    printf("Generating Report from: [%s] log file\n", logFile);

    struct Hashtable ht = {1000, 0};
    int iPackages = 0;
    int rPackages = 0;
    int uPackages = 0;
    int cInstalled = 0;

    int fd = open(logFile, O_RDONLY);
    if (fd == -1)
    {
        printf("Failed to open the file.\n");
        return;
    }
    int size = lseek(fd, sizeof(char), SEEK_END);
    close(fd);
    fd = open(logFile, O_RDONLY);
    if (fd == -1)
    {
        printf("Failed to open the file.\n");
        return;
    }
    char buf[size];
    read(fd, buf, size);
    close(fd);
    buf[size - 1] = '\0';

    int i = 0;
    int j = 0;
    int state = FETCHING_DATE;
    char date[17];
    char name[50];
    char action[10];
    bool validLine = false;
    while (i < size)
    {
        switch (state)
        {
        case FETCHING_DATE:
            if (buf[i] != 'f')
            {
                i++;
                j = 0;
                while (buf[i] != ']')
                {
                    date[j] = buf[i];
                    j++;
                    i++;
                }
                date[j] = '\0';
                i = i + 2;
                state = FINDING_ACTION;
            }
            else
            {
                state = FINDING_ACTION;
            }
            break;

        case FINDING_ACTION:
            while (buf[i] != ' ')
            {
                i++;
            }
            i++;
            state = FETCHING_ACTION;
            break;

        case FETCHING_ACTION:
            j = 0;
            if (isAction(buf[i], buf[i + 1]))
            {
                validLine = true;
                while (buf[i] != ' ')
                {
                    action[j] = buf[i];
                    i++;
                    j++;
                }
                action[j] = '\0';
                i++;
                state = FETCHING_NAME;
            }
            else
            {
                state = FINDING_NEW_LINE;
            }
            break;

        case FETCHING_NAME:
            j = 0;
            while (buf[i] != ' ')
            {
                name[j] = buf[i];
                i++;
                j++;
            }
            name[j] = '\0';
            i++;
            state = FINDING_NEW_LINE;
            break;

        case FINDING_NEW_LINE:
            while (!(buf[i] == '\n' || buf[i] == '\0'))
            {
                i++;
            }
            i++;
            if (validLine)
            {
                if (!findInHashtable(&ht, name))
                {
                    struct Package p = {"", "", "", 0, "-", INSTALLED};
                    strcpy(p.name, name);
                    strcpy(p.idate, date);
                    addToHashtable(&ht, &p);

                    iPackages++;
                }
                else
                {
                    struct Package *p1 = get(&ht, name);
                    if (strcmp(action, "installed") == 0)
                    {
                        if (p1->status == REMOVED)
                        {
                            p1->status = INSTALLED;
                            strcpy(p1->rdate, "-");
                            rPackages--;
                        }
                    }
                    else if (strcmp(action, "removed") == 0)
                    {
                        if (p1->status == INSTALLED || p1->status == UPGRADED)
                        {
                            p1->status = REMOVED;
                            strcpy(p1->rdate, date);
                            strcpy(p1->lupdate, date);
                            p1->updates = p1->updates + 1;
                            rPackages++;
                        }
                    }
                    else if (strcmp(action, "upgraded") == 0)
                    {
                        if (p1->status == INSTALLED)
                        {
                            p1->status = UPGRADED;
                            strcpy(p1->lupdate, date);
                            p1->updates = p1->updates + 1;
                            uPackages++;
                        }
                        else if (p1->status == UPGRADED)
                        {
                            strcpy(p1->lupdate, date);
                            p1->updates = p1->updates + 1;
                        }
                    }
                }
            }
            validLine = false;
            state = FETCHING_DATE;
            if (i >= size - 1)
            {
                i = i + 1;
            }
            break;
        }
    }
    cInstalled = iPackages - rPackages;
    char reportS[100000];
    makeReport(reportS, iPackages, rPackages, uPackages, cInstalled, &ht);
    fd = open(report, O_CREAT | O_WRONLY, 0600);
    if (fd == -1)
    {
        printf("Failed to open the file.\n");
        return;
    }
    write(fd, reportS, strlen(reportS));
    close(fd);

    printf("Report is generated at: [%s]\n", report);
}

bool isAction(char c1, char c2)
{
    if (c1 == 'i' || c1 == 'u')
    {
        return true;
    }
    else if (c1 == 'r' && c2 == 'e')
    {
        return true;
    }
    else
    {
        return false;
    }
}

void addToHashtable(struct Hashtable *ht, struct Package *p)
{
    for (int i = 0; i < ht->nelements + 1; i++)
    {
        int hashValue = getHashCode(p->name) + i;
        int index = hashValue % ht->size;
        if (strcmp(ht->array[index].name, "") == 0)
        {
            ht->array[index] = *p;
            break;
        }
    }
    ht->nelements += 1;
}

int getHashCode(char s[])
{
    int n = strlen(s);
    int hashValue = 0;

    for (int i = 0; i < n; i++)
    {
        hashValue = hashValue * 31 + s[i];
    }

    hashValue = hashValue & 0x7fffffff;
    return hashValue;
}

bool findInHashtable(struct Hashtable *ht, char key[])
{
    for (int i = 0; i < ht->nelements + 1; i++)
    {
        int hashValue = getHashCode(key) + i;
        int index = hashValue % ht->size;
        if (strcmp(ht->array[index].name, key) == 0)
        {
            return true;
        }
        else if (strcmp(ht->array[index].name, "") == 0)
        {
            return false;
        }
    }
    return false;
}

struct Package *get(struct Hashtable *ht, char key[])
{
    for (int i = 0; i < ht->nelements + 1; i++)
    {
        int hashValue = getHashCode(key) + i;
        int index = hashValue % ht->size;
        if (strcmp(ht->array[index].name, key) == 0)
        {
            return &ht->array[index];
        }
        else if (strcmp(ht->array[index].name, "") == 0)
        {
            return NULL;
        }
    }
    return NULL;
}

void printHashtable(struct Hashtable *ht)
{
    printf("ht.size: %d\n", ht->size);
    printf("ht.nelements: %d\n", ht->nelements);
    printf("ht.array: \n");
    for (int i = 0; i < ht->size; i++)
    {
        if (strcmp(ht->array[i].name, "") != 0)
        {
            printPackage(&ht->array[i]);
            printf("\n");
        }
    }
}

void htToString(char string[], struct Hashtable *ht)
{
    for (int i = 0; i < ht->size; i++)
    {
        if (strcmp(ht->array[i].name, "") != 0)
        {
            pToString(string, &ht->array[i]);
            strcat(string, "\n\n");
        }
    }
}

void printPackage(struct Package *p)
{
    printf("- Package Name        : %s\n", p->name);
    printf("  - Install date      : %s\n", p->idate);
    printf("  - Last update date  : %s\n", p->lupdate);
    printf("  - How many updates  : %d\n", p->updates);
    printf("  - Removal date      : %s\n", p->rdate);
}

void pToString(char string[], struct Package *p)
{
    strcat(string, "- Package Name        : ");
    strcat(string, p->name);
    strcat(string, "\n");
    strcat(string, "  - Install date      : ");
    strcat(string, p->idate);
    strcat(string, "\n");
    strcat(string, "  - Last update date  : ");
    strcat(string, p->lupdate);
    strcat(string, "\n");
    strcat(string, "  - How many updates  : ");
    char numBuf[20];
    sprintf(numBuf, "%d\n", p->updates);
    strcat(string, numBuf);
    strcat(string, "  - Removal date      : ");
    strcat(string, p->rdate);
}

void makeReport(char *reportS, int iPackages, int rPackages, int uPackages, int cInstalled, struct Hashtable *ht)
{
    strcat(reportS, "Pacman Packages Report\n");
    strcat(reportS, "----------------------\n");
    char numBuf[120];
    sprintf(numBuf, "- Installed packages : %d\n- Removed packages   : %d\n- Upgraded packages  : %d\n- Current installed  : %d\n\n", iPackages, rPackages, uPackages, cInstalled);
    strcat(reportS, numBuf);
    htToString(reportS, ht);
}
