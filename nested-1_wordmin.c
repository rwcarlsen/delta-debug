printf (const *, ...);
*memset (void size_t);
bar (int (*)(), int, void *);

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

