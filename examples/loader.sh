#!/bin/bash

# ======================================================================================================
# USER CONFIGURATION
# ======================================================================================================

MODPACK_NAME="Valhelsia 2"
MODPACK_URL="https://github.com/MinecraftMachina/valhelsia-2-corn/archive/master.zip"

# ======================================================================================================

CORNSTONE_VERSION="1.0.1"
CORNSTONE_URL="https://github.com/MinecraftMachina/cornstone/releases/download/v${CORNSTONE_VERSION}/cornstone_${CORNSTONE_VERSION}_linux_amd64"

if [ ! -f "cornstone" ]; then
    echo "Downloading loader..."
    curl -L -s "$CORNSTONE_URL" -o "cornstone"
    chmod +x cornstone
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
    echo " 3 - Play offline"
    echo " 4 - Exit"
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
            EXIT ;;
        *) MENU ;;
    esac
}

function INSTALL {
    if [ ! -d "MultiMC" ]; then
        $PWD/cornstone multimc -m "MultiMC" init || ERROR
    fi
    $PWD/cornstone multimc -m "MultiMC" install -u -n "$MODPACK_NAME" -i "$MODPACK_URL" || ERROR
    MENU
}

function PLAY {
    $PWD/cornstone multimc -m "MultiMC" run || ERROR
    EXIT
}

function OFFLINE {
    $PWD/cornstone multimc -m "MultiMC" offline || ERROR
    PLAY
}

function EXIT {
    exit 1
}

function ERROR {
    echo "Failed with error $?"
    read -n1 -r -p "Press any key to continue..." key
    EXIT
}

MENU
