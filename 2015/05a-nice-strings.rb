#!/usr/bin/env ruby

def three_vocals?(str)
  count = 0
  'aeiou'.each_char do |v|
    count += str.count v
  end
  count >= 3
end

def duplicated_letters?(str)
  str.chars.each_index do |i|
    if i.positive? && str[i - 1] == str[i]
      return true
    end
  end
  false
end

def doesnt_contain?(str, banned)
  banned.each do |s|
    unless (str.index s).nil?
      return false
    end
  end
  true
end

def nice?(str)
  three_vocals?(str) &&
    duplicated_letters?(str) &&
    doesnt_contain?(str, %w[ab cd pq xy])
end

count = 0

while (line = gets)
  nice?(line.strip) && count += 1
end

puts count
