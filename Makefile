CC=gcc

LIB_OBJ = \
    internal/quirc/lib/decode.o \
    internal/quirc/lib/identify.o \
    internal/quirc/lib/quirc.o \
    internal/quirc/lib/version_db.o

libquirc.a: $(LIB_OBJ)
	ar cru $@ $^
	ranlib $@