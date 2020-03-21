package bot

import (
	"io"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TelegramBotExt interface {
	TelegramBot
	ConditionalHandle(endpoint interface{}, handler interface{}, condFunc interface{})
}

type TelegramBot interface {
	Handle(endpoint interface{}, handler interface{})
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
	SendAlbum(to tb.Recipient, a tb.Album, options ...interface{}) ([]tb.Message, error)
	Reply(to *tb.Message, what interface{}, options ...interface{}) (*tb.Message, error)
	Forward(to tb.Recipient, what *tb.Message, options ...interface{}) (*tb.Message, error)
	Edit(message tb.Editable, what interface{}, options ...interface{}) (*tb.Message, error)
	EditReplyMarkup(message tb.Editable, markup *tb.ReplyMarkup) (*tb.Message, error)
	EditCaption(message tb.Editable, caption string) (*tb.Message, error)
	EditMedia(message tb.Editable, inputMedia tb.InputMedia, options ...interface{}) (*tb.Message, error)
	Delete(message tb.Editable) error
	Notify(recipient tb.Recipient, action tb.ChatAction) error
	Accept(query *tb.PreCheckoutQuery, errorMessage ...string) error
	Answer(query *tb.Query, response *tb.QueryResponse) error
	Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error
	FileByID(fileID string) (tb.File, error)
	Download(file *tb.File, localFilename string) error
	GetFile(file *tb.File) (io.ReadCloser, error)
	StopLiveLocation(message tb.Editable, options ...interface{}) (*tb.Message, error)
	GetInviteLink(chat *tb.Chat) (string, error)
	SetGroupTitle(chat *tb.Chat, newTitle string) error
	SetGroupDescription(chat *tb.Chat, description string) error
	SetGroupPhoto(chat *tb.Chat, p *tb.Photo) error
	SetGroupStickerSet(chat *tb.Chat, setName string) error
	DeleteGroupPhoto(chat *tb.Chat) error
	DeleteGroupStickerSet(chat *tb.Chat) error
	Leave(chat *tb.Chat) error
	Pin(message tb.Editable, options ...interface{}) error
	Unpin(chat *tb.Chat) error
	ChatByID(id string) (*tb.Chat, error)
	ProfilePhotosOf(user *tb.User) ([]tb.Photo, error)
	ChatMemberOf(chat *tb.Chat, user *tb.User) (*tb.ChatMember, error)
	FileURLByID(fileID string) (string, error)
	UploadStickerFile(userID int, pngSticker *tb.File) (*tb.File, error)
	GetStickerSet(name string) (*tb.StickerSet, error)
	CreateNewStickerSet(sp tb.StickerSetParams, containsMasks bool, maskPosition tb.MaskPosition) error
	AddStickerToSet(sp tb.StickerSetParams, maskPosition tb.MaskPosition) error
	SetStickerPositionInSet(sticker string, position int) error
	DeleteStickerFromSet(sticker string) error
	Start()
	Stop()
}
