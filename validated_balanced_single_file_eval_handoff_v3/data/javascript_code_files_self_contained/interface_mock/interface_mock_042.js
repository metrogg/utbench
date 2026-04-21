/*You have to write a function which validates a given date string and
  returns true if the date is valid otherwise false.
  The date is valid if all of the following rules are satisfied:
  1. The date string is not empty.
  2. The number of days is not less than 1 or higher than 31 days for months 1,3,5,7,8,10,12. And the number of days is not less than 1 or higher than 30 days for months 4,6,9,11. And, the number of days is not less than 1 or higher than 29 for the month 2.
  3. The months should not be less than 1 or higher than 12.
  4. The date should be in the format: mm-dd-yyyy

  for example: 
  validDate('03-11-2000') => true

  validDate('15-01-2012') => false

  validDate('04-0-2040') => false

  validDate('06-04-2020') => true

  validDate('06/04/2020') => false
  */
const validDate = (date) => {
  let t = date.split(/-/)
  if (t.length != 3) { return false }
  if (t[0] < 1 || t[0] > 12 || t[1] < 1) { return false }
  if (t[0] == 2 && t[1] > 29) { return false }
  if ((t[0] == 1 || t[0] == 3 || t[0] == 5 || t[0] == 7 || t[0] == 8 || t[0] == 10 || t[0] == 12) && t[1] > 31) { return false }
  if ((t[0] == 4 || t[0] == 6 || t[0] == 9 || t[0] == 11) && t[1] > 30) { return false }
  return true
}
