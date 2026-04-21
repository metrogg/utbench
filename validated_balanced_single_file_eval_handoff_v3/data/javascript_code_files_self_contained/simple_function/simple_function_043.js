/*
  Write a function that takes an integer a and returns true
  if this ingeger is a cube of some integer number.
  Note: you may assume the input is always valid.
  Examples:
  iscube(1) ==> true
  iscube(2) ==> false
  iscube(-1) ==> true
  iscube(64) ==> true
  iscube(0) ==> true
  iscube(180) ==> false
  */
const iscube = (a) => {
  a = Math.abs(a);
  return (Math.pow(Math.round(Math.pow(a, 1.0 / 3.0)), 3) == a);
}
