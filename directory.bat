@echo off
cd /d "%~dp0"
tree /f > structure.txt
echo 项目目录结构已保存到 structure.txt
pause
