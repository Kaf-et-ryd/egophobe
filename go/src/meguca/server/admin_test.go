package server

import (
	"bytes"
	"fmt"
	"io"
	"meguca/auth"
	"meguca/common"
	"meguca/config"
	"meguca/db"
	. "meguca/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

var adminLoginCreds = auth.SessionCreds{
	UserID:  "admin",
	Session: genSession(),
}

func newJSONPair(t *testing.T, url string, data interface{}) (
	*httptest.ResponseRecorder, *http.Request,
) {
	t.Helper()

	body := encodeBody(t, data)
	return httptest.NewRecorder(), httptest.NewRequest("POST", url, body)
}

func encodeBody(t *testing.T, data interface{}) io.Reader {
	t.Helper()
	return bytes.NewReader(marshalJSON(t, data))
}

func TestNotBoardOwner(t *testing.T) {
	assertTableClear(t, "accounts", "boards")
	writeSampleBoard(t)
	writeSampleUser(t)

	paths := [...]string{
		"/api/configure-board/a",
		"/html/configure-board/a",
	}
	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			t.Parallel()

			rec, req := newJSONPair(t, p, map[string]string{})
			router.ServeHTTP(rec, req)
			assertError(t, rec, 403, errAccessDenied)
		})
	}
}

func TestServePrivateBoardConfigs(t *testing.T) {
	assertTableClear(t, "boards", "accounts")

	config.ClearBoards()
	conf := db.BoardConfigs{
		BoardConfigs: config.BoardConfigs{
			ID:        "a",
			Eightball: []string{"a", "b", "c"},
		},
	}
	_, err := config.SetBoardConfigs(conf.BoardConfigs)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.WriteBoard(nil, conf); err != nil {
		t.Fatal(err)
	}
	writeSampleUser(t)
	writeSampleBoardOwner(t)

	rec, req := newJSONPair(t, "/api/board-config/a", nil)
	setLoginCookies(req, sampleLoginCreds)
	router.ServeHTTP(rec, req)
	assertCode(t, rec, 200)
	assertBody(t, rec, string(marshalJSON(t, conf.BoardConfigs)))
}

func TestBoardConfiguration(t *testing.T) {
	assertTableClear(t, "accounts", "boards")

	const board = "a"
	conf := config.BoardConfigs{
		ID:        board,
		Eightball: []string{},
		BoardPublic: config.BoardPublic{
			ForcedAnon: true,
			DefaultCSS: "egophobe",
		},
	}
	init := db.BoardConfigs{
		BoardConfigs: config.BoardConfigs{
			BoardPublic: config.BoardPublic{
				DefaultCSS: "egophobe",
			},
			ID:        board,
			Eightball: []string{},
		},
	}
	if err := db.WriteBoard(nil, init); err != nil {
		t.Fatal(err)
	}

	writeSampleUser(t)
	writeSampleBoardOwner(t)

	data := boardConfigSettingRequest{
		BoardConfigs: conf,
	}
	rec, req := newJSONPair(t, "/api/configure-board/a", data)
	setLoginCookies(req, sampleLoginCreds)
	router.ServeHTTP(rec, req)
	assertCode(t, rec, 200)

	res, err := db.GetBoardConfigs(board)
	if err != nil {
		t.Fatal(err)
	}
	AssertDeepEquals(t, res, conf)
}

func TestValidateBoardConfigs(t *testing.T) {
	t.Parallel()

	cases := [...]struct {
		name string
		config.BoardConfigs
		err error
	}{
		{
			"all is well",
			config.BoardConfigs{
				BoardPublic: config.BoardPublic{
					DefaultCSS: "egophobe",
				},
			},
			nil,
		},
		{
			"too many eightball answers",
			config.BoardConfigs{
				Eightball: make([]string, maxEightballLen+1),
			},
			errTooManyAnswers,
		},
		{
			"compound eightball length to big",
			config.BoardConfigs{
				Eightball: []string{GenString(maxEightballLen + 1)},
			},
			errEightballTooLong,
		},
		{
			"notice too long",
			config.BoardConfigs{
				BoardPublic: config.BoardPublic{
					Notice: GenString(common.MaxLenNotice + 1),
				},
			},
			errNoticeTooLong,
		},
		{
			"rules too long",
			config.BoardConfigs{
				BoardPublic: config.BoardPublic{
					Rules: GenString(common.MaxLenRules + 1),
				},
			},
			errRulesTooLong,
		},
		{
			"title too long",
			config.BoardConfigs{
				BoardPublic: config.BoardPublic{
					Title: GenString(common.MaxLenBoardTitle + 1),
				},
			},
			errTitleTooLong,
		},
	}

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			if b := validateBoardConfigs(rec, c.BoardConfigs); b != (c.err == nil) {
				t.Fatal("unexpected result")
			}
			if c.err != nil {
				assertCode(t, rec, 400)
				assertBody(t, rec, fmt.Sprintf("400 %s\n", c.err))
			}
		})
	}
}

func TestValidateBoardCreation(t *testing.T) {
	assertTableClear(t, "boards", "accounts")
	writeSampleBoard(t)
	writeSampleUser(t)

	cases := [...]struct {
		name, id, title string
		err             error
	}{
		{
			name:  "board name too long",
			id:    GenString(common.MaxLenBoardID + 1),
			title: "foo",
			err:   errInvalidBoardName,
		},
		{
			name:  "empty board name",
			id:    "",
			title: "foo",
			err:   errInvalidBoardName,
		},
		{
			name:  "invalid chars in board name",
			id:    ":^)",
			title: "foo",
			err:   errInvalidBoardName,
		},
		{
			name:  "title too long",
			id:    "b",
			title: GenString(101),
			err:   errTitleTooLong,
		},
		{
			name:  "board name taken",
			id:    "a",
			title: "foo",
			err:   errBoardNameTaken,
		},
	}

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			msg := boardCreationRequest{
				ID:    c.id,
				Title: c.title,
			}
			rec, req := newJSONPair(t, "/api/create-board", msg)
			setLoginCookies(req, sampleLoginCreds)
			router.ServeHTTP(rec, req)

			assertCode(t, rec, 400)
			assertBody(t, rec, fmt.Sprintf("400 %s\n", c.err))
		})
	}
}

func writeSampleBoard(t testing.TB) {
	t.Helper()

	b := db.BoardConfigs{
		BoardConfigs: config.BoardConfigs{
			ID:        "a",
			Eightball: []string{"yes"},
		},
	}
	if err := db.WriteBoard(nil, b); err != nil {
		t.Fatal(err)
	}
	if _, err := config.SetBoardConfigs(b.BoardConfigs); err != nil {
		t.Fatal(err)
	}
}

func writeSampleBoardOwner(t *testing.T) {
	t.Helper()

	tx, err := db.StartTransaction()
	if err != nil {
		t.Fatal(err)
	}
	defer db.RollbackOnError(tx, &err)

	err = db.WriteStaff(tx, "a", map[string][]string{
		"owners": {"user1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestBoardCreation(t *testing.T) {
	assertTableClear(t, "boards", "accounts")
	writeSampleUser(t)

	const (
		id    = "a"
		title = "/a/ - Animu & Mango"
	)

	msg := boardCreationRequest{
		ID:    id,
		Title: title,
	}
	rec, req := newJSONPair(t, "/api/create-board", msg)
	setLoginCookies(req, sampleLoginCreds)
	router.ServeHTTP(rec, req)

	assertCode(t, rec, 200)

	board, err := db.GetBoardConfigs(id)
	if err != nil {
		t.Fatal(err)
	}

	std := config.BoardConfigs{
		ID: id,
		BoardPublic: config.BoardPublic{
			Title: title,
		},
		Eightball: config.EightballDefaults,
	}
	AssertDeepEquals(t, board, std)
}

func TestServePrivateServerConfigs(t *testing.T) {
	assertTableClear(t, "accounts")
	writeSampleUser(t)
	writeAdminAccount(t)
	if err := config.Set(config.Defaults); err != nil {
		t.Fatal(err)
	}

	cases := [...]struct {
		name string
		auth.SessionCreds
		code int
		err  error
	}{
		{
			name:         "not admin",
			SessionCreds: sampleLoginCreds,
			code:         403,
			err:          errAccessDenied,
		},
		{
			name:         "admin",
			SessionCreds: adminLoginCreds,
			code:         200,
		},
	}

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			rec, req := newJSONPair(t, "/api/config", nil)
			setLoginCookies(req, c.SessionCreds)
			router.ServeHTTP(rec, req)

			assertCode(t, rec, c.code)
			if c.err != nil {
				assertBody(t, rec, fmt.Sprintf("%d %s\n", c.code, c.err))
			}
		})
	}
}

func writeAdminAccount(t *testing.T) {
	t.Helper()

	err := db.CreateAdminAccount()
	if err != nil {
		t.Fatal(err)
	}
	err = db.WriteLoginSession("admin", adminLoginCreds.Session)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerConfigSetting(t *testing.T) {
	assertTableClear(t, "accounts")
	if err := db.WriteConfigs(config.Defaults); err != nil {
		t.Fatal(err)
	}
	writeAdminAccount(t)

	msg := config.Defaults
	msg.DefaultCSS = "egophobe"
	rec, req := newJSONPair(t, "/api/configure-server", msg)
	setLoginCookies(req, adminLoginCreds)
	router.ServeHTTP(rec, req)

	assertCode(t, rec, 200)

	conf, err := db.GetConfigs()
	if err != nil {
		t.Fatal(err)
	}
	std := config.Defaults
	std.DefaultCSS = "egophobe"
	AssertDeepEquals(t, conf, std)
}

func TestDeleteBoard(t *testing.T) {
	assertTableClear(t, "accounts", "boards")
	writeSampleUser(t)
	writeSampleBoard(t)
	writeSampleBoardOwner(t)

	rec, req := newJSONPair(t, "/api/delete-board", boardActionRequest{
		Board: "a",
	})
	setLoginCookies(req, sampleLoginCreds)
	router.ServeHTTP(rec, req)

	assertCode(t, rec, 200)
}

func TestDeletePost(t *testing.T) {
	assertTableClear(t, "accounts", "boards")
	writeSampleBoard(t)
	writeSampleThread(t)
	writeSampleUser(t)
	writeSampleBoardOwner(t)

	cConfigs := db.BoardConfigs{
		BoardConfigs: config.BoardConfigs{
			ID:        "c",
			Eightball: []string{"yes"},
		},
	}
	err := db.WriteBoard(nil, cConfigs)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := config.SetBoardConfigs(cConfigs.BoardConfigs); err != nil {
		t.Fatal(err)
	}

	thread := db.Thread{
		ID:    3,
		Board: "c",
	}
	op := db.Post{
		StandalonePost: common.StandalonePost{
			Board: "c",
			Post: common.Post{
				ID: 3,
			},
			OP: 3,
		},
	}
	if err := db.WriteThread(nil, thread, op); err != nil {
		t.Fatal(err)
	}

	posts := [...]db.Post{
		{
			StandalonePost: common.StandalonePost{
				Board: "a",
				Post: common.Post{
					ID: 2,
				},
				OP: 1,
			},
		},
		{
			StandalonePost: common.StandalonePost{
				Board: "a",
				Post: common.Post{
					ID: 4,
				},
				OP: 1,
			},
		},
	}
	for _, p := range posts {
		if err := db.WritePost(nil, p); err != nil {
			t.Fatal(err)
		}
	}

	data := []uint64{2, 4}
	const url = "/api/delete-post"
	rec, req := newJSONPair(t, url, data)
	setLoginCookies(req, sampleLoginCreds)
	router.ServeHTTP(rec, req)
	assertCode(t, rec, 200)

	data = []uint64{3}
	rec, req = newJSONPair(t, url, data)
	router.ServeHTTP(rec, req)
	assertCode(t, rec, 403)

	cases := [...]struct {
		name    string
		id      uint64
		deleted bool
	}{
		{"from target board", 2, true},
		{"from target board", 4, true},
		{"different board", 3, false},
	}

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			post, err := db.GetPost(c.id)
			switch {
			case err != nil:
				t.Fatal(err)
			case post.Deleted != c.deleted:
				LogUnexpected(t, post.Deleted, c.deleted)
			}
		})
	}
}

func writeSampleThread(t *testing.T) {
	t.Helper()

	thread := db.Thread{
		ID:        1,
		Board:     "a",
		ReplyTime: 11,
	}
	op := db.Post{
		StandalonePost: common.StandalonePost{
			Post: common.Post{
				ID:   1,
				Time: 345351,
			},
			OP:    1,
			Board: "a",
		},
	}
	if err := db.WriteThread(nil, thread, op); err != nil {
		t.Fatal(err)
	}
}
