name: nvim conf
action:
    install vim:
        dir: ~
        cmd: curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim.appimage
        - wait
        cmd: chmod u+x nvim.appimage
        cmd: ./nvim.appimage --appimage-extract 
        cmd: wait
        cmd: mv squashfs-root /
        cmd: ln -s /squashfs-root/AppRun /usr/bin/nvim
    packer: 
        dir: ~
        cmd: git clone --depth 1 https://github.com/wbthomason/packer.nvim .local/share/nvim/site/pack/packer/start/packer.nvim
        cmd: wait
    preconfig:
        dir: ~
        create config folder: mkdir -p .config/nvim
        initialize config module:
            dir: ~/.config/nvim
            cmd: touch init.lua
            cmd: echo "require('mxk')" >> init.lua
        create lua module: 
            dir: ~/.config/nvim
            cmd: mkdir -p lua/mxk 
            cmd: cp /app/nvim_conf/lua/mxk/packer.lua lua/mxk/
            cmd: touch lua/mxk/init.lua 
            cmd: echo "require('mxk.packer')" >> lua/mxk/init.lua
        sync packer packages: nvim --headless -c "so" -c "PackerSync" -c "q"
        wait: wait
    config:
        dir: ~
        copy entire config:
            cmd: cp -a /app/nvim_conf/after .config/nvim/
            cmd: cp -a /app/nvim_conf/lua .config/nvim/
        sync once again:
            dir: ~
            cmd: nvim --headless -c "so" -c "PackerSync" -c "q"
