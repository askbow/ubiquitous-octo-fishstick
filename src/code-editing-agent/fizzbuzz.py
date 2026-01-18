#!/usr/bin/env python3

for i in range(1, 16):
    print("Fizz"*(i%3==0) + "Buzz"*(i%5==0) or i)
