@ECHO off

rem ======================================================================================================
rem USER CONFIGURATION
rem ======================================================================================================

SET "MODPACK_NAME=Valhelsia2"
SET "MODPACK_URL=https://github.com/MinecraftMachina/valhelsia-2-corn/archive/master.zip"
SET "MODPACK_SERVER_URL=https://github.com/MinecraftMachina/valhelsia-2-corn/archive/server.zip"

SET "CORNSTONE_VERSION=1.1.0"
SET "CORNSTONE_FILE=%CD%\cornstone-%MODPACK_NAME%.exe"
SET "LAUNCHER_DIR=%CD%\corn-%MODPACK_NAME%"
SET "SERVER_DIR=%CD%\corn-%MODPACK_NAME%-server"

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
ECHO  4 - Reset
ECHO  5 - Install or update server
ECHO  6 - Exit
ECHO.

SET /P M=Type a number then press ENTER: 
CLS

IF %M%==1 GOTO :INSTALL
IF %M%==2 GOTO :PLAY
IF %M%==3 GOTO :OFFLINE
IF %M%==4 GOTO :RESET
IF %M%==5 GOTO :SERVER
IF %M%==6 GOTO :EXIT
GOTO :MENU

:INSTALL
IF NOT EXIST "%LAUNCHER_DIR%" (
    "%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" init || GOTO :ERROR
)
"%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" install -n "%MODPACK_NAME%" -i "%MODPACK_URL%" || GOTO :ERROR
pause
GOTO :MENU

:PLAY
"%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" run || GOTO :ERROR
GOTO :EXIT

:OFFLINE
"%CORNSTONE_FILE%" multimc -m "%LAUNCHER_DIR%" offline || GOTO :ERROR
pause
GOTO :MENU

:RESET
ECHO.
ECHO WARNING: This will delete the modpack with all your data!
ECHO.
pause
rd /s /q "%LAUNCHER_DIR%" || GOTO :ERROR
GOTO :INSTALL

:SERVER
"%CORNSTONE_FILE%" server -s "%SERVER_DIR%" install -i "%MODPACK_SERVER_URL%" || GOTO :ERROR
pause
GOTO :MENU

:EXIT
exit

:ERROR
ECHO Failed with error %errorlevel%
pause
GOTO :EXIT
