#!/bin/sh

sudo chown -R $(whoami): /home/vscode/.cache
sudo chown -R $(whoami): /home/vscode/.config/gh
# setup github-cli completion 
sudo gh completion -s zsh | sudo tee /usr/local/share/zsh/site-functions/_gh > /dev/null
