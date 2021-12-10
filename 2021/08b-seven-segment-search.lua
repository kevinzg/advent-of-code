local get_digit = (function()
   local numbers = {
      [0] = { 1, 2, 3, 5, 6, 7 },
      [1] = { 3, 6 },
      [2] = { 1, 3, 4, 5, 7 },
      [3] = { 1, 3, 4, 6, 7 },
      [4] = { 2, 3, 4, 6 },
      [5] = { 1, 2, 4, 6, 7 },
      [6] = { 1, 2, 4, 5, 6, 7 },
      [7] = { 1, 3, 6 },
      [8] = { 1, 2, 3, 4, 5, 6, 7 },
      [9] = { 1, 2, 3, 4, 6, 7 },
   }

   local bitmasks = {}
   for n, t in pairs(numbers) do
      local m = 0
      for i = 1, #t do
         m = m | (1 << t[i])
      end
      bitmasks[m] = n
   end

   return function(m)
      return bitmasks[m]
   end
end)()

local function permutate(arr, mem, m, cb)
   if #mem == #arr then
      cb(mem)
      return
   end

   for i = 1, #arr do
      if m & (1 << i) == 0 then
         table.insert(mem, arr[i])
         permutate(arr, mem, m | (1 << i), cb)
         table.remove(mem)
      end
   end
end

local function string_to_number(s, key)
   local x = 0
   for c in s:gmatch '%w' do
      x = x | (1 << key[c])
   end
   return get_digit(x)
end

local keys = { 'a', 'b', 'c', 'd', 'e', 'f', 'g' }

local function solve(digits)
   local key

   local function check_permutation(p)
      local rev = {}
      for k, v in ipairs(p) do
         rev[v] = k
      end

      for j = 1, #digits do
         if string_to_number(digits[j], rev) == nil then
            return
         end
      end

      key = rev
   end

   permutate(keys, {}, 0, check_permutation)

   return key
end

local sum = 0
for line in io.lines() do
   local digits = {}
   for word in line:gmatch '%w+' do
      table.insert(digits, word)
   end

   local key = solve(digits)
   local n = 0
   for i = #digits - 3, #digits do
      local d = string_to_number(digits[i], key) * 10 ^ (#digits - i)
      n = n + d
   end
   sum = sum + n
end

print(sum)
