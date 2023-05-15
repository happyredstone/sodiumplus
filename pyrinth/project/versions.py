import os
import json

from typing import Any
from pyrinth import api

class VersionDependency:
    version_id: str
    project_id: str
    file_name: str
    dependency_type: str
    
    def __init__(self, data: Any):
        self.version_id = data["version_id"]
        self.project_id = data["project_id"]
        self.file_name = data["file_name"]
        self.dependency_type = data["dependency_type"]

class VersionFile:
    hashes: dict[str, str]
    url: str
    filename: str
    primary: bool
    size: int
    file_type: str
    
    def __init__(self, data: Any):
        self.hashes = data["hashes"]
        self.url = data["url"]
        self.filename = data["filename"]
        self.primary = data["primary"]
        self.size = data["size"]
        self.file_type = data["file_type"]

class ProjectVersion:
    name: str
    version_number: str
    changelog: str
    dependencies: list[VersionDependency]
    game_versions: list[str]
    version_type: str
    loaders: list[str]
    featured: bool
    status: str
    requested_status: str
    id: str
    project_id: str
    author_id: str
    date_published: str
    downloads: int
    changelog_url: str | None
    files: list[VersionFile]
    
    def __init__(self, data: Any):
        self.name = data["name"]
        self.version_number = data["version_number"]
        self.changelog = data["changelog"]
        self.dependencies = [VersionDependency(dep) for dep in data["dependencies"]]
        self.game_versions = data["game_versions"]
        self.version_type = data["version_type"]
        self.loaders = data["loaders"]
        self.featured = data["featured"]
        self.status = data["status"]
        self.requested_status = data["requested_status"]
        self.id = data["id"]
        self.project_id = data["project_id"]
        self.author_id = data["author_id"]
        self.date_published = data["date_published"]
        self.downloads = data["downloads"]
        self.changelog_url = data["changelog_url"]
        self.files = [VersionFile(file) for file in data["files"]]

def get_versions(project: str) -> list[ProjectVersion]:
    route = f"project/{project}/version"
    text = api.get_resource(route)
    
    return [ProjectVersion(item) for item in json.loads(text)]
