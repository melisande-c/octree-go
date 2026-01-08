import os
import subprocess
import platform
import sys
from pathlib import Path
from typing import Any

from hatchling.builders.hooks.plugin.interface import BuildHookInterface


class GoBuildHook(BuildHookInterface):
    """Custom build hook to compile Go code into a shared library."""
    
    PLUGIN_NAME = "go-build"
    
    def initialize(self, version: str, build_data: dict[str, Any]) -> None:
        """
        Compile Go code and add the shared library to the wheel.
        """
        build_data['pure_python'] = False
        build_data["infer_tag"] = True

        # Get configuration from pyproject.toml
        go_source = self.config.get("go-source", ".")
        go_package = self.config.get("go-package", "export.go")
        output_name = self.config.get(
            "output-name", "octree"
        )
        target_dir = self.config.get("target-dir", "pyoctree/_lib")
        
        # Determine the shared library extension based on platform
        if sys.platform == "linux":
            lib_ext = "so"
        elif sys.platform == "darwin":
            lib_ext = "dylib"
        else:
            raise RuntimeError(f"Unsupported platform: {sys.platform}")

        machine = platform.machine()
        if machine == "x86_64":
            machine = "amd64"
        elif machine == "aarch64":
            machine = "arm64"
        if machine not in ["arm64", "amd64"]:
            raise RuntimeError(
                f"Unsupported CPU '{machine}'. Currently only arm64 and amd64 is supported."
            )
        lib_name = f"{output_name}-{platform.system().lower()}-{machine}.{lib_ext}"
        
        # Paths
        project_root = Path(self.root)
        go_source_path = project_root / go_source / go_package
        target_path = project_root / target_dir
        output_path = target_path / lib_name
        
        # Ensure target directory exists
        target_path.mkdir(parents=True, exist_ok=True)
        
        print(f"Compiling Go code from {go_source_path}")
        print(f"Output will be written to {output_path}")
        
        # Compile Go code to shared library
        try:
            cmd = [
                # "CGO_ENABLED=1",
                "go",
                "build",
                "-buildmode=c-shared",
                "-o",
                str(output_path),
                str(go_source_path),
            ]
            
            result = subprocess.run(
                cmd,
                check=True,
                capture_output=True,
                text=True,
                cwd=project_root,
            )
            
            print(f"Successfully compiled {lib_name}")
            
            # Remove the .h file that Go generates (optional)
            h_file = output_path.with_suffix(".h")
            if h_file.exists():
                h_file.unlink()
                
        except subprocess.CalledProcessError as e:
            print(f"Error compiling Go code: {e.stderr}", file=sys.stderr)
            raise
        except FileNotFoundError:
            print("Go compiler not found. Please install Go.", file=sys.stderr)
            raise
        
        # Add the shared library to the wheel
        # The force_include tells Hatch to include this file in the wheel
        if "force_include" not in build_data:
            build_data["force_include"] = {}
        
        # Map the compiled library to its location in the wheel
        relative_output = str(output_path.relative_to(project_root))
        build_data["force_include"][relative_output] = f"pyoctree/_lib/{lib_name}"
        
        print(f"Added {lib_name} to wheel")