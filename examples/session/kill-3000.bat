@echo off
set PORT=3000

echo Checking port %PORT%...

:: Find the PID running on the specified port
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :%PORT% ^| findstr LISTENING') do (
    set PID=%%a
)

:: If a PID was found, kill it
if defined PID (
    echo Found process running on port %PORT% with PID: %PID%
    echo Killing process...
    taskkill /F /PID %PID%
    echo Port %PORT% has been cleared successfully.
) else (
    echo No process found listening on port %PORT%.
)

pause