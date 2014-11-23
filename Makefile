GO=go

all: gokumail

gokumail: clean
	$(GO) build

install: gokumail
	# bin
	install -Dm755 gokumail $(DESTDIR)/usr/bin/
	# templates
	install -Dm644 views/* $(DESTDIR)/usr/share/gokumail/views/
	# static
	install -Dm644 static/* $(DESTDIR)/usr/share/gokumail/static/
	# config
	install -Dm644 gokumail.conf $(DESTDIR)/etc/gokumail.conf
	# service
	install -Dm644 contrib/gokumail.service $(DESTDIR)/usr/lib/systemd/system/

clean:
	-@rm -f gokumail
