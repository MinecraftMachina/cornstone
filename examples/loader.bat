@ECHO off

rem ======================================================================================================
rem USER CONFIGURATION
rem ======================================================================================================

SET "MODPACK_NAME=Valhelsia2"
SET "MODPACK_URL=https://github.com/MinecraftMachina/valhelsia-2-corn/archive/master.zip"

SET "CORNSTONE_VERSION=1.0.2"
SET "CORNSTONE_FILE=%CD%\cornstone-%MODPACK_NAME%.exe"
SET "LAUNCHER_DIR=%CD%\corn-%MODPACK_NAME%"

rem ======================================================================================================

SET "CORNSTONE_URL=https://github.com/MinecraftMachina/cornstone/releases/download/v%CORNSTONE_VERSION%/cornstone_%CORNSTONE_VERSION%_windows_amd64.exe"

IF NOT EXIST "%CORNSTONE_FILE%" (
    ECHO Downloading loader...
    powershell -Command "(New-Object Net.WebClient).DownloadFile('%CORNSTONE_URL%', '%CORNSTONE_FILE%')" || GOTO :ERROR
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
ECHO  3 - Add offline account
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
IF NOT EXIST "%LAUNCHER_DIR%" (
    "%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" init || GOTO :ERROR
)
"%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" install -u -n "%MODPACK_NAME%" -i "%MODPACK_URL%" || GOTO :ERROR
pause
GOTO :MENU

:PLAY
"%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" run || GOTO :ERROR
GOTO :EXIT

:OFFLINE
"%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" offline || GOTO :ERROR
pause
GOTO :MENU

:EXIT
exit

:ERROR
ECHO Failed with error %errorlevel%
pause
GOTO :EXIT
