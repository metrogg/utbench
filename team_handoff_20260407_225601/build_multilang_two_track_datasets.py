from build_cpp_module_level_dataset import main as build_cpp
from build_go_two_track_datasets import main as build_go
from build_java_module_level_dataset import main as build_java
from build_javascript_two_track_datasets import main as build_javascript
from build_python_module_level_dataset import main as build_python_module
from build_self_contained_python_dataset import main as build_python_self_contained
from sync_external_self_contained_sources import main as sync_external_sources


def main():
    sync_external_sources()
    build_python_self_contained()
    build_python_module()
    build_go()
    build_java()
    build_cpp()
    build_javascript()
    print("=== All two-track datasets finished ===")


if __name__ == "__main__":
    main()
