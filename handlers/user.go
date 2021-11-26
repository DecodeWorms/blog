package handlers

import (
	"blog/storage"
	"blog/types"
	"blog/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	user storage.User
}

func NewUserHandler(use storage.User) UserHandler {
	return UserHandler{
		user: use,
	}

}

var TokenString string
var newToken string

func (u UserHandler) AutoMigrate(w http.ResponseWriter, r *http.Request) {
	var err error
	var data types.User
	err = u.user.Automigrate(data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to create users table")
	}
	json.NewEncoder(w).Encode("Users table created successfully...")
}

func (u UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var data types.User
	var err error
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to decode JSON")
	}

	err = u.user.Create(data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to save data")
	}
	json.NewEncoder(w).Encode("Resources are created successfully")
}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	util.SetHeader(w)
	var data types.User
	var d types.User
	var err error

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to decode json")
		return
	}
	json.NewEncoder(w).Encode("Decoding json..")

	d, err = u.user.Login(data)
	if err != nil {
		json.NewEncoder(w).Encode("username or password incorrect")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(d.Password), []byte(data.Password))
	if err != nil {
		json.NewEncoder(w).Encode("Unable to hash password")
		return
	}

	res := &types.TokenDetails{}
	res, err = createToken(data.Username, data.ID, w)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to create token string")
	}
	TokenString = res.AccessToken
	//This http.SetCookie is not working and cookie is not set at the client(browser)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   res.AccessToken,
		Expires: res.AtExp,
	})
	json.NewEncoder(w).Encode(TokenString)
}

func (u UserHandler) MyProfiles(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var err error

	var claims *types.Claim
	claims, err = verifyToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode("Unauthorized token or you are being logged out")
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	var data []types.User
	data, err = u.user.MyProfiles(claims.Username)
	json.NewEncoder(w).Encode(data)

}

func (u UserHandler) Posts(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var err error
	var data []types.Post
	var claims *types.Claim
	claims, err = verifyToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode("Unauthorized token or you are being logged out")
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
	data, err = u.user.Posts()
	if err != nil {
		json.NewEncoder(w).Encode("Unable to load posts")
	}
	json.NewEncoder(w).Encode(data)

}

func (u UserHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var err error
	var data types.User
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to marshal json")
		return
	}
	var clm *types.Claim
	clm, err = verifyToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode("Unauthorized token or you are being logged out")
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", clm.Username)))
	err = u.user.UpdateName(clm.Username, data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to update user profile")
		return
	}
	json.NewEncoder(w).Encode("Username Changed successfully")

}

func (u UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	var err error
	var clm *types.Claim
	clm, err = verifyToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode("Unauthorized token or you are being logged out")
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", clm.Username)))
	fmt.Println(clm.UserId)

	var data types.Post
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to unmarshal json or empty input")
	}

	err = u.user.Post(uint64(clm.UserId), data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to create post..")
	}
	json.NewEncoder(w).Encode("Post created successfully..")

}

func (u UserHandler) Comment(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var err error
	var clm *types.Claim
	clm, err = verifyToken(w, r)
	if err != nil {
		json.NewEncoder(w).Encode("Unauthorized token or you are being logged out")
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %s!", clm.Username)))

	var data types.Comment
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to unmarshal json...")
	}

	err = u.user.Coment(clm.Username, data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to save ur comment")
	}
	json.NewEncoder(w).Encode("Comment saved successfully..")

}

func (u UserHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	TokenString = "deleted"
	json.NewEncoder(w).Encode("You are logged out..see u next time")

}

func createToken(username string, userId uint, w http.ResponseWriter) (*types.TokenDetails, error) {
	td := &types.TokenDetails{}
	td.AtExp = time.Now().Add(time.Minute * 15)

	td.RtExp = time.Now().Add(time.Hour * 24 * 7).Unix()

	var err error
	//Creating Access Token
	atClaims := &types.Claim{
		Username: username,
		UserId:   userId,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: td.AtExp.Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rClaims := &types.Claim{
		Username: username,
		UserId:   userId,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: td.RtExp,
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil

}

func verifyToken(w http.ResponseWriter, r *http.Request) (*types.Claim, error) {
	// c, err := r.Cookie("token") //this code is not generating a token stored in the cookie
	// tk := c.Value               //nil,so it creates a panic
	// fmt.Println(tk)

	tknString := TokenString //token generatd is store here to proceed with authentication

	clm := &types.Claim{}

	tkn, err := jwt.ParseWithClaims(tknString, clm, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusBadRequest)

	}
	return clm, nil

}
