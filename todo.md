- [ ] Parse the configuration file.

The config file look like the file below. 

```yaml
source: /home/aristide/Downloads # the data source, the root folder to listen for changes.
destination:
    - path: /home/aristide/Videos # The same as file:///home/aristide/Videos  
      pattern: ["*.mp4", "*.flv", "*.avi", "*.mov", "*.wmn", "*.webm", "*.mkv"]

    - path: file:///home/aristide/Documents
      pattern: ["*.pdf", "*.epub", "*.djvu", "*.mobi", "*.doc", "*.docx", "*.odt"]

    - path: sftp://aristide@server.com/Images
      pattern: ["*.png", "*.jpg", "*.jpeg", "*.svg", "*gif"]
```

- [ ] Fill in the routing logic.

- [ ] Fill in the workers logic.

- [ ] Add support of database logging.

- [ ] Add an api on top of database logs.
