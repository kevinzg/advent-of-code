local depth = {}

while true do
   local n = io.read '*n'
   if n == 0 then
      break
   end
   table.insert(depth, n)
end

local prev = 0
local count = 0

for i = 2, (#depth - 2) do
   if depth[i - 1] < depth[i + 2] then
      count = count + 1
   end
end

print(count)
