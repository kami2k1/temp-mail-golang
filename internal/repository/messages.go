package repository

import (
    "sort"

    "kami/internal/stmp"
)

type Message struct {
    Mailbox string      `json:"mailbox"`
    Data    []stmp.MAIL `json:"data"`
}

func GetMessages(mailbox string) Message {
    list := stmp.DATAMAIL[mailbox]
    
    sort.SliceStable(list, func(i, j int) bool {
        ti := list[i].ReceivedAt
        tj := list[j].ReceivedAt
        if ti == 0 && tj == 0 {
            return list[i].UID > list[j].UID
        }
        if tj == 0 {
            return true
        }
        if ti == 0 {
            return false
        }
        return ti > tj
    })
    
    out := make([]stmp.MAIL, len(list))
    copy(out, list)
    for i := range out {
       
        out[i].Body = ""
   
        out[i].Attc = nil
    }
    return Message{Mailbox: mailbox, Data: out}
}


func NewEmptyMessage(mailbox string) Message {
    return Message{Mailbox: mailbox, Data: []stmp.MAIL{}}
}


func GetMessageByUID(mailbox string, uid uint32) (stmp.MAIL, bool) {
    list := stmp.DATAMAIL[mailbox]
    for _, m := range list {
        if m.UID == uid {
            return m, true
        }
    }
    return stmp.MAIL{}, false
}
