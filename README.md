# gazer
Hot reload, exception recovery, and other controls for locally running multi-project monorepos.

## Config: `gazer.yaml`
```yaml
shell:
  exec: [shell executable]
  args:
    - <optional>
gaze-at:
  - name: [process name]
    cmd: [command to run process]
    watch:
      - <files/directories to watch, restart on changes>
    color: <optional>
    restart: <optional | always, retry>
    retries: <optional | int>
```

Examaple
```yaml
shell:
  exec: cmd
  args:
    - /C
gaze-at:
  - name: frontend
    cmd: cd ./frontend && ng serve
    color: blue
    restart: always
  - name: server
    cmd: npm run start:dev --prefix ./backend/src
    watch:
      - ./backend/src
      - ./backend/package.json
    restart: retry
    retries: 3
```