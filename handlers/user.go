package handlers

import (
	"blog/storage"
	"blog/types"
	"blog/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

type UserHandler struct {
	user storage.User
}

func NewUserHandler(use storage.User) UserHandler {
	return UserHandler{
		user: use,
	}

}

func (u UserHandler) Table(w http.ResponseWriter, r *http.Request) {
	var data types.User
	var err error
	err = u.user.Table(data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to create a table")
	}
	json.NewEncoder(w).Encode("Tbale created successfully")

}

func (u UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	var data types.User
	var err error

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to decode json or json is empty")
	}
	json.NewEncoder(w).Encode("Fectching json results..")
	err = u.user.Create(data)

}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	util.SetHeader(w)
	var data types.User
	var err error
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to decode json")
	}
	json.NewEncoder(w).Encode("Decoding json..")

	err = u.user.Login(data)

	if err != nil {
		json.NewEncoder(w).Encode("username or password incorrect")
		return
	}
	var tokenDetails *types.TokenDetails
	tokenDetails, err = CreateToken(uint64(data.ID), r)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to create token")
	}
	json.NewEncoder(w).Encode("Successfully login")

	err = saveToken(uint64(data.ID), tokenDetails)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to save token")
	}
	json.NewEncoder(w).Encode("Token saved successfully")
	tokens := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}
	json.NewEncoder(w).Encode(tokens)
	return

}

func (u UserHandler) Details(w http.ResponseWriter, r *http.Request) {
	util.SetHeader(w)
	params := mux.Vars(r)
	userName := params["username"]
	var err error
	var data types.User

	data, err = u.user.Details(userName)
	if err != nil {
		json.NewEncoder(w).Encode("User info not available")
	}
	json.NewEncoder(w).Encode("User infos...")
	json.NewEncoder(w).Encode(data)

	var re *jwt.Token
	re, err = VerifyToken(r)
	fmt.Println(re)
	json.NewEncoder(w).Encode(err)
	json.NewEncoder(w).Encode(re)

}

func CreateToken(userId uint64, r *http.Request) (*types.TokenDetails, error) {
	var err error

	Td := types.TokenDetails{}
	Td.AccessUuid = uuid.NewV4().String()
	//td.UserId = userId
	Td.AtExp = time.Now().Add(time.Minute * 15).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["authorization"] = true
	atClaims["access_uuid"] = Td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = Td.AtExp

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	Td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	Td.RefreshUuid = uuid.NewV4().String()
	Td.RtExp = time.Now().Add(time.Hour * 24 * 7).Unix()

	rfClaims := jwt.MapClaims{}
	rfClaims["refresh_uuid"] = Td.RefreshUuid
	rfClaims["exp"] = Td.RtExp
	rfClaims["user_id"] = userId

	ft := jwt.NewWithClaims(jwt.SigningMethodHS256, rfClaims)
	Td.RefreshToken, err = ft.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return &Td, nil

}

func saveToken(userId uint64, td *types.TokenDetails) error {
	var client *redis.Client

	dsn := os.Getenv("REDIS_PORT")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	res, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	at := time.Unix(td.AtExp, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExp, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	fmt.Println("Access uuid: ", td.AccessUuid, "and", "refresh uuid:", td.RefreshUuid)
	return nil

}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	tokenString := strings.Split(bearToken, ".")
	//fmt.Println("Bear Token", bearToken)

	//code from the source tutorial
	if len(tokenString) >= 2 {
		return tokenString[1]
	}
	return ""
}

// 	if len(tokenString) >= 2 {
// 		return tokenString[1]
// 	}
// 	return ""
// }

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	var err error
	exToken := ExtractToken(r)

	token, _ := jwt.Parse(exToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid method specified %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("An error %v", err))
	}
	fmt.Println("Parsed token", token.Header)
	return token, nil

}

// func VerifyToken(r *http.Request) (*jwt.Token, error) {
// 	tokenString := ExtractToken(r)
// 	fmt.Println("This is the extracted token", tokenString)

// 	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Invalid method specified %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("ACCESS_SECRET")), nil
// 	})

// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	fmt.Println("Parsed token", token)
// 	return token, nil
// }

// func ValidateToken(r *http.Request) error {
// 	token, err := VerifyToken(r)
// 	if err != nil {
// 		return err
// 	}
// 	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
// 		return err
// 	}
// 	return nil
// }

// //No1 to call

// func ExtractTokenMetaData(r *http.Request) (*types.AccessDetails, error) {
// 	var accessUuid string
// 	var err error
// 	var token *jwt.Token
// 	var c *types.AccessDetails

// 	token, err = VerifyToken(r)
// 	fmt.Println("New Token", token)

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	fmt.Println(claims)

// 	if ok && token.Valid {
// 		accessUuid, ok = claims["access_uuid"].(string)

// 		if !ok {
// 			return nil, err
// 		}
// 		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
// 		if err != nil {
// 			return nil, err
// 		}
// 		c = &types.AccessDetails{
// 			AccessUuid: accessUuid,
// 			UserId:     userId,
// 		}
// 		return c, nil
// 	}

// 	return nil, err

// }

// //No2 to call

// func FetchAuth(auth *types.AccessDetails) (uint64, error) {
// 	// var r *redis.Client
// 	// var cli storage.RedisClient
// 	// cli = storage.NewRedisClient(r, os.Getenv("ACCESS_SECRET"))
// 	dsn := os.Getenv("ACCESS_SECRET")
// 	cli := redis.NewClient(&redis.Options{
// 		Addr: dsn, //redis port
// 	})

// 	userid, err := cli.Get(auth.AccessUuid).Result()
// 	if err != nil {
// 		return 0, nil
// 	}
// 	userId, _ := strconv.ParseUint(userid, 10, 64)

// 	return userId, nil
// }
