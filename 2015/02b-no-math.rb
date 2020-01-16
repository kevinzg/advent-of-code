#!/usr/bin/env ruby

def ribbon(l, w, h)
  perimeter = [l + w, l + h, w + h]
  l * w * h + 2 * perimeter.min
end

total = 0

while (line = gets)
  dim = line.strip.split('x').map(&:to_i)
  total += ribbon(dim[0], dim[1], dim[2])
end

puts total
