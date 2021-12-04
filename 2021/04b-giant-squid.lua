local numbers = nil
local boards = {}

for line in io.lines() do
   if numbers == nil then
      numbers = {}
      for x in line:gmatch '%d+' do
         table.insert(numbers, tonumber(x))
      end
   end

   local last_board = #boards == 0 and nil or boards[#boards]
   if last_board == nil or #last_board == 5 then
      table.insert(boards, {})
      last_board = boards[#boards]
   end

   local row = {}

   for x in line:gmatch '%d+' do
      table.insert(row, tonumber(x))
   end

   if #row == 5 then
      table.insert(last_board, row)
   end
end

local order = {}

for i, n in ipairs(numbers) do
   order[n] = i
end

local scored = {}

for i = 1, #boards do
   local b = boards[i]
   local s = #numbers

   for j = 1, 5 do
      local t = 0

      for k = 1, 5 do
         if t < order[b[j][k]] then
            t = order[b[j][k]]
         end
      end

      if t < s then
         s = t
      end

      t = 0

      for k = 1, 5 do
         if t < order[b[k][j]] then
            t = order[b[k][j]]
         end
      end

      if t < s then
         s = t
      end
   end

   table.insert(scored, {
      score = s,
      board = b,
   })
end

local best_score = 0
local best_board = nil

for _, b in ipairs(scored) do
   if best_score < b.score then
      best_score = b.score
      best_board = b.board
   end
end

local selected_numbers = {}
for i = 1, best_score do
   selected_numbers[numbers[i]] = 1
end

local unmarked = 0

for j = 1, 5 do
   for k = 1, 5 do
      if selected_numbers[best_board[j][k]] == nil then
         unmarked = unmarked + best_board[j][k]
      end
   end
end

print(unmarked * numbers[best_score])
