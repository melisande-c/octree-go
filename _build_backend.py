
"""Custom build backend to mark wheels as platform-specific while preserving setuptools-scm."""

from setuptools import build_meta as _orig

# Re-export all the standard build backend functions
prepare_metadata_for_build_wheel = _orig.prepare_metadata_for_build_wheel
build_wheel = _orig.build_wheel
build_sdist = _orig.build_sdist
get_requires_for_build_wheel = _orig.get_requires_for_build_wheel
get_requires_for_build_sdist = _orig.get_requires_for_build_sdist

# Re-export optional hooks that setuptools-scm might use
try:
    build_editable = _orig.build_editable
except AttributeError:
    pass

try:
    get_requires_for_build_editable = _orig.get_requires_for_build_editable
except AttributeError:
    pass

try:
    prepare_metadata_for_build_editable = _orig.prepare_metadata_for_build_editable
except AttributeError:
    pass


def _patch_distribution():
    """Patch setuptools to mark this package as having platform-specific content."""
    from setuptools.dist import Distribution
    
    def has_ext_modules(self):
        # Return True to generate platform-specific wheel tags
        # (e.g., cp310-cp310-macosx_11_0_arm64.whl instead of py3-none-any.whl)
        return True
    
    Distribution.has_ext_modules = has_ext_modules


# Apply the patch when the module is imported
_patch_distribution()