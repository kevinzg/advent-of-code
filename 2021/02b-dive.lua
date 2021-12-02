local x = 0
local depth = 0
local aim = 0

for line in io.lines() do
   local action = line:match '%a+'
   local units = tonumber(line:match '%d+')

   if action == 'forward' then
      x = x + units
      depth = depth + aim * units
   elseif action == 'down' then
      aim = aim + units
   else
      aim = aim - units
   end
end

print(x * depth)
