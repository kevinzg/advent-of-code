#!/usr/bin/env ruby

def wrapping_paper(l, w, h)
  faces = [l * w, l * h, w * h]
  faces.sum * 2 + faces.min
end

total = 0

while (line = gets)
  dim = line.strip.split('x').map(&:to_i)
  total += wrapping_paper(dim[0], dim[1], dim[2])
end

puts total
