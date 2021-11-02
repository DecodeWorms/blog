//Am not getting verification of token right

//Everything seems perfect at createToken() because i passed the generated token to json and i saw it there

func CreateToken(userId uint64, r *http.Request) (*types.TokenDetails, error) {
	var err error

	td := &types.TokenDetails{}
	td.AccessUuid = uuid.NewV4().String()
	//td.UserId = userId
	td.AtExp = time.Now().Add(time.Minute * 15).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["authorization"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExp

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	td.RefreshUuid = uuid.NewV4().String()
	td.RtExp = time.Now().Add(time.Hour * 24 * 7).Unix()

	rfClaims := jwt.MapClaims{}
	rfClaims["refresh_uuid"] = td.RefreshUuid
	rfClaims["exp"] = td.RtExp
	rfClaims["user_id"] = userId

	ft := jwt.NewWithClaims(jwt.SigningMethodHS256, rfClaims)
	td.RefreshToken, err = ft.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil

}

//At saveToken everything seems perfetc because i printed out the access uuid and refersh uuid at the CMD to confirm

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

//At the ExtractToken seems the prof made mistake,dont really sure though because
//1. He split the bear token by empty (space) "" and the result returned set of a single characters and space is btw them
//2. He checked if the length of the split string to be 2 as in == 2 and the legnth is more than 2(since the split string is now set of characters) so its more than 2 and definately it will return empty string ""

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	tokenString := strings.Split(bearToken, "")
	fmt.Println("Bear Token", bearToken)

	if len(tokenString) == 2 {
		return tokenString[1]
	}
	return ""
}

//CORRECTION
//This is how i fixed it

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

    //I used "." to split them and returned header,payload and signature string

	tokenString := strings.Split(bearToken, ".")
	fmt.Println("Bear Token", bearToken)

	if len(tokenString) == 2 {
         // returning payload string here.. and it returned it
		return tokenString[1]
	}
	return ""
}

//At VerifyToken(),exToken variable is confirmed to have token string,BUT after the string token is passed to jwt.Parse() to parse,validate ,return it using ACCESS_SECRET and is RETURNING nil

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	var err error
	exToken := ExtractToken(r)
	token, err := jwt.Parse(exToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid method specified %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("An error %v", err))
	}

	return token, nil

}

//NOTE := KINDLY CHECK THE MAIN SOURCE CODE FOR DETAILS
