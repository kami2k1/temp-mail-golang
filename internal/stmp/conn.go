package stmp

import (
	"crypto/tls"
	"fmt"
	"io"
	"kami/internal/config"
	"log"
	"mime"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	msgmail "github.com/emersion/go-message/mail"
)

func readMailbox(c *client.Client, mboxName string, lastN uint32) error {
	mbox, err := c.Select(mboxName, false)
	if err != nil {
		ConnectIMAP()
		return fmt.Errorf("select %s: %w", mboxName, err)

	}
	if mbox.Messages == 0 {

		return nil
	}

	fromSeq := uint32(1)
	if mbox.Messages > lastN {
		fromSeq = mbox.Messages - lastN + 1
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(fromSeq, mbox.Messages)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchInternalDate, imap.FetchUid, section.FetchItem()}
	messages := make(chan *imap.Message, lastN)
	done := make(chan error, 1)
	go func() { done <- c.Fetch(seqset, items, messages) }()

	for msg := range messages {
		if msg.Envelope == nil {
			continue
		}

		if len(msg.Envelope.To) < 1 {
			continue
		}
		toAddr := msg.Envelope.To[0]

		if toAddr == nil || toAddr.HostName == "gmail.com" {
			continue
		}
		key := toAddr.MailboxName + "@" + toAddr.HostName

		if msg.Uid != 0 {
			if !MarkSeenUID(mboxName, msg.Uid) {
				continue
			}
		}

		var bodyReader io.Reader
		if r := msg.GetBody(section); r != nil {
			bodyReader = r
		}

		var (
			attachmentsMeta []Attc
			bodyPreview     string
			bodyHTML        string
		)

		if bodyReader != nil {
			mr, err := msgmail.CreateReader(bodyReader)
			if err == nil {
				for {
					p, err := mr.NextPart()
					if err == io.EOF {
						break
					}
					if err != nil {
						break
					}
					switch h := p.Header.(type) {
					case *msgmail.AttachmentHeader:

						filename, _ := h.Filename()
						ctype, _, _ := h.ContentType()
						cid := strings.Trim(h.Header.Get("Content-Id"), "<>")

						data, _ := io.ReadAll(p.Body)
						size := int64(len(data))

						id := fmt.Sprintf("att-%d-%d", msg.Uid, time.Now().UnixNano())
						att := Attc{ID: id, Filename: filename, Size: size, MIME: ctype, CID: cid}
						attachmentsMeta = append(attachmentsMeta, att)
						PutAttachment(attachmentBlob{Attc: att, Data: data})
					case *msgmail.InlineHeader:

						ctype, params, _ := h.ContentType()
						mediatype, _, _ := mime.ParseMediaType(ctype)
						b, _ := io.ReadAll(p.Body)
						if strings.HasPrefix(mediatype, "text/") {
							txt := string(b)
							if strings.Contains(strings.ToLower(mediatype), "html") || mediatype == "text/html" {
								if bodyHTML == "" {
									bodyHTML = txt
								}
							} else {
								if bodyPreview == "" {
									if charset := params["charset"]; charset != "" {

									}
									if len(txt) > 200 {
										bodyPreview = txt[:200]
									} else {
										bodyPreview = txt
									}
								}
							}
						} else if strings.HasPrefix(mediatype, "image/") {

							data := b
							size := int64(len(data))
							cid := strings.Trim(h.Header.Get("Content-Id"), "<>")

							filename := ""
							if disp := h.Header.Get("Content-Disposition"); disp != "" {
								if _, p, err := mime.ParseMediaType(disp); err == nil {
									if fn := p["filename"]; fn != "" {
										filename = fn
									}
								}
							}
							if filename == "" {
								if _, p, err := mime.ParseMediaType(h.Header.Get("Content-Type")); err == nil {
									if fn := p["name"]; fn != "" {
										filename = fn
									}
								}
							}
							if filename == "" {
								if loc := h.Header.Get("Content-Location"); loc != "" {
									filename = path.Base(loc)
								}
							}

							if filename == "" {
								filename = cid
							}
							if filename == "" {
								filename = "inline"
							}
							if !strings.Contains(filename, ".") {
								ext := ""
								switch mediatype {
								case "image/jpeg":
									ext = "jpg"
								case "image/png":
									ext = "png"
								case "image/gif":
									ext = "gif"
								case "image/webp":
									ext = "webp"
								}
								if ext != "" {
									filename = filename + "." + ext
								}
							}
							id := fmt.Sprintf("att-%d-%d", msg.Uid, time.Now().UnixNano())
							att := Attc{ID: id, Filename: filename, Size: size, MIME: mediatype, CID: cid}
							attachmentsMeta = append(attachmentsMeta, att)
							PutAttachment(attachmentBlob{Attc: att, Data: data})
						}
					}
				}
			}
		}

		from := ""
		if len(msg.Envelope.From) > 0 && msg.Envelope.From[0] != nil {
			from = msg.Envelope.From[0].MailboxName + "@" + msg.Envelope.From[0].HostName
		}
		receivedAt := msg.InternalDate.Unix()

		if bodyHTML != "" {
			replaced := bodyHTML
			for _, a := range attachmentsMeta {
				if a.CID != "" {
					replaced = strings.ReplaceAll(replaced, "cid:"+a.CID, "/attachments/"+a.ID)
				}
			}
			bodyHTML = replaced
		}

		m := MAIL{
			UID:        msg.Uid,
			From:       from,
			To:         key,
			Subject:    msg.Envelope.Subject,
			Body:       bodyHTML,
			Bodypw:     bodyPreview,
			ReceivedAt: receivedAt,
			AttcCount:  len(attachmentsMeta),
			Attc:       attachmentsMeta,
		}

		PutMail(key, m)
	}

	if err := <-done; err != nil {
		return fmt.Errorf("fetch %s: %w", mboxName, err)
	}
	return nil
}
func GetMailinbox() {
	ticker := time.NewTicker(2 * time.Second)
	for {
		<-ticker.C

		if !ActiveWithin(3 * time.Minute) {
			continue
		}
		for _, box := range TARGET {
			if err := readMailbox(IMAP.client, box, 20); err != nil {
				log.Println("error:", err)
			}
		}
	}
}

func ConnectSTMP() {

}
func ConnectIMAP() {
	host := fmt.Sprintf("%s:%s", config.Config.IMAP_HOST, config.Config.IMAP_PORT)
	fmt.Print("Connecting to IMAP server... ")
	fmt.Printf("%s\n", host)
	tlsconfig := &tls.Config{ServerName: config.Config.IMAP_HOST}

	c, err := client.DialTLS(host, tlsconfig)
	if err != nil {
		log.Println("dial:", err)

	}

	err = c.Login(config.Config.STMP_USER, config.Config.STMP_PASS)
	if err != nil {
		log.Println("login: fal", err)

	}
	IMAP.client = c
	log.Println("Connected to IMAP server")
	boxes, err := listMailboxes(c)
	if err != nil {
		log.Fatal("list:", err)
	}
	all := findByAttr(boxes, `\All`)
	spam := findByAttr(boxes, `\Junk`)
	trash := findByAttr(boxes, `\Trash`)

	if all != "" {
		TARGET = append(TARGET, all)
	}
	if spam != "" {
		TARGET = append(TARGET, spam)
	}
	if trash != "" {
		TARGET = append(TARGET, trash)
	}
	fmt.Printf("Reading mailboxes: %v\n", TARGET)

}

func listMailboxes(c *client.Client) ([]*imap.MailboxInfo, error) {
	ch := make(chan *imap.MailboxInfo, 50)
	done := make(chan error, 1)
	go func() { done <- c.List("", "*", ch) }()

	var boxes []*imap.MailboxInfo
	for m := range ch {
		boxes = append(boxes, m)
	}
	return boxes, <-done
}
func findByAttr(boxes []*imap.MailboxInfo, attr string) string {
	attr = strings.ToUpper(attr)
	for _, b := range boxes {
		for _, a := range b.Attributes {
			if strings.ToUpper(a) == attr {
				return b.Name
			}
		}
	}
	return ""
}

// Helpers to interact with stores
func PutMail(key string, m MAIL) {
	storeMu.Lock()
	defer storeMu.Unlock()
	// Append newest at the front
	list := DATAMAIL[key]
	// If duplicate UID already present, skip
	for _, ex := range list {
		if ex.UID == m.UID {
			return
		}
	}
	list = append([]MAIL{m}, list...)
	DATAMAIL[key] = list
}

func MarkSeenUID(mbox string, uid uint32) bool {
	storeMu.Lock()
	defer storeMu.Unlock()
	set, ok := seenUID[mbox]
	if !ok {
		set = make(map[uint32]bool)
		seenUID[mbox] = set
	}
	if set[uid] {
		return false
	}
	set[uid] = true
	return true
}

func PutAttachment(a attachmentBlob) {
	storeMu.Lock()
	defer storeMu.Unlock()
	attachments[a.ID] = a
}

func GetAttachment(id string) (attachmentBlob, bool) {
	storeMu.RLock()
	defer storeMu.RUnlock()
	a, ok := attachments[id]
	return a, ok
}

// Activity touch and check
func TouchActivity() { atomic.StoreInt64(&lastMessagesHit, time.Now().Unix()) }
func ActiveWithin(d time.Duration) bool {
	ts := atomic.LoadInt64(&lastMessagesHit)
	if ts == 0 {
		return false
	}
	return time.Since(time.Unix(ts, 0)) <= d
}
