

windows:
	env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_LDFLAGS="-L/usr/local/Cellar/mingw-w64/6.0.0_1/toolchain-x86_64/x86_64-w64-mingw32/lib -lSDL2" CGO_CFLAGS="-I/usr/local/Cellar/mingw-w64/6.0.0_1/toolchain-x86_64/x86_64-w64-mingw32/include -D_REENTRANT" go build -x