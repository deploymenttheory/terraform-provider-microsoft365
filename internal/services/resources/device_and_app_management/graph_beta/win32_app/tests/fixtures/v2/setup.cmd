@echo off
if not exist "%ProgramData%\TerraformWin32Acceptance" mkdir "%ProgramData%\TerraformWin32Acceptance"
echo v2>"%ProgramData%\TerraformWin32Acceptance\version.txt"
exit /b 0
