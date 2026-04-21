/*
  Given an array of integers nums, find the minimum sum of any non-empty sub-array
  of nums.
  Example
  minSubArraySum([2, 3, 4, 1, 2, 4]) == 1
  minSubArraySum([-1, -2, -3]) == -6
  */
const minSubArraySum = (nums) => {
  let min = nums[0]
  for (let i = 0; i < nums.length; i++) {
    for (let j = i + 1; j <= nums.length; j++) {
      let s = 0;
      for (let k = i; k < j; k++) {
        s += nums[k]
      }
      if (s < min) { min = s }
    }
  }
  return min
}
