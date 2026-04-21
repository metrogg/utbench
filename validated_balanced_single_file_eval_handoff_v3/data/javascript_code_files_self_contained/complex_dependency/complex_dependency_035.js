/*
  Given an array of non-negative integers, return a copy of the given array after sorting,
  you will sort the given array in ascending order if the sum( first index value, last index value) is odd,
  or sort it in descending order if the sum( first index value, last index value) is even.

  Note:
  * don't change the given array.

  Examples:
  * sortArray([]) => []
  * sortArray([5]) => [5]
  * sortArray([2, 4, 3, 0, 1, 5]) => [0, 1, 2, 3, 4, 5]
  * sortArray([2, 4, 3, 0, 1, 5, 6]) => [6, 5, 4, 3, 2, 1, 0]
  */
const sortArray = (array) => {
  let arr = array
  let tot = arr[0] + arr[arr.length-1]
  for (let j = 0; j < arr.length; j++) {
    let ind = j
    for (let k = j + 1; k < arr.length; k++) {
      if ((tot % 2 == 1 && arr[k] < arr[ind]) || (tot % 2 == 0 && arr[k] > arr[ind])) {
        ind = k
      }
    }
    let tmp = arr[j]
    arr[j] = arr[ind]
    arr[ind] = tmp
  }
  return arr
}
