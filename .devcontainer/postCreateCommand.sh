#!/bin/sh

# setup github-cli completion 
sudo gh completion -s zsh | sudo tee /usr/local/share/zsh/site-functions/_gh > /dev/null
