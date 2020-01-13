#!/usr/bin/env ruby

def say(number)
  ans = []
  i = 0

  while i < number.length
    c = number[i]
    k = 0

    while i < number.length && number[i] == c
      i = i + 1
      k = k + 1
    end

    ans.append(k)
    ans.append(c)
  end

  ans.join('')
end

number = gets.strip

for i in 1..40
  number = say number
end

puts number.length
