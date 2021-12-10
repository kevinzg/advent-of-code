local input = io.read '*a'

local crabs = {}
for i = 0, 2000 do
   crabs[i] = 0
end

local right = 0
local cost = 0

for d in input:gmatch '%d+' do
   d = tonumber(d)
   crabs[d] = crabs[d] + 1
   cost = cost + d

   if d > 0 then
      right = right + 1
   end
end

local best = cost
local left = 0

for i = 1, 2000 do
   left = left + crabs[i - 1]
   right = right - crabs[i]

   cost = cost + left - right - crabs[i]
   if cost < best then
      best = cost
   end
end

print(best)
