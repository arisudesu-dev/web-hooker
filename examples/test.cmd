@echo off
echo Status: 200 OK
echo Content-type: text/plain
echo.
echo Hello from cmd.exe
echo.
echo Script is: %0
echo.
echo Working dir is: %cd%
echo.
echo Environment is:
set | sort
echo.
