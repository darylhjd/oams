x64 = amd64
env = windows

all:
	@echo "Creating Intervention Service Build";
	@echo "Environment: $(env), Architecture: $(x64)";
	env GOOS=$(env) GOARCH=$(x64) go build -o main.exe;