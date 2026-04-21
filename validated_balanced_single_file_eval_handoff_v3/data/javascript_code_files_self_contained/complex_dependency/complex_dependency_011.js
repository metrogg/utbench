/* From a supplied list of numbers (of length at least two) select and return two that are the closest to each
  other and return them in order (smaller number, larger number).
  >>> findClosestElements([1.0, 2.0, 3.0, 4.0, 5.0, 2.2])
  (2.0, 2.2)
  >>> findClosestElements([1.0, 2.0, 3.0, 4.0, 5.0, 2.0])
  (2.0, 2.0)
  */
const findClosestElements = (numbers) => {
  var closest_pair, distance;
  for (let i = 0; i < numbers.length; i++)
    for (let j = 0; j < numbers.length; j++)
      if (i != j) {
        let a = numbers[i], b = numbers[j];
        if (distance == null) {
          distance = Math.abs(a - b);
          closest_pair = [Math.min(a, b), Math.max(a, b)];
        } else {
          let new_distance = Math.abs(a - b);
          if (new_distance < distance) {
            distance = new_distance;
            closest_pair = [Math.min(a, b), Math.max(a, b)];
          }
        }
      }
  return closest_pair;
}
