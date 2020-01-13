#!/usr/bin/env ruby

floor = 0

direction = {
  '(' => +1,
  ')' => -1
}

gets.strip.split('').each_with_index do |c, i|
  floor += direction[c]
  if floor == -1
    puts i + 1
    break
  end
end
