/*
  You are given a list of integers.
  Write a function nextSmallest() that returns the 2nd smallest element of the list.
  Return null if there is no such element.
  
  nextSmallest([1, 2, 3, 4, 5]) == 2
  nextSmallest([5, 1, 4, 3, 2]) == 2
  nextSmallest([]) == null
  nextSmallest([1, 1]) == null
  */
const nextSmallest = (lst) => {
  let arr = lst
  for (let j = 0; j < arr.length; j++) {
    let ind = j
    for (let k = j + 1; k < arr.length; k++) {
      if (arr[k] < arr[ind]) {
        ind = k
      }
    }
    let tmp = arr[j]
    arr[j] = arr[ind]
    arr[ind] = tmp
  }
  let smallest = arr[0]
  let pt = 1
  while(pt<arr.length){
    if(arr[pt]>smallest){
      return arr[pt]
    }
    pt++
  }
  return null
}
