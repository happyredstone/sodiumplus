# CF: /projects/update.curseforge.project-id
# MR: /mod/update.modrinth.mod-id

import os
import toml

rp_toml_files = [file for file in os.listdir("resourcepacks") if file.endswith(".toml")]
mod_toml_files = [file for file in os.listdir("mods") if file.endswith(".toml")]

def write_file(path: str, content: str) -> None:
    with open(path, "w") as file:
        file.write(content + "\n")

def append_file(path: str, content: str) -> None:
    with open(path, "a") as file:
        file.write(content + "\n")

def get_files(root: str, files: list[str]):
    for file in files:
        with open(root + file, "r") as file:
            content = file.read()
            data = toml.loads(content)
            
            mod_name = data["name"]
            update_config = data["update"]
            
            if "curseforge" in update_config:
                project_id = update_config["curseforge"]["project-id"]
                mod_url = f"https://www.curseforge.com/projects/{project_id}"
            elif "modrinth" in update_config:
                mod_id = update_config["modrinth"]["mod-id"]
                mod_url = f"https://modrinth.com/mod/{mod_id}"
            
            line = f"- [{mod_name}]({mod_url})"
            
            append_file("MODLIST.md", line)

def main():
    write_file("MODLIST.md", "# Mod List\n")
    
    get_files("resourcepacks/", rp_toml_files)
    get_files("mods/", mod_toml_files)

if __name__ == "__main__":
    main()
