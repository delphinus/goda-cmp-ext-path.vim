-- Types the demo into demo/sample/notes.txt, then quits (GODA_DEMO_DRIVE=1).

local steps = {
  { 1500, "G" },
  { 300, "o" },
}

local function input(text)
  for c in text:gmatch "." do
    table.insert(steps, { 180, c })
  end
end

-- The whole names with their modification times, docs/ out of "dc", then
-- getting-started.md out of "gs".
input "See ./"
vim.list_extend(steps, { { 2200, "d" }, { 180, "c" }, { 1500, "<C-n>" }, { 1200, "<C-y>" } })
vim.list_extend(steps, { { 2200, "g" }, { 180, "s" }, { 1500, "<C-n>" }, { 1200, "<C-y>" } })
input " and ./."
vim.list_extend(steps, { { 2500, "<Esc>" } })

local function feed(keys)
  vim.api.nvim_feedkeys(vim.api.nvim_replace_termcodes(keys, true, false, true), "t", false)
end

local function run(i)
  if i > #steps then
    vim.defer_fn(function()
      if vim.env.GODA_DEMO_LOG then
        vim.fn.writefile(vim.fn["goda#log"](), vim.env.GODA_DEMO_LOG)
      end
      vim.cmd "qa!"
    end, 1500)
    return
  end
  vim.defer_fn(function()
    feed(steps[i][2])
    run(i + 1)
  end, steps[i][1])
end

-- Clear what startup said (the recording terminal answers no queries).
vim.cmd "echo '' | redraw"
run(1)
