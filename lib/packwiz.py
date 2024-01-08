import subprocess

MODRINTH = 0
CURSEFORGE = 1

class Packwiz:
    def run(dir: str, *args: str) -> int:
        return subprocess.run(["packwiz", *args], cwd=dir, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL).returncode
    
    def add(dir: str, source: int, mod: str) -> bool:
        if source == MODRINTH:
            return Packwiz.run(dir, "modrinth", "add", mod, "-y") == 0
        elif source == CURSEFORGE:
            return Packwiz.run(dir, "curseforge", "add", mod, "-y") == 0
    
        return False

    def init(dir: str, loader: str, minecraft: str) -> bool:
        if loader == "fabric":
            arg = "--fabric-latest"
        elif loader == "forge":
            arg = "--forge-latest"
        elif loader == "quilt":
            arg = "--quilt-latest"
        
        return Packwiz.run(
            dir,
            "init",
            arg,
            "--modloader",
            loader,
            "--version",
            "0.0.0",
            "--author",
            "HappyRedstone Modding",
            "--mc-version",
            minecraft,
            "--name",
            "Sodium Plus"
        ) == 0
