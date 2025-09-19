#!/bin/bash

gh repo edit --enable-projects=false
gh repo edit --enable-merge-commit=false
gh repo edit --enable-rebase-merge=false
gh repo edit --delete-branch-on-merge=true
gh repo edit --enable-wiki=false
gh repo edit --enable-discussions=false
