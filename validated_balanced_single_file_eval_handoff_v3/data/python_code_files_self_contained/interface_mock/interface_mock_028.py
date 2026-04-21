from random import choice
import turtle
import time

def task_func(colors):
    window = turtle.Screen()
    window.bgcolor('white')
    t = turtle.Turtle()
    t.speed(1)
    for _ in range(5):
        t.color(choice(colors))
        for _ in range(4):
            t.forward(100)
            t.right(90)
        time.sleep(1)
    window.mainloop()
