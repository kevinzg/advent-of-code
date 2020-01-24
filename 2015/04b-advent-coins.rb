#!/usr/bin/env ruby

require 'digest'

key = gets.strip

i = 0

while true
  md5 = Digest::MD5.hexdigest(key + i.to_s)
  if md5.start_with? '000000'
    puts i
    break
  end
  i += 1
end
