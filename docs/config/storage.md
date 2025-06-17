---
icon: material/database
---
<span class="badge badge-red" title="If this option is required or optional">required</span>

The `storage` option is used to configure how and where data is stored.

## `backend`
<span class="badge badge-purple" title="Value Type">enum</span>
<span class="badge badge-blue" title="Default Value">badger</span>
<span class="badge badge-orange" title="If this option is required or optional">recommended</span>

The `backend` option is used to set which storage backend should be used to store data. This defines how data is 
stored. Depending on the chosen backend different further configuration options might be available or not.

In the following the supported storage backend and their available options are detailed.

### `json`
If `backend` is set to `json` the JSON files backend is used. This backend stores data in simple json files (in 
multiple directories). This option is great to see which data is stored, since it is the most human-readable storage 
format supported. It is also great if data is manipulated externally.
Performance-wise other options are better.

??? file "config.yaml"

    ```yaml
    storage:
        backend: json
        data_dir: /path/to/data
    ```

The following configuration options are defined for the `json` backend:

#### `data_dir`
<span class="badge badge-purple" title="Value Type">directory path</span>
<span class="badge badge-red" title="If this option is required or optional">required</span>

The `data_dir` option sets the root directory where data is stored on disk. LightHouse creates subdirectories and 
places the JSON files in those directories.

### `badger`
If `backend` is set to `badger` the [BadgerDB](https://github.com/hypermodeinc/badger) backend is used. BadgerDB is 
an embeddable, persistent key-value database. No external dependencies are needed, as BadgerDB is embedded into 
LightHouse. The data is stored on disk is not suitable to be read or manipulated by humans.

??? file "config.yaml"

    ```yaml
    storage:
        backend: badger
        data_dir: /path/to/data
    ```

The following configuration options are defined for the `badger` backend:

#### `data_dir`
<span class="badge badge-purple" title="Value Type">directory path</span>
<span class="badge badge-red" title="If this option is required or optional">required</span>

The `data_dir` option sets the root directory where the badger data is stored on disk. 
