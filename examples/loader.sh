#!/bin/bash

# ======================================================================================================
# USER CONFIGURATION
# ======================================================================================================

MODPACK_NAME="Valhelsia2"
MODPACK_URL="https://github.com/MinecraftMachina/valhelsia-2-corn/archive/master.zip"
MODPACK_SERVER_URL="https://github.com/MinecraftMachina/valhelsia-2-corn/archive/server.zip"

CORNSTONE_VERSION="1.2.1"
CORNSTONE_FILE="$PWD/cornstone-$MODPACK_NAME"
LAUNCHER_DIR="$PWD/corn-$MODPACK_NAME"

# ======================================================================================================

if [[ "$OSTYPE" == "darwin"* ]]; then
    CORNSTONE_OS="darwin"
else
    CORNSTONE_OS="linux"
fi
CORNSTONE_URL="https://github.com/MinecraftMachina/cornstone/releases/download/v${CORNSTONE_VERSION}/cornstone_${CORNSTONE_VERSION}_${CORNSTONE_OS}_amd64"

if [ ! -f "$CORNSTONE_FILE" ]; then
    echo "Downloading loader..."
    curl -L -s "$CORNSTONE_URL" -o "$CORNSTONE_FILE" || ERROR
    chmod +x "$CORNSTONE_FILE" || ERROR
fi

function MENU {
    clear

    echo ""
    echo "..............................................."
    echo " $MODPACK_NAME Launcher, powered by cornstone"
    echo "..............................................."
    echo ""
    echo " 1 - Install or update"
    echo " 2 - Play"
    echo " 3 - Add offline account"
    echo " 4 - Reset"
    echo " 5 - Install or update server"
    echo " 6 - Exit"
    echo ""
    read -p $'Type a number then press ENTER: ' CHOICE

    clear

    case $CHOICE in
        "1")
            INSTALL ;;
        "2")
            PLAY ;;
        "3")
            OFFLINE ;;
        "4")
            RESET ;;
        "5")
            SERVER ;;
        "6")
            EXIT ;;
        *) ;;
    esac
}

function INSTALL {
    if [ ! -d "$LAUNCHER_DIR" ]; then
        "$CORNSTONE_FILE" multimc -m "$LAUNCHER_DIR" init || ERROR
    fi
    "$CORNSTONE_FILE" multimc -m "$LAUNCHER_DIR" install -n "$MODPACK_NAME" -i "$MODPACK_URL" || ERROR
    pause
}

function PLAY {
    "$CORNSTONE_FILE" multimc -m "$LAUNCHER_DIR" run || ERROR
    EXIT
}

function OFFLINE {
    "$CORNSTONE_FILE" multimc -m "$LAUNCHER_DIR" offline || ERROR
    pause
}

function RESET {
    echo ""
    echo "WARNING: This will delete the modpack with all your data!"
    echo ""
    pause
    rm -rf "$LAUNCHER_DIR" || ERROR
    INSTALL
}

function SERVER {
    "$CORNSTONE_FILE" server -s "$SERVER_DIR" install -i "$MODPACK_SERVER_URL" || ERROR
    pause
}

function EXIT {
    exit 1
}

function ERROR {
    echo "Failed with error $?"
    pause
    EXIT
}

function pause {
    read -n1 -r -p "Press any key to continue..." key
}

while true
do
    MENU
done
