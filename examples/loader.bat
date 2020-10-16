@ECHO off

rem ======================================================================================================
rem CONFIGURATION
rem ======================================================================================================

SET "CORNSTONE_URL=https://github.com/MinecraftMachina/cornstone/releases/download/v1.0.0/cornstone_1.0.0_windows_amd64.exe"
SET "MODPACK_NAME=Valhelsia 2"
SET "MODPACK_URL=https://github.com/MinecraftMachina/valhelsia-2-corn/archive/master.zip"

rem ======================================================================================================

IF NOT EXIST "cornstone.exe" (
    ECHO Downloading loader...
    powershell -Command "(New-Object Net.WebClient).DownloadFile('%CORNSTONE_URL%', '%~dp0cornstone.exe')" || GOTO :ERROR
)

:MENU
CLS

ECHO.
ECHO ...............................................
ECHO  %MODPACK_NAME% Launcher, powered by cornstone
ECHO ...............................................
ECHO.
ECHO  1 - Install or update
ECHO  2 - Play
ECHO  3 - Play offline
ECHO  4 - Exit
ECHO.

SET /P M=Type a number then press ENTER: 
CLS

IF %M%==1 GOTO :INSTALL
IF %M%==2 GOTO :PLAY
IF %M%==3 GOTO :OFFLINE
IF %M%==4 GOTO :EXIT
GOTO :MENU

:INSTALL
IF NOT EXIST "MultiMC" (
    cornstone multimc -m "MultiMC" init || GOTO :ERROR
)
cornstone multimc -m "MultiMC" install -u -n "%MODPACK_NAME%" -i "%MODPACK_URL%" || GOTO :ERROR
GOTO :MENU

:PLAY
cornstone multimc -m "MultiMC" run || GOTO :ERROR
GOTO :EXIT

:OFFLINE
cornstone multimc -m "MultiMC" offline || GOTO :ERROR
GOTO :PLAY

:EXIT
exit

:ERROR
ECHO Failed with error #%errorlevel%
pause
exit /b %errorlevel%
