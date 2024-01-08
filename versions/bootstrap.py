import os
import sys

sys.path.append("..")

from lib.packwiz import Packwiz, MODRINTH, CURSEFORGE

mods = [
    "roughly-searchable",
    "debugify",
    "fabricskyboxes",
    "fabricskyboxes-interop",
    "mixintrace",
    "borderless-mining",
    "dynamic-fps",
    "immediatelyfast",
    "wakes",
    "fastquit",
    "optigui",
    "roughly-enough-professions-rep",
    "e4mc",
    "ferritecore",
    "fpsdisplay",
    "language-reload",
    "rei",
    "fabrishot",
    "lazy-language-loader",
    "eating-animation",
    "chat-heads",
    "jade",
    "presence-footsteps",
    "antighost",
    "blur-fabric",
    "held-item-info",
    "betterf3",
    "capes",
    "xaeros-minimap",
    "xaeros-world-map",
    "simple-voice-chat",
    "dripsounds-fabric",
    "entity-model-features",
    "full-brightness-toggle",
    "modelfix",
    "fallingleaves",
    "appleskin",
    "cit-resewn",
    "sodium",
    "morechathistory",
    "shulkerboxtooltip",
    "yosbr",
    "lambdabettergrass",
    "modmenu",
    "animatica",
    "lambdynamiclights",
    "not-enough-animations",
    "xaeros-world-map",
    "fadeless",
    "no-chat-reports",
    "cull-leaves",
    "stack-to-nearby-chests",
    "sound-physics-remastered",
    "entity-view-distance",
    "dashloader",
    "durability-tooltip",
    "hold-that-chunk",
    "sodium-extra",
    "entitytexturefeatures",
    "iris",
    "moreculling",
    "indium",
    "wavey-capes",
    "advancementinfo",
    "better-mount-hud",
    "main-menu-credits",
    "zoomify",
    "visuality",
    "calcmod",
    "mouse-wheelie",
    "continuity",
    "reeses-sodium-options",
    "starlight",
    "boat-item-view",
    "entityculling",
    "curse/litematica",
]

dir = sys.argv[2] + "/" + sys.argv[1]
dir = os.path.abspath(dir)

if not os.path.exists(dir):
    os.mkdir(dir)

print("[INFO] Initializing pack...")

Packwiz.init(dir, sys.argv[1], sys.argv[2])

failed = []

for mod in mods:
    if mod.startswith("curse/"):
        mod = mod[6:]
        source = CURSEFORGE
    else:
        source = MODRINTH
    
    print(f"[INFO] Adding {mod}...")
    
    if not Packwiz.add(dir, source, mod):
        failed.append(mod)

for mod in failed:
    print(f"** [ERROR] Failed to add: {mod}")
