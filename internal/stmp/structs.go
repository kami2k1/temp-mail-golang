package stmp

import (
    "sync"


    "github.com/emersion/go-imap/client"
)

type IMAPConfig struct {
    client *client.Client
}


type Attc struct {
    ID       string `json:"_id"`
    Filename string `json:"filename"`
    Size     int64  `json:"size"`
    MIME     string `json:"mimetype"`
    CID      string `json:"cid,omitempty"`
}

type MAIL struct {
    UID        uint32    `json:"uid"`
    From       string    `json:"from"`
    To         string    `json:"to"`
    Subject    string    `json:"subject"`
    Body       string    `json:"bodyHtml,omitempty"`
    Bodypw     string    `json:"bodyPreview,omitempty"`
    ReceivedAt int64     `json:"receivedAt"`
    AttcCount  int       `json:"attachmentsCount"`
    Attc       []Attc    `json:"attachments"`
}


type attachmentBlob struct {
    Attc
    Data []byte
}

var (
    IMAP     = &IMAPConfig{}
    TARGET   = []string{}
    DATAMAIL = make(map[string][]MAIL) // key: recipient email (To)


    attachments   = make(map[string]attachmentBlob)
    seenUID       = make(map[string]map[uint32]bool) 
    storeMu       sync.RWMutex
    nextAttachSeq uint64

   
    lastMessagesHit int64
)

