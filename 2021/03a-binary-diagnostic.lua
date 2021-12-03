local counter = {}

local max_bits = 12

for _ = 1, max_bits do
   table.insert(counter, 0)
end

local total_lines = 0

for line in io.lines() do
   total_lines = total_lines + 1

   for i = 1, #line do
      local c = line:sub(i, i)
      if c == '1' then
         local idx = #line - i + 1
         counter[idx] = counter[idx] + 1
      end
   end
end

local gamma = 0
local epsilon = 0

for i = max_bits, 1, -1 do
   if counter[i] > total_lines - counter[i] then
      gamma = gamma + 2 ^ (i - 1)
   elseif counter[i] > 0 then
      epsilon = epsilon + 2 ^ (i - 1)
   end
end

print(gamma * epsilon)
