### shell脚本

---

### 编译 git_merge
```bash
go build -ldflags="-s -w" -o ./bin/git_merge.exe ./git_merge/main.go
```

### 编译 go_cmd
```bash
go build -ldflags="-s -w" -o ./bin/go_cmd.exe ./go_cmd/main.go
```

### 编译 google_sheet
```bash
go build -ldflags="-s -w" -o ./bin/google_sheet.exe ./google_sheet/main.go
```

### 编译 go_shell
```bash
go build -ldflags="-s -w" -o ./bin/go_shell.exe ./go_shell/main.go
```