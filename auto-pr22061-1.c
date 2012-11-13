int N = 1;
foo() {} /* */
bar (char a[2][N]) { a[1][0] = N; }
main (void)
{
  void *x;

  (x, 2 * N);
  if (N[(char *) x] != N)
    (0);
}

