package main

// You're a hungry rabbit, and you already have Eaten a certain number of carrots,
// but now you need to Eat more carrots to complete the day's meals.
// you should return an array of [ total number of Eaten carrots after your meals,
// the number of carrots left after your meals ]
// if there are not enough remaining carrots, you will Eat all remaining carrots, but will still be hungry.
// 
// Example:
// * Eat(5, 6, 10) -> [11, 4]
// * Eat(4, 8, 9) -> [12, 1]
// * Eat(1, 10, 10) -> [11, 0]
// * Eat(2, 11, 5) -> [7, 0]
// 
// Variables:
// @number : integer
// the number of carrots that you have Eaten.
// @need : integer
// the number of carrots that you need to Eat.
// @remaining : integer
// the number of remaining carrots thet exist in stock
// 
// Constrain:
// * 0 <= number <= 1000
// * 0 <= need <= 1000
// * 0 <= remaining <= 1000
// 
// Have fun :)
func Eat(number, need, remaining int) []int {
    if(need <= remaining) {
        return []int{ number + need , remaining-need }
    }
    return []int{ number + remaining , 0}
}
