local input = io.read '*a'

local crabs = {}
for i = 0, 2000 do
   crabs[i] = 0
end

local right = 0
local cost = 0
local distance = 0

for d in input:gmatch '%d+' do
   d = tonumber(d)
   crabs[d] = crabs[d] + 1
   cost = cost + (((d + 1) * d) / 2)
   distance = distance + d

   if d > 0 then
      right = right + 1
   end
end

local best = cost
local left = 0
local cost_left = 0
local cost_right = cost
local dl = 0
local dr = distance

for i = 1, 2000 do
   left = left + crabs[i - 1]
   dl = dl + left
   cost_left = cost_left + dl

   cost_right = cost_right - dr
   dr = dr - right
   right = right - crabs[i]

   cost = cost_left + cost_right

   if cost < best then
      best = cost
   end
end

print(best)
