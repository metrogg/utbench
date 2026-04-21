/* For a given list of integers, return a tuple consisting of a sum and a product of all the integers in a list.
  Empty sum should be equal to 0 and empty product should be equal to 1.
  >>> sumProduct([])
  (0, 1)
  >>> sumProduct([1, 2, 3, 4])
  (10, 24)
  */
const sumProduct = (numbers, int) => {
  var sum_value = 0, prod_value = 1;
  for (const n of numbers) {
    sum_value += n;
    prod_value *= n;
  }
  return [sum_value, prod_value];
}
