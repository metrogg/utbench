/*We have an array 'arr' of N integers arr[1], arr[2], ..., arr[N].The
  numbers in the array will be randomly ordered. Your task is to determine if
  it is possible to get an array sorted in non-decreasing order by performing 
  the following operation on the given array:
      You are allowed to perform right shift operation any number of times.
  
  One right shift operation means shifting all elements of the array by one
  position in the right direction. The last element of the array will be moved to
  the starting position in the array i.e. 0th index. 

  If it is possible to obtain the sorted array by performing the above operation
  then return true else return false.
  If the given array is empty then return true.

  Note: The given list is guaranteed to have unique elements.

  For Example:
  
  moveOneBall([3, 4, 5, 1, 2])==>true
  Explanation: By performin 2 right shift operations, non-decreasing order can
               be achieved for the given array.
  moveOneBall([3, 5, 4, 1, 2])==>false
  Explanation:It is not possible to get non-decreasing order for the given
              array by performing any number of right shift operations.
              
  */
const moveOneBall = (arr) => {
  if (arr.length == 0) { return true }
  let k = 0
  let len = arr.length
  for (let i = 0; i < len; i++) {
    let t = 1;
    for (let j = 1; j < len; j++) {
      if (arr[j] < arr[j - 1]) {
        t = 0;
        break;
      }
    }
    if (t == 1) {
      k = 1;
      break;
    }
    arr.push(arr[0]);
    arr.shift()
  }
  if (k == 1) { return true }
  return false
}
