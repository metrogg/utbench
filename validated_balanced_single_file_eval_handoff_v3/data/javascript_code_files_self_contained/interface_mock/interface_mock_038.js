/* There are eight planets in our solar system: the closerst to the Sun
  is Mercury, the next one is Venus, then Earth, Mars, Jupiter, Saturn,
  Uranus, Neptune.
  Write a function that takes two planet names as strings planet1 and planet2.
  The function should return a tuple containing all planets whose orbits are
  located between the orbit of planet1 and the orbit of planet2, sorted by
  the proximity to the sun.
  The function should return an empty tuple if planet1 or planet2
  are not correct planet names.
  Examples
  bf("Jupiter", "Neptune") ==> ("Saturn", "Uranus")
  bf("Earth", "Mercury") ==> ("Venus")
  bf("Mercury", "Uranus") ==> ("Venus", "Earth", "Mars", "Jupiter", "Saturn")
  */
const bf = (planet1, planet2) => {
  let y = ['Mercury', 'Venus', 'Earth', 'Mars', 'Jupiter', 'Saturn', 'Uranus', 'Neptune']
  let u = []
  let lo = -1
  let hi = -1
  for (let i = 0; i < 8; i++) {
    if (y[i] == planet1) { lo = i }
  }
  for (let i = 0; i < 8; i++) {
    if (y[i] == planet2) { hi = i }
  }
  if (lo == -1 || hi == -1 || lo == hi) { return [] }
  if (lo > hi) {
    let tmp = lo;
    lo = hi;
    hi = tmp;
  }
  for (let i = lo + 1; i < hi; i++) {
    u.push(y[i])
  }
  return u
}
