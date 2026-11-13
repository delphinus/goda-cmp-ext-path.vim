-- A minimal Neovim configuration for the demo: goda.vim, goda-cmp.vim,
-- goda-cmp-ui-pum.vim, goda-cmp-ext-fuzzy.vim and pum.vim from
-- $GODA_DEMO_PLUGINS, and this checkout. Run by demo/record.sh:
--
--   nvim -u demo/init.lua -i NONE demo/sample/notes.txt

local demo = vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":p:h")
local repo = vim.fn.fnamemodify(demo, ":h")
local plugins = assert(vim.env.GODA_DEMO_PLUGINS, "set GODA_DEMO_PLUGINS")
local cache = vim.env.GODA_DEMO_CACHE or (plugins .. "/../cache")

local names = { "goda.vim", "goda-cmp.vim", "goda-cmp-ui-pum.vim", "goda-cmp-ext-fuzzy.vim", "pum.vim" }
for _, name in ipairs(names) do
  vim.opt.rtp:prepend(plugins .. "/" .. name)
end
vim.opt.rtp:prepend(repo)

-- Not asked to the terminal: the recording terminal does not answer.
vim.o.background = vim.env.GODA_DEMO_BACKGROUND or "light"
vim.o.swapfile = false
-- Bundled with Neovim 0.12 and later: Catppuccin Latte with a light background
-- (Mocha with a dark one).
pcall(vim.cmd.colorscheme, vim.env.GODA_DEMO_COLORSCHEME or "catppuccin")
vim.o.termguicolors = true
vim.o.number = true
vim.o.signcolumn = "no"
vim.o.laststatus = 0
vim.o.ruler = false
vim.o.showmode = false

-- Build with the plugins above and this checkout (go.work).
vim.g["goda#workspace"] = {
  plugins .. "/goda.vim",
  plugins .. "/goda-cmp.vim",
  plugins .. "/goda-cmp-ui-pum.vim",
  plugins .. "/goda-cmp-ext-fuzzy.vim",
  repo,
}
vim.g["goda_cmp#config_dir"] = demo .. "/config"
vim.g["goda_cmp#binary"] = cache .. "/goda-cmp"

vim.api.nvim_create_autocmd("VimEnter", {
  once = true,
  callback = function()
    vim.fn["pum#set_option"] {
      -- pum.vim highlights the typed characters as a pattern over the whole
      -- line, the menu included, which does not follow the fuzzy matcher.
      highlight_matches = "",
      max_height = 8,
    }
    if vim.fn.executable(vim.g["goda_cmp#binary"]) == 0 then
      vim.fn["goda_cmp#build"]()
    end
    vim.fn["goda_cmp#enable"]()
    if vim.env.GODA_DEMO_DRIVE then
      dofile(demo .. "/drive.lua")
    end
  end,
})

vim.keymap.set("i", "<C-n>", "<Cmd>call pum#map#select_relative(+1)<CR>")
vim.keymap.set("i", "<C-p>", "<Cmd>call pum#map#select_relative(-1)<CR>")
vim.keymap.set("i", "<C-y>", "<Cmd>call pum#map#confirm()<CR>")
vim.keymap.set("i", "<C-e>", "<Cmd>call pum#map#cancel()<CR>")
