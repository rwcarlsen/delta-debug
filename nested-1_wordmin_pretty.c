typedef __SIZE_TYPE__ size_t; 
int printf (const char *, ...); 
void *memset (void *, int, size_t); 
int bar (int (*)(), int, void *); 
main(int argc, char **argv) {
  struct s {
    int b[argc];
  };
  int nested (struct s x) {
    sizeof(x); 
  } 
  struct s t; 
  printf("%d\n", bar (nested, argc, &t)); 
}
