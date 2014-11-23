GO=go

all: gokumail

gokumail: clean
	$(GO) build

install:
	# bin
	install -Dm755 gokumail $(DESTDIR)/usr/bin/gokumail
	# templates
	install -d $(DESTDIR)/usr/share/gokumail/views/
	install -m644 views/* $(DESTDIR)/usr/share/gokumail/views/
	# static
	install -d $(DESTDIR)/usr/share/gokumail/static/css
	install -d $(DESTDIR)/usr/share/gokumail/static/js
	install -Dm644 static/css/* $(DESTDIR)/usr/share/gokumail/static/css/
	install -Dm644 static/js/* $(DESTDIR)/usr/share/gokumail/static/js/
	# config
	install -Dm644 gokumail.conf $(DESTDIR)/etc/gokumail.conf
	# service
	install -d $(DESTDIR)/usr/lib/systemd/system/
	install -m644 contrib/gokumail.service $(DESTDIR)/usr/lib/systemd/system/

clean:
	-@rm -f gokumail
