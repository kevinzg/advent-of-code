local count = 0

local function is_unique(s)
   return #s == 2 or #s == 3 or #s == 4 or #s == 7
end

for line in io.lines() do
   local digits = {}
   for number in line:gmatch '%w+' do
      table.insert(digits, number)
   end

   for i = #digits - 3, #digits do
      if is_unique(digits[i]) then
         count = count + 1
      end
   end
end

print(count)
