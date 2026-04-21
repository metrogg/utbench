/*
 Given the lengths of the three sides of a triangle. Return true if the three
  sides form a right-angled triangle, false otherwise.
  A right-angled triangle is a triangle in which one angle is right angle or
  90 degree.
  Example:
  rightAngleTriangle(3, 4, 5) == true
  rightAngleTriangle(1, 2, 3) == false
  */
const rightAngleTriangle = (a, b, c) => {
  return (a * a + b * b == c * c || a * a == b * b + c * c || b * b == a * a + c * c)
}
