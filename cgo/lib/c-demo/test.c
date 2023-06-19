#include <stdio.h>

// 可以通过手工方式声明了 goPrintln 和 number_add_mod 两个导出函数。
// 这样我们就实现了从多个 Go 包导出 C 函数了。
void goPrintln(char*);
int number_add_mod(int a, int b, int mod);

int main() {
    int a = 10;
    int b = 5;
    int c = 12;

    int x = number_add_mod(a, b, c);
    printf("(%d+%d)%%%d = %d\n", a, b, c, x);

    goPrintln("done");
    return 0;
}
