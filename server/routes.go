package main

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/uuid"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

//URI is th domain
const URI = "localhost:8084"

//Session struct describes the session (for DB)
type Session struct {
	ID            uuid.UUID `json: "session_id`
	UserID        uuid.UUID `json: "_id"`
	Refresh       string    `json: "refresh"`
	ExpiresAt     time.Time `json: "expires_at"`
	IsSessionOver bool      `json: "is_session_over"`
}

//SuccessResponseNewAccess struct for describes response in /refresh and /setTokens paths
type SuccessResponseNewAccess struct {
	SessionID uuid.UUID
	Access    *jwt.Token
	ExpiresAt time.Time
}

//RefreshKitFromClient struct describes data sent w request on /refresh
type RefreshKitFromClient struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Refresh   string
}

var response SuccessResponseNewAccess

//Marsh to marshall everything
func Marsh(toM interface{}) (m []byte) {
	m, err := json.Marshal(toM)
	if err != nil {
		log.Fatal("error marshalling", toM, err)
	}
	return m
}

func home(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("This is an auth app.Users' IDs:\n"))

	for _, v := range IDs {
		ctx.Write(Marsh(v))
		ctx.Write([]byte("\n"))
	}
}

func setTokens(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("set tokens\n")
	reqUser := uuid.FromBytesOrNil(ctx.Request.Header.PeekBytes([]byte("id")))
	found, err := findUser(reqUser)
	if err != nil {
		log.Println(err)
	}
	if found == true {
		insertNewSession(createNewTokens(reqUser))
		ctx.Write([]byte("Tokens created\n"))
	} else {
		ctx.Write([]byte("No user with this ID is present is the database\n"))
	}
}

func createNewTokens(userID uuid.UUID) (session Session) {
	var response SuccessResponseNewAccess
	refresh := generateRefresh()
	sessionID := uuid.Must(uuid.NewV4())
	refreshExpiresAt := time.Now().Add(time.Hour * 48)
	isSessionOver := false
	c := fasthttp.AcquireCookie()
	c.SetKey("refreshToken")
	c.SetValue(refresh)
	c.HTTPOnly()
	c.SetExpire(refreshExpiresAt)
	ctx.Response.Header.SetCookie(c)
	response.SessionID = sessionID
	response.Access = generateAccess(userID)
	response.ExpiresAt = time.Now().Add(time.Minute * 30)
	ctx.Write(Marsh(response))
	ctx.Write([]byte("refresh token passed in Cookie"))
	session = Session{
		userID,
		sessionID,
		refresh,
		refreshExpiresAt,
		isSessionOver,
	}
	return session
}

func refresh(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("refresh")
	var reqUser RefreshKitFromClient
	reqUser.UserID = uuid.FromBytesOrNil(ctx.Request.Header.PeekBytes([]byte("id")))
	DecodeFromHeader(ctx, "id", reqUser.UserID)
	DecodeFromHeader(ctx, "refresh", reqUser.Refresh)
	DecodeFromHeader(ctx, "session_id", reqUser.SessionID)
	if reqUser.Refresh == "" {
		ctx.Write([]byte("No token found"))
		ctx.Redirect(URI+"/setTokens", 308)
	} else {
		var sessionInfo = readSessionInfo(reqUser.UserID)
		if time.Now().After(sessionInfo.ExpiresAt) || time.Now().Equal(sessionInfo.ExpiresAt) || sessionInfo.IsSessionOver == true {
			ctx.Write([]byte("Your session is over, please log in again"))
			ctx.Redirect(URI+"/setTokens", 308)
		}

		err := bcrypt.CompareHashAndPassword([]byte(sessionInfo.Refresh), []byte(reqUser.Refresh))
		if err != nil || sessionInfo.ID != reqUser.SessionID {
			ctx.Write([]byte("Your session is over, please log in again"))
			ctx.Redirect(URI+"/setTokens", 308)
		}
		insertNewSession(createNewTokens(reqUser.UserID))
	}
}

//DecodeFromHeader decodes header's value for given key to dest
func DecodeFromHeader(ctx *fasthttp.RequestCtx, key string, dest interface{}) {
	header := ctx.Request.Header.PeekBytes([]byte(key))
	decoder := json.NewDecoder(bytes.NewReader(header))
	err := decoder.Decode(&dest)
	if err != nil {
		log.Fatal("decoder error\n", err)
	}
}

//edit 18:50
func delOne(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("del one")
	//create new session there
	var id uuid.UUID
	var refresh string
	DecodeFromHeader(ctx, "id", id)
	DecodeFromHeader(ctx, "refresh", refresh)
	session := createNewTokens(id)
	delRefresh(refresh, session)

	//this route refreshes and deletes the prev refresh token
}

//edit 18:50
func delAll(ctx *fasthttp.RequestCtx) {
	var id uuid.UUID
	ctx.WriteString("del all")
	DecodeFromHeader(ctx, "id", id)
	delAllRefresh(id)
	//this route deletes all refresh tokens
}
