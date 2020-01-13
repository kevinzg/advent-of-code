#!/usr/bin/env ruby

floor = 0

direction = {
  '(' => +1,
  ')' => -1
}

gets.strip.split('').each do |c|
  floor += direction[c]
end

puts floor
