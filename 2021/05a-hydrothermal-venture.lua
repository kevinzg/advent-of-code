local lines = {}

for line in io.lines() do
   local next = line:gmatch '%d+'
   local p = {
      x = 1 + next(),
      y = 1 + next(),
   }
   local q = {
      x = 1 + next(),
      y = 1 + next(),
   }
   table.insert(lines, {
      p = p,
      q = q,
   })
end

local function dir(a, b)
   if a > b then
      return -1
   elseif b > a then
      return 1
   else
      return 0
   end
end

local map = {}

local mark_map = function(x, y)
   local key = x .. ',' .. y
   if map[key] == nil then
      map[key] = 1
   else
      map[key] = map[key] + 1
   end
end

for _, line in ipairs(lines) do
   local dx = dir(line.p.x, line.q.x)
   local dy = dir(line.p.y, line.q.y)

   -- Ignore diagonal lines
   if dx == 0 or dy == 0 then
      local x = line.p.x
      local y = line.p.y

      while x ~= line.q.x or y ~= line.q.y do
         mark_map(x, y)
         x = x + dx
         y = y + dy
      end

      mark_map(x, y)
   end
end

local count = 0

for _, v in pairs(map) do
   if v > 1 then
      count = count + 1
   end
end

print(count)
