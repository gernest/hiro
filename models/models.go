package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"

	valid "github.com/asaskevich/govalidator"
	"github.com/gernest/hiro/keys"
	"github.com/gernest/hiro/resource"
	uuid "github.com/satori/go.uuid"
)

// MaxLimit max limit of list of items to be returned on each database query.
const MaxLimit = 10

// QR information associated with the qr code.
type QR struct {
	UUID           uuid.UUID      `json:"uuid,omitempty"`
	URL            string         `json:"url,omitempty"`
	Name           string         `json:"name,omitempty"`
	ImageData      []byte         `json:"image_data,omitempty"`
	UserID         uuid.UUID      `json:"-"`
	ImageURL       string         `json:"image_url,omitempty"`
	Context        []*ContextItem `json:"context,omitempty"`
	ShouldRedirect bool           `json:"should_redirect,omitempty"`
	RedirectURL    string         `json:"redirect_url,omitempty"`
	Private        bool           `json:"private,omitempty"`
	CreatedAt      time.Time      `json:"created_at,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at,omitempty"`
	Groups         []*Collection  `json:"groups,omitempty"`
	id             string
}

func (q *QR) DownloadName() string {
	downloadName := fmt.Sprintf("bq_%s", q.UUID)
	if q.Name != "" {
		downloadName = slug.Make(q.Name)
	}
	return downloadName
}

func (q *QR) ID() string {
	if q.id != "" {
		return q.id
	}
	q.id = q.UUID.String()
	return q.id
}

type Collection struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	AccountID uuid.UUID `json:"account_id"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Collection) ID() string {
	return c.UUID.String()
}

type CollectionList struct {
	Collections []*Collection `json:"collections"`
}

// QRList is an object with list of qrcodes.
type QRList struct {
	QRCodes []*QR        `json:"qrcodes"`
	Options *ListOptions `json:"options"`
	// Total is the total number of the qrcodes stored  belonging to the user.
	Total int64 `json:"total"`
}

// Range calls f on every grcode in the list. If the call returns false
// execution is stopped.
func (list *QRList) Range(f func(*QR) bool) {
	for i := 0; i < len(list.QRCodes); i++ {
		if !f(list.QRCodes[i]) {
			return
		}
	}
}

// Token is a jwt token
type Token struct {
	UUID uuid.UUID

	// Issuer the user issuing the token.
	Issuer        uuid.UUID
	IssuerAccount *Account

	// Subject the user who will benefit from the issued token.
	Subject        uuid.UUID
	SubjectAccount *Account

	ExpiresOn time.Time
	NotBefore time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewToken returns new *Token with uuid set.
func NewToken(issuer, subject uuid.UUID) *Token {
	now := time.Now()
	return &Token{
		UUID:      uuid.NewV4(),
		ExpiresOn: now.Add(keys.MaxTokenLife),
		NotBefore: now,
		Issuer:    issuer,
		Subject:   subject,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// DefaultLoginToken returns default token for a logged in user.
func DefaultLoginToken(user uuid.UUID) *Token {
	return NewToken(user, user)
}

// JWT is a struct for generating and parsing jwt tokens using hmac signing
// method.
type JWT struct {
	Secret []byte
}

// Generate generate a jwt token with HMAC signing method.

// Item is a keyvalue pair.
type Item struct {
	Key   interface{}
	Value interface{}
}

// Account is a bq user account
type Account struct {
	UUID      uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListOptions struct {
	Limit  int
	Offset int
}

// Login json object for login.
type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// CreateAccount form for creating new account.
type CreateAccount struct {
	Email           string `schema:"email" json:"email"`
	Password        string `schema:"password" json:"password"`
	ConfirmPassword string `schema:"confirm_password" json:"confirm_password"`
}

// Message is the error message.
//
// swagger:model message
type Message struct {
	Resource string `json:"resource,omitempty"`
	Field    string `json:"field,omitempty"`
	Desc     string `json:"desc,omitempty"`
}

// APIError is an error returned from api calls.
//
// swagger:response apiError
type APIError struct {
	Message string    `json:"message"`
	Errors  []Message `json:"errors,omitempty"`
}

func (a *APIError) Error() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// Add appens msg to errors.
func (a *APIError) Add(msg Message) {
	a.Errors = append(a.Errors, msg)
}

// LoginRes json object for login response.
type LoginRes struct {
	Token string `json:"token"`
}

// Status os the response status in json.
type Status struct {
	Status string `json:"status"`
}

// QRReq contains fields needed to generate a new qr code.
//
// swagger:parameters qrReq createQRCode
type QRReq struct {

	// Width is the size in width for the generated png image of the qe code
	//
	// in: body
	// required: false
	Width int64 `json:"width,omitempty"`

	// Height is the height dimension of the generated png file of the qr code
	//
	// in:body
	// required: false
	Height int64 `json:"height,omitempty"`

	ShouldRedirect bool           `json:"should_redirect"`
	RedirectURL    string         `json:"redirect_url,omitempty"`
	Name           string         `json:"name"`
	Groups         []string       `json:"groups,omitempty"`
	Context        []*ContextItem `json:"context,omitempty"`
}

// CollectionReq is an object reqpresenting options for creating a new collection.
type CollectionReq struct {

	// Name is a collection name it doesn't have to be unique
	Name string `json:"name,omitempty"`

	// Color is a hex color string which aid to visually identify a collection.
	Color string `json:"color,omitempty"`
}

// CollectionAssignReq is an object with options for assigning a qrcode to a
// given collection.
type CollectionAssignReq struct {

	// ID  is the uuid for the collection we want to assign the qrcode into.
	ID uuid.UUID `json:"collection_id"`

	// QRID is the uuid of the qrcode. A signle qrccode can be assigned to as many
	// collections as the user wants.
	QRID uuid.UUID `json:"qr_id"`
}

// SortBy sorts the wrapped collections by field.
func (c *CollectionList) SortBy(field string) {
	switch field {
	case "updated_at":
		sort.Slice(c.Collections, func(i, j int) bool {
			return c.Collections[i].UpdatedAt.Before(
				c.Collections[j].UpdatedAt,
			)
		})
	case "created_at":
		sort.Slice(c.Collections, func(i, j int) bool {
			return c.Collections[i].CreatedAt.Before(
				c.Collections[j].CreatedAt,
			)
		})
	}
}

func (c *CollectionList) Filter(name string) []*Collection {
	if name == "" {
		return c.Collections
	}
	return c.filter(func(k string) bool {
		return strings.Contains(k, name)
	})
}

func (c *CollectionList) filter(f func(string) bool) []*Collection {
	var out []*Collection
	for _, v := range c.Collections {
		if f(v.Name) {
			out = append(out, v)
		}
	}
	return out
}

// QRCodeCacheRequest json object for a request to the qrcode image data.
type QRCodeCacheRequest struct {
	UUID   uuid.UUID `json:"uuid"`
	Owner  uuid.UUID `json:"owner"`
	Width  int       `json:"width"`
	Height int       `json:"height"`
}

func (q *QRCodeCacheRequest) Key() string {
	return url.QueryEscape(q.key())
}
func (q *QRCodeCacheRequest) key() string {
	return fmt.Sprintf("usr=%s&id=%s&w=%d&h=%d",
		q.Owner, q.UUID, q.Width, q.Height,
	)
}

func (q *QRCodeCacheRequest) FromKey(key string) error {
	key, err := url.QueryUnescape(key)
	if err != nil {
		return err
	}
	parts := strings.Split(key, "&")
	if len(parts) != 4 {
		return errors.New("bad key")
	}
	usr, err := uuid.FromString(strings.Split(parts[0], "=")[1])
	if err != nil {
		return err
	}
	id, err := uuid.FromString(strings.Split(parts[1], "=")[1])
	if err != nil {
		return err
	}
	w, err := strconv.Atoi(strings.Split(parts[2], "=")[1])
	if err != nil {
		return err
	}
	h, err := strconv.Atoi(strings.Split(parts[3], "=")[1])
	if err != nil {
		return err
	}
	q.Owner = usr
	q.UUID = id
	q.Width = w
	q.Height = h
	return nil
}

type ContextItemKind string

const (
	ItemBool   ContextItemKind = "bool"
	ItemInt    ContextItemKind = "int64"
	ItemFloat  ContextItemKind = "float64"
	ItemString ContextItemKind = "string"
)

type ContextItem struct {
	Kind  ContextItemKind `json:"kind"`
	Key   string          `json:"key"`
	Value interface{}     `json:"value"`
}

func (c *ContextItem) Valid() bool {
	if c.Value == nil {
		return false
	}
	switch c.Kind {
	case ItemBool:
		_, ok := c.Value.(bool)
		return ok
	case ItemInt:
		_, ok := c.Value.(int64)
		return ok
	case ItemFloat:
		_, ok := c.Value.(float64)
		return ok
	case ItemString:
		_, ok := c.Value.(string)
		return ok
	default:
		return false
	}
}

func (c *ContextItem) GetValue() string {
	return fmt.Sprint(c.Value)
}

// Validate returns true if the form is valid. Different rules are applied on
// deterrent fields.
//	name
//		- can not be empty
//		- can not be one of the reserved names.
//		- should only be a combination of characters and numbers.
//	email
//		- can not be empty
//		- should be a valid email address
//	password
//		- can not be empty
//		- TODO: set minimum number of characters to encourage strong password.
//	confirm_password
//		- must match the given  password.
func (c *CreateAccount) Validate() *APIError {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)
	c.ConfirmPassword = strings.TrimSpace(c.ConfirmPassword)
	a := &APIError{
		Message: keys.FailedValidation,
	}

	if c.Email == "" {
		a.Add(Message{
			Resource: resource.Account,
			Field:    "email",
			Desc:     keys.MissingEmail,
		})
	} else {
		if !valid.IsEmail(c.Email) {
			a.Add(Message{
				Resource: resource.Account,
				Field:    "email",
				Desc:     keys.InvalidEmail,
			})
		} else {
			n, _ := valid.NormalizeEmail(c.Email)
			c.Email = n
		}
	}
	if c.Password == "" {
		a.Add(Message{
			Resource: resource.Account,
			Field:    "password",
			Desc:     keys.MissingPassword,
		})
	} else {
		if c.Password != c.ConfirmPassword {
			a.Add(Message{
				Resource: resource.Account,
				Field:    "confirm_password",
				Desc:     keys.ConfirmPasswordMismatch,
			})
		}
	}
	if len(a.Errors) > 0 {
		return a
	}
	return nil
}
