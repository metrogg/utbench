/*You are given a string s.
  Your task is to check if the string is happy or not.
  A string is happy if its length is at least 3 and every 3 consecutive letters are distinct
  For example:
  isHappy(a) => false
  isHappy(aa) => false
  isHappy(abcd) => true
  isHappy(aabb) => false
  isHappy(adb) => true
  isHappy(xyy) => false
  */
const isHappy = (s) => {
  if (s.length < 3)
    return false;
  for (let i = 0; i < s.length - 2; i++)
    if (s[i] == s[i+1] || s[i+1] == s[i+2] || s[i] == s[i+2])
      return false;
  return true;
}
