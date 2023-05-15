import os
import requests

from pyrinth.project import versions

project = "sodiumplus"

def write_file(path: str, content: str) -> None:
    with open(path, "w") as file:
        file.write(content + "\n")

def append_file(path: str, content: str) -> None:
    with open(path, "a") as file:
        file.write(content + "\n")

def get_versions(project: str) -> list[versions.ProjectVersion]:
    return versions.get_versions(project)

def main():
    versions = get_versions(project)
    changelogs = [(version.name, version.changelog) for version in versions]
    
    for (name, changelog) in changelogs:
        head = f"# v{name}\n"
        body = changelog + "\n"
        
        append_file("CHANGELIST.md", head)
        append_file("CHANGELIST.md", body)

if __name__ == "__main__":
    main()
