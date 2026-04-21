/*Given a string representing a space separated lowercase letters, return a dictionary
  of the letter with the most repetition and containing the corresponding count.
  If several letters have the same occurrence, return all of them.
  
  Example:
  histogram('a b c') == {'a': 1, 'b': 1, 'c': 1}
  histogram('a b b a') == {'a': 2, 'b': 2}
  histogram('a b c a b') == {'a': 2, 'b': 2}
  histogram('b b b b a') == {'b': 4}
  histogram('') == {}

  */
const histogram = (test) => {
  let d = {}
  let t = test.split(/\s/)
  if (test == '') { t = [] }
  for (m in t) {
    if (t[m] in d) {
      d[t[m]]++
    }
    else {
      d[t[m]] = 1
    }
  }
  s = Object.keys(d).sort(function (a, b) { return - d[a] + d[b]; });
  if (s.length == 0) { return {} }
  let g = d[s[0]]
  let l = {}
  for (let ss=0; ss<s.length; ss++) {
    if (d[s[ss]] == g) {
      l[s[ss]] = d[s[ss]]
    }
  }
  return l
}
