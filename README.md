# gokumail

`gokumail` is a simple `POP3` to `IMAP` proxy which will split your KUmails
between student mails and work mails, and expose the student mails through
`POP3`. This makes it possible to have Gmail fetch all your student mails,
while leaving the work mails at KUs mail servers.

## TODO

- [x] Create `alumni` folder if it doesn't exist
- [x] Add proper logging
- [x] Enable TLS for `POP3` server
- [x] Use external config instead of hardcoded values
- [x] Add user DB
- [ ] Simple webinterface for configuring white/blacklists

## DB table

``` sql
CREATE TABLE IF NOT EXISTS `user_settings` (
    `username` varchar(255) NOT NULL,
    `workmail` varchar(255) NOT NULL,
    `fromwhitelist` varchar(255) NOT NULL,
    `towhitelist` varchar(255) NOT NULL,
    `blacklist` varchar(255) NOT NULL
    PRIMARY KEY (`username`)
);
```

## LICENSE

Copyright (C) 2014  Mikkel Oscar Lyderik Larsen

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
