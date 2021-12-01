local depth = {}

while true do
   local n = io.read '*n'
   if n == nil then
      break
   end
   table.insert(depth, n)
end

local prev = 0
local count = -1 -- first one doesn't count

for _, d in ipairs(depth) do
   if d - prev > 0 then
      count = count + 1
   end
   prev = d
end

print(count)
