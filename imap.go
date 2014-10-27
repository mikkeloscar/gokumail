package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/mikkeloscar/goimap"
)

const (
	kuImapServer     = "exchange.ku.dk"
	kuImapServerPort = 993
	kuUsernameFmt    = "%s@ku.dk"
	kuAlumniFmt      = "%s@alumni.ku.dk"
)

// KUmail defines an special IMAP client for KUmail
type KUmail struct {
	User   string
	Pass   string
	client *imap.IMAPClient
	conf   *Config
}

// MsgInfo defines a struct to hold a message ID and the corresponding message
// size
type MsgInfo struct {
	ID   string
	size int // message size in octets
}

// MsgUID defines a struct to hold a message ID and the UID of that message
type MsgUID struct {
	ID  string
	UID int
}

// Init setup a connection, authenticate with IMAP server and organize mails.
// After this, the server will be ready to send the mails requested from the
// subfolder
// assumes User and Pass has been initialized in k
func (k *KUmail) Init(config *Config) bool {
	k.conf = config
	alumniMail := fmt.Sprintf(kuAlumniFmt, k.User)
	k.conf.ToWhitelist = append(k.conf.ToWhitelist, alumniMail)
	service := fmt.Sprintf("%s:%d", kuImapServer, kuImapServerPort)

	conn, err := net.Dial("tcp", service)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false
	}

	client, errr := imap.NewClient(conn, kuImapServer)
	if errr != nil {
		fmt.Printf("error: %s\n", errr)
		return false
	}

	k.client = client
	user := fmt.Sprintf(kuUsernameFmt, k.User)

	err = k.client.Login(user, k.Pass)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false
	}

	// Organize Mails just after login
	err = k.organizeMails()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false
	}

	return true
}

// Close logout of IMAP session and close connection
func (k *KUmail) Close() {
	k.client.Logout()
	k.client.Close()
}

func (k *KUmail) organizeMails() error {
	k.client.Select("INBOX")

	uids, err := k.search()
	if err != nil {
		return err
	}

	// TODO make alumni configurable
	return k.moveMails(uids, "INBOX", "INBOX/alumni")
}

func (k *KUmail) search() (map[string]string, error) {
	uids := make(map[string]string)

	// TO
	for _, to := range k.conf.ToWhitelist {
		e, err := k.searchHeader("TO", to)
		if err != nil {
			return nil, err
		}
		addToMap(&uids, e)
	}
	// Received
	for _, to := range k.conf.ToWhitelist {
		e, err := k.searchHeader("Received", to)
		if err != nil {
			return nil, err
		}
		addToMap(&uids, e)
	}
	// From
	for _, from := range k.conf.FromWhitelist {
		e, err := k.searchHeader("FROM", from)
		if err != nil {
			return nil, err
		}
		addToMap(&uids, e)
	}

	return uids, nil
}

// Add all elements of a slice to map m
func addToMap(m *map[string]string, s []string) {
	for _, val := range s {
		(*m)[val] = val
	}
}

func (k *KUmail) moveMails(msgUIDs map[string]string, src string, dst string) error {
	moved := 0
	notMoved := 0

	for _, uid := range msgUIDs {
		if k.validateMail(uid) {
			k.moveMail(uid, src, dst)
			moved++
		} else {
			notMoved++
		}
	}

	// expunge after moving all mails and marking them Deleted in INBOX
	_, err := k.client.Expunge()
	if err != nil {
		return err
	}

	fmt.Printf("Moved %d of %d possible mails\n", moved, moved+notMoved)
	return nil
}

func (k *KUmail) moveMail(msgUID string, src string, dst string) error {
	err := k.client.Copy(msgUID, dst)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}

	// TODO handle err in StoreAddFlag
	return k.client.StoreAddFlag(msgUID, imap.Deleted)
}

// Make sure that the mail was not sent to work mail account
func (k *KUmail) validateMail(msgUID string) bool {
	fields := "BODY.PEEK[HEADER.FIELDS (FROM TO CC)]"
	resp, err := k.client.Fetch(msgUID, fields)
	if err != nil {
		return false
	}

	body := resp.Body

	return hasSubstring(body, k.conf.Whitelist) || !hasSubstring(body, k.conf.Blacklist)
}

// check if some element of slice l is a substring of s
func hasSubstring(s string, l []string) bool {
	for _, elem := range l {
		if strings.Contains(s, elem) {
			return true
		}
	}
	return false
}

func (k *KUmail) searchHeader(header string, query string) ([]string, error) {
	resp, err := k.client.Search(fmt.Sprintf("(HEADER %s \"%s\")", header, query))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}
	return resp, nil
}

// ListAll lists all the messages in the alumni folder in KUmail
func (k *KUmail) ListAll() ([]*MsgInfo, int, error) {
	k.client.Select("INBOX/alumni")

	resp, err := k.client.Search("ALL")
	if err != nil {
		return []*MsgInfo{}, 0, err
	}

	msgs := make([]*MsgInfo, len(resp))

	total := 0

	for i, id := range resp {
		res, err := k.client.GetMessageSize(id)
		if err != nil {
			return []*MsgInfo{}, 0, err
		}
		total += res

		msgs[i] = &MsgInfo{id, res}
	}

	return msgs, total, nil
}

// UIDL lists all the messages in the alumni folder along with there UID
func (k *KUmail) UIDL() ([]*MsgUID, error) {
	k.client.Select("INBOX/alumni")

	resp, err := k.client.Search("ALL")
	if err != nil {
		return []*MsgUID{}, err
	}

	msgs := make([]*MsgUID, len(resp))

	for i, id := range resp {
		res, err := k.client.Fetch(id, "UID")
		if err != nil {
			return []*MsgUID{}, err
		}

		msgs[i] = &MsgUID{id, res.Value}
	}

	return msgs, nil
}

func (k *KUmail) GetMessage(id string) (string, int, error) {
	k.client.Select("INBOX/alumni")

	resp, err := k.client.Fetch(id, imap.RFC822)
	if err != nil {
		return "", 0, err
	}

	return resp.Body, resp.Length, nil
}
