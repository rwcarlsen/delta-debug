main(int argc, char **argv)
{
  struct s {
    int a; 
    char b[argc];
  };

  int nested (struct s x) {
    sizeof(x);
  }
  struct s t;
  (&t, sizeof(t));
  t.a = printf("%d\n", bar (nested, &t));
  }


