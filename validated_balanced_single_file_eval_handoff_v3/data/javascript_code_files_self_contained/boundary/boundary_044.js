/*
  Evaluates polynomial with coefficients xs at point x.
  return xs[0] + xs[1] * x + xs[1] * x^2 + .... xs[n] * x^n
  */
const poly = (xs, x) => {
  return xs.reduce((prev, item, index) => {
    return prev + item * Math.pow(x, index);
  }, 0);
}

/*
  xs are coefficients of a polynomial.
  findZero find x such that poly(x) = 0.
  findZero returns only only zero point, even if there are many.
  Moreover, findZero only takes list xs having even number of coefficients
  and largest non zero coefficient as it guarantees
  a solution.
  >>> round(findZero([1, 2]), 2) # f(x) = 1 + 2x
  -0.5
  >>> round(findZero([-6, 11, -6, 1]), 2) # (x - 1) * (x - 2) * (x - 3) = -6 + 11x - 6x^2 + x^3
  1.0
  */
const findZero = (xs) => {
  var begin = -1.0, end = 1.0;
  while (poly(xs, begin) * poly(xs, end) > 0) {
    begin *= 2.0;
    end *= 2.0;
  }
  while (end - begin > 1e-10) {
    let center = (begin + end) / 2.0;
    if (poly(xs, center) * poly(xs, begin) > 0)
      begin = center;
    else
      end = center;
  }
  return begin;
}
