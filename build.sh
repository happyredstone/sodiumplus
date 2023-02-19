#!/bin/bash

set -e

# Update submodules
git submodule update --init --recursive

# Litematica build
cd litematica
./gradlew build
cp build/libs/litematica-*.jar ../mods
cd ..

packwiz refresh
packwiz modrinth export
packwiz curseforge export
