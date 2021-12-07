local input = io.read '*a'
local state = {}

for d in input:gmatch '%d+' do
   d = tonumber(d)
   state[d] = (state[d] or 0) + 1
end

for _ = 1, 80 do
   local next_gen = {}
   for t, count in pairs(state) do
      if t == 0 then
         next_gen[6] = (next_gen[6] or 0) + count
         next_gen[8] = count
      else
         next_gen[t - 1] = (next_gen[t - 1] or 0) + count
      end
   end
   state = next_gen
end

local total = 0

for _, count in pairs(state) do
   total = total + count
end

print(total)
