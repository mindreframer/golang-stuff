#include "stdio.h"
#include "stdlib.h"
#include "string.h"

int c(const void *l, const void *r) {
    return strcmp(*(char * const *)l, *(char * const *)r);
}

void print(int c, char **ss) {
    int i;
    for (i = 0; i < c; i++) {
        printf("%d: %s\n", i, ss[i]);
    }
}

int main(int argc, char **argv) {
    printf("Hello, World!\n");
    printf("ARGV before:\n");
    print(argc, argv);
    qsort(argv, argc, sizeof(char *), c);
    printf("ARGV after:\n");
    print(argc, argv);
    return 0;
}
