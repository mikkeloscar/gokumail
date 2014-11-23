GO=go

all: gokumail

gokumail: clean
	$(GO) build

install: gokumail
	# bin
	install -Dm755 gokumail $(DESTDIR)/usr/bin/gokumail
	# templates
	install -d $(DESTDIR)/usr/share/gokumail/views/
	install -m644 views/* $(DESTDIR)/usr/share/gokumail/views/
	# static
	install -d $(DESTDIR)/usr/share/gokumail/static/
	install -m644 static/* $(DESTDIR)/usr/share/gokumail/static/
	# config
	install -Dm644 gokumail.conf $(DESTDIR)/etc/gokumail.conf
	# service
	install -d $(DESTDIR)/usr/lib/systemd/system/
	install -m644 contrib/gokumail.service $(DESTDIR)/usr/lib/systemd/system/

clean:
	-@rm -f gokumail
