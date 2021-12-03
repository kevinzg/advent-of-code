local max_bits = 12

local function get_bit(number, k)
   if k > #number then
      return '0'
   end
   k = #number - k + 1
   return number:sub(k, k)
end

local function to_int(number)
   local x = 0
   for i = max_bits, 1, -1 do
      if get_bit(number, i) == '1' then
         x = x + 2 ^ (i - 1)
      end
   end
   return x
end

local function search(numbers, type, k)
   local ones = 0
   local zeros = 0

   for _, n in ipairs(numbers) do
      local bit = get_bit(n, k)
      if bit == '1' then
         ones = ones + 1
      else
         zeros = zeros + 1
      end
   end

   if ones == 0 then
      return false, numbers
   end

   local target = '1'
   if zeros > ones then
      target = '0'
   end

   if type == 'co2' then
      target = target == '1' and '0' or '1'
   end

   local candidates = {}

   for _, n in ipairs(numbers) do
      local bit = get_bit(n, k)
      if bit == target then
         table.insert(candidates, n)
      end
   end

   if #candidates == 0 then
      return true, numbers[#numbers]
   elseif #candidates == 1 then
      return true, candidates[1]
   end

   return false, candidates
end

local function get_rating(numbers, type)
   for i = max_bits, 1, -1 do
      local last_one, result = search(numbers, type, i)
      if last_one then
         return result
      end
      numbers = result
   end
   return nil
end

local numbers = {}

for line in io.lines() do
   table.insert(numbers, line)
end

local oxygen = get_rating(numbers, 'oxygen')
local co2 = get_rating(numbers, 'co2')

print(to_int(oxygen) * to_int(co2))
