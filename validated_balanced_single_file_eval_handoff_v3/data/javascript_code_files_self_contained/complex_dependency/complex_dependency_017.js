/*You will be given the name of a class (a string) and a list of extensions.
  The extensions are to be used to load additional classes to the class. The
  strength of the extension is as follows: Let CAP be the number of the uppercase
  letters in the extension's name, and let SM be the number of lowercase letters
  in the extension's name, the strength is given by the fraction CAP - SM.
  You should find the strongest extension and return a string in this
  format: ClassName.StrongestExtensionName.
  If there are two or more extensions with the same strength, you should
  choose the one that comes first in the list.
  For example, if you are given "Slices" as the class and a list of the
  extensions: ['SErviNGSliCes', 'Cheese', 'StuFfed'] then you should
  return 'Slices.SErviNGSliCes' since 'SErviNGSliCes' is the strongest extension
  (its strength is -1).
  Example:
  for strongestExtension('my_class', ['AA', 'Be', 'CC']) == 'my_class.AA'
  */
const strongestExtension = (class_name, extensions) => {
  let u = 0
  let s = -Infinity
  for (let i = extensions.length - 1; i >= 0; i--) {
    let y = 0
    for (let j = 0; j < extensions[i].length; j++) {
      let k = extensions[i][j].charCodeAt()
      if (k >= 65 && k <= 90) { y += 1 }
      if (k >= 97 && k <= 122) { y -= 1 }
    }
    if (y >= s) {
      s = y;
      u = i;
    }
  }
  return class_name + '.' + extensions[u]
}
