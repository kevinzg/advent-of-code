#!/usr/bin/env ruby

def pair_twice?(str)
  str.chars.each_index do |i|
    unless i.positive?
      next
    end

    pair = str.slice(i - 1, 2)
    ridx = str.rindex(pair)
    if ridx > i
      return true
    end
  end
  false
end

def repeated_one?(str)
  str.chars.each_index do |i|
    if i > 1 && str[i - 2] == str[i]
      return true
    end
  end
  false
end

def nice?(str)
  pair_twice?(str) &&
    repeated_one?(str)
end

count = 0

while (line = gets)
  nice?(line.strip) && count += 1
end

puts count
