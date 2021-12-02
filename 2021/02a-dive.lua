local x = 0
local depth = 0

for line in io.lines() do
   local action = line:match '%a+'
   local units = tonumber(line:match '%d+')

   if action == 'forward' then
      x = x + units
   elseif action == 'down' then
      depth = depth + units
   else
      depth = depth - units
   end
end

print(x * depth)
