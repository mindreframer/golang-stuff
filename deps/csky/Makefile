################################################################################
# Variables
################################################################################

MAJVER=0
MINVER=3
RELVER=0

PREFIX?=/usr/local
CFLAGS=-g -O2 -Wall -Wextra -std=c99 -Iinclude -fPIC -Wno-pointer-to-int-cast
LDFLAGS=

SOURCES=$(wildcard src/*.c)
OBJECTS=$(patsubst %.c,%.o,${SOURCES})
TEST_SOURCES=$(wildcard tests/*_tests.c)
TEST_OBJECTS=$(patsubst %.c,%,${TEST_SOURCES})

UNAME=$(shell uname)
ifeq ($(UNAME), Darwin)
LDFLAGS+=-dynamiclib -undefined dynamic_lookup -fPIC -compatibility_version $(MAJVER).$(MINVER) -current_version $(MAJVER).$(MINVER).$(RELVER)
SONAME=libcsky.dylib
SONAME_VER1=libcsky.${MAJVER}.dylib
SONAME_VER2=libcsky.${MAJVER}.${MINVER}.dylib
endif
ifeq ($(UNAME), Linux)
SONAME=libcsky.so
SONAME_VER1=libcsky.so.${MAJVER}
SONAME_VER2=libcsky.so.${MAJVER}.${MINVER}
LDFLAGS+=-shared -Wl,-soname,${SONAME}.${MAJVER} -rdynamic -fPIC -o ${SONAME}.${MAJOR_VERSION}.${MINOR_VERSION}
endif


################################################################################
# Targets
################################################################################

default: build

build: libcsky.a ${SONAME_VER2}

libcsky.a: ${OBJECTS}
	rm -f bin/libsky.a
	ar rcs libcsky.a ${OBJECTS}
	ranlib libcsky.a

${SONAME_VER2}: ${OBJECTS}
	$(CXX) ${LDFLAGS} ${OBJECTS} -o ${SONAME_VER2}

install: build
	install -d $(DESTDIR)/$(PREFIX)/include/sky
	install include/sky/* $(DESTDIR)/$(PREFIX)/include/sky
	install include/sky.h $(DESTDIR)/$(PREFIX)/include/sky.h
	install -d $(DESTDIR)/$(PREFIX)/lib
	install libcsky.a $(DESTDIR)/$(PREFIX)/lib
	install ${SONAME_VER2} $(DESTDIR)/$(PREFIX)/lib
	ln -sf $(DESTDIR)/$(PREFIX)/lib/${SONAME_VER2} $(DESTDIR)/$(PREFIX)/lib/${SONAME}
	ln -sf $(DESTDIR)/$(PREFIX)/lib/${SONAME_VER2} $(DESTDIR)/$(PREFIX)/lib/${SONAME_VER1}


clean: 
	rm -rf bin ${OBJECTS} ${TEST_OBJECTS} tmp
	rm -rf tests/*.dSYM tests/**/*.dSYM
	rm -rf tests/*.o tests/**/*.o
	rm -rf libcsky.*


################################################################################
# Tests
################################################################################

check: $(TEST_OBJECTS)
	@sh ./tests/runtests.sh

$(TEST_OBJECTS): %: %.c build
	$(CC) $(CFLAGS) -Itests -o $@ $< libcsky.a
