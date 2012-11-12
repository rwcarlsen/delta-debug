
main(int argc)
{
  struct s {
    char b[argc];
  };

  int nested (struct s x) {
    sizeof(x);
  }

  struct s t;
  printf("%d\n", bar(nested, &t));
}


