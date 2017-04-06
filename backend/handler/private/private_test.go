/*
	Copyright (C) 2017  Kagucho <kagucho.net@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or (at
	your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package private

import (
	"bytes"
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/file"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type queryTest struct {
	query    string
	expected graph
}

func TestParseQuery(t *testing.T) {
	for _, test := range [...]struct {
		description string
		routeQuery  []string
		expected    map[string]string
	}{
		{`no2LenRouteQuery`, []string{``}, nil},
		{`no2LenComponent`, []string{``, ``}, map[string]string{}},
		{`unescapeError`, []string{``, `id=%`}, map[string]string{}},
		{`valid`, []string{``, `id=+`}, map[string]string{"id": " "}},
	} {
		t.Run(test.description, func(t *testing.T) {
			if result := parseQuery(test.routeQuery); !reflect.DeepEqual(result, test.expected) {
				t.Error(`expected `, test.expected,
					`, got `, result)
			}
		})
	}
}

func TestPrivate(t *testing.T) {
	db, err := db.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	t.Run(`Private`, func(t *testing.T) {
		fileError, err := file.NewError(`test/valid`)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(`valid`, func(t *testing.T) {
			t.Parallel()

			var private Private

			if !t.Run(`New`, func(t *testing.T) {
				private, err = New(`test/valid`, db, fileError)
				if err != nil {
					t.Error(err)
				}
			}) {
				t.FailNow()
			}

			t.Run(`ServeHTTP`, func(t *testing.T) {
				t.Run(`/private`, func(t *testing.T) {
					t.Parallel()

					t.Run(`unescapedFragment`, func(t *testing.T) {
						t.Parallel()

						t.Run(`withQuery`, func(t *testing.T) {
							t.Parallel()

							recorder := httptest.NewRecorder()
							request := httptest.NewRequest(`GET`,
								`https://kagucho.net/private?key=value`, nil)

							private.ServeHTTP(recorder, request)

							if result := recorder.HeaderMap.Get(`Location`); result != `https://kagucho.net/private` {
								t.Errorf(`invalid Location field in header; expected "http://kagucho.net/private", got %q`,
									result)
							}
						})

						t.Run(`withoutQuery`, func(t *testing.T) {
							t.Parallel()

							recorder := httptest.NewRecorder()
							request := httptest.NewRequest(`GET`,
								`https://kagucho.net/private`, nil)

							private.ServeHTTP(recorder, request)

							body, err := ioutil.ReadFile(private.file)
							if err != nil {
								t.Fatal(err)
							}

							if result := recorder.HeaderMap.Get(`Content-Language`); result != `ja` {
								t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
									result)
							}

							if result := recorder.Body.Bytes(); !bytes.Equal(result, body) {
								t.Errorf(`invalid body; expected %q, got %q`,
									body, result)
							}
						})
					})

					t.Run(`escapedFragment`, func(t *testing.T) {
						t.Parallel()

						for _, test := range [...]struct {
							escapedFragment string
							expected        string
						}{
							{
								``,
								`<!DOCTYPE html>

<html lang="ja">
	<head prefix="og: http://ogp.me/ns#">
		
		<meta content=神楽坂一丁目通信局内で利用しているWebサービスです。 property=og:description>
		
		<meta content=https://kagucho.net/favicon.ico property=og:image>
		
		<meta content=ja_JP property=og:locale>
		
		<meta content=TsuboneSystem property=og:title>
		
		<meta content=website property=og:type>
		
		<meta content=https://kagucho.net/private property=og:url>
		
	</head>
</html>
`,
							}, {
								`club?id=prog`,
								`<!DOCTYPE html>

<html lang="ja">
	<head prefix="og: http://ogp.me/ns#">
		
		<meta content=神楽坂一丁目通信局のProg部の詳細情報です。 property=og:description>
		
		<meta content=https://kagucho.net/favicon.ico property=og:image>
		
		<meta content=ja_JP property=og:locale>
		
		<meta content=TsuboneSystem&#32;Prog部の詳細情報 property=og:title>
		
		<meta content=website property=og:type>
		
		<meta content=https://kagucho.net/private#!club?id&#61;prog property=og:url>
		
	</head>
</html>
`,
							},
						} {
							test := test

							t.Run(test.escapedFragment, func(t *testing.T) {
								recorder := httptest.NewRecorder()
								request := httptest.NewRequest(`GET`,
									`https://kagucho.net/private?_escaped_fragment_=`+
										test.escapedFragment, nil)

								private.ServeHTTP(recorder, request)

								if result := recorder.HeaderMap.Get(`Content-Language`); result != `ja` {
									t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
										result)
								}

								expectedLen := strconv.Itoa(len(test.expected))
								if result := recorder.HeaderMap.Get(`Content-Length`); result != expectedLen {
									t.Errorf(`invalid Content-Length field in header; expected %q, got %q`,
										expectedLen, result)
								}

								if result := recorder.Body.String(); result != test.expected {
									t.Errorf(`invalid body; expected %q, got %q`, test.expected, result)
								}
							})
						}
					})
				})

				t.Run(`/private/`, func(t *testing.T) {
					t.Parallel()

					recorder := httptest.NewRecorder()
					request := httptest.NewRequest(`GET`,
						`https://kagucho.net/private/?key=value`, nil)

					private.ServeHTTP(recorder, request)

					if result := recorder.HeaderMap.Get(`Location`); result != `https://kagucho.net/private` {
						t.Errorf(`invalid Location field in header; expected "https://kagucho.net/private", got %q`,
							result)
					}
				})

				t.Run(`notFound`, func(t *testing.T) {
					t.Parallel()

					recorder := httptest.NewRecorder()
					request := httptest.NewRequest(`GET`, `https://kagucho.net/private!`, nil)

					private.ServeHTTP(recorder, request)

					if recorder.Code != http.StatusNotFound {
						t.Error(`invalid status code; expected `, http.StatusNotFound,
							`, got `, recorder.Code)
					}
				})
			})
		})

		t.Run(`invalid`, func(t *testing.T) {
			t.Parallel()

			var private Private

			if !t.Run(`New`, func(t *testing.T) {
				private, err = New(`test/invalid`, db, fileError)
				if err != nil {
					t.Error(err)
				}
			}) {
				t.FailNow()
			}

			t.Run(`ServeHTTP`, func(t *testing.T) {
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(`GET`,
					`https://kagucho.net/private?_escaped_fragment_=`, nil)

				private.ServeHTTP(recorder, request)

				if recorder.Code != http.StatusInternalServerError {
					t.Error(`invalid status code; expected `,
						http.StatusInternalServerError,
						`, got `, recorder.Code)
				}
			})
		})

		t.Run(`na`, func(t *testing.T) {
			t.Parallel()

			private, err := New(`test/na`, db, fileError)
			if (private != Private{}) {
				t.Error(`expected zero value, got `, private)
			}

			if err == nil {
				t.Error(`expected an error, got nil`)
			}
		})
	})

	t.Run(`graphs`, func(t *testing.T) {
		for _, route := range [...]struct {
			route   string
			queries []queryTest
		}{
			{
				`club`,
				[]queryTest{
					{
						`id=prog`,
						graph{
							[]property{
								/*
									The Open Graph protocol

									Basic Metadata
									http://ogp.me/#metadata
									> og:title - The title of your object as it should appear within the graph, e.g., "The Rock".
									> og:type - The type of your object, e.g., "video.movie". Depending on the type you specify, other properties may also be required.
									> og:image - An image URL which should represent your object within the graph.
									> og:url - The canonical URL of your object that will be used as its permanent ID in the graph, e.g., "http://www.imdb.com/title/tt0117500/".

									Optional Metadata
									http://ogp.me/#optional
									> og:description - A one to two sentence description of your object.
									> og:locale - The locale these tags are marked up in. Of the format language_TERRITORY. Default is en_US.

									http://ogp.me/#type_website
								*/
								{`og:description`, `神楽坂一丁目通信局のProg部の詳細情報です。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem Prog部の詳細情報`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!club?id=prog`},
							}, `https://kagucho.net/private#!club?id=prog`,
						},
					}, {
						``,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の不明な部門です。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 不明な部門です`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!club?id=`},
							}, `https://kagucho.net/private#!club?id=`,
						},
					},
				},
			}, {
				`clubs`,
				[]queryTest{
					{
						``,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の部門一覧です。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 部門一覧`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!clubs`},
							}, `https://kagucho.net/private#!clubs`,
						},
					},
				},
			}, {
				`member`,
				[]queryTest{
					{
						``,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の不明な局員です。`},
								{`og:title`, `TsuboneSystem 不明な局員です`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:type`, `profile`},
								{`og:url`, `https://kagucho.net/private#!member?id=`},
							}, `https://kagucho.net/private#!member?id=`,
						},
					}, {
						`id=1stDisplayID`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員、1 !\%_1"#の詳細情報です。`},
								{`og:title`, `TsuboneSystem 1 !\%_1"#の詳細情報`},
								{`og:profile:username`, `1stDisplayID`},
								{`og:profile:gender`, `male`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:type`, `profile`},
								{`og:url`, `https://kagucho.net/private#!member?id=1stDisplayID`},
							}, `https://kagucho.net/private#!member?id=1stDisplayID`,
						},
					}, {
						`id=2ndDisplayID`,
						graph{
							[]property{
								/*
									http://ogp.me/#type_profile
									> profile:username - string - A short unique string to identify them.
									> profile:gender - enum(male, female) - Their gender.
								*/
								{`og:description`, `神楽坂一丁目通信局の局員、2 !%_1"#の詳細情報です。`},
								{`og:title`, `TsuboneSystem 2 !%_1"#の詳細情報`},
								{`og:profile:username`, `2ndDisplayID`},
								{`og:profile:gender`, `female`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:type`, `profile`},
								{`og:url`, `https://kagucho.net/private#!member?id=2ndDisplayID`},
							}, `https://kagucho.net/private#!member?id=2ndDisplayID`,
						},
					}, {
						`id=3rdDisplayID`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員、3 !\%*1"#の詳細情報です。`},
								{`og:title`, `TsuboneSystem 3 !\%*1"#の詳細情報`},
								{`og:profile:username`, `3rdDisplayID`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:type`, `profile`},
								{`og:url`, `https://kagucho.net/private#!member?id=3rdDisplayID`},
							}, `https://kagucho.net/private#!member?id=3rdDisplayID`,
						},
					},
				},
			}, {
				`members`,
				[]queryTest{
					{
						``,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。7 件`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members`},
							}, `https://kagucho.net/private#!members`,
						},
					}, {
						`entrance=invalid`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。不正な入学年度の指定がありました。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members`},
							}, `https://kagucho.net/private#!members`,
						},
					}, {
						`entrance=0`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。不正な入学年度の指定がありました。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members`},
							}, `https://kagucho.net/private#!members`,
						},
					}, {
						`entrance=2155`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。1 件`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members?entrance=2155`},
							}, `https://kagucho.net/private#!members?entrance=2155`,
						},
					}, {
						`nickname=%25`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。6 件`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members?nickname=%25`},
							}, `https://kagucho.net/private#!members?nickname=%25`,
						},
					}, {
						`ob=0`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。6 件`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members?ob=0`},
							}, `https://kagucho.net/private#!members?ob=0`,
						},
					}, {
						`ob=1`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。1 件`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members?ob=1`},
							}, `https://kagucho.net/private#!members?ob=1`,
						},
					}, {
						`realname=%24%26%5C%25_`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。4 件`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members?realname=%24%26%5C%25_`},
							}, `https://kagucho.net/private#!members?realname=%24%26%5C%25_`,
						},
					}, {
						`entrance=1901&ob=0`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局員検索結果です。5 件`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局員検索`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!members?entrance=1901&ob=0`},
							}, `https://kagucho.net/private#!members?entrance=1901&ob=0`,
						},
					},
				},
			}, {
				`officer`,
				[]queryTest{
					{
						``,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の不明な役職です。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 不明な役職です`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!officer?id=`},
							}, `https://kagucho.net/private#!officer?id=`,
						},
					}, {
						`id=president`,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の局長の詳細情報です。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 局長の詳細情報`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!officer?id=president`},
							}, `https://kagucho.net/private#!officer?id=president`,
						},
					},
				},
			}, {
				`officers`,
				[]queryTest{
					{
						``,
						graph{
							[]property{
								{`og:description`, `神楽坂一丁目通信局の役員一覧です。`},
								{`og:image`, `https://kagucho.net/favicon.ico`},
								{`og:locale`, `ja_JP`},
								{`og:title`, `TsuboneSystem 役員一覧`},
								{`og:type`, `website`},
								{`og:url`, `https://kagucho.net/private#!officers`},
							}, `https://kagucho.net/private#!officers`,
						},
					},
				},
			},
		} {
			route := route

			t.Run(route.route, func(t *testing.T) {
				graphFunc := graphFuncs[route.route]
				for _, query := range route.queries {
					query := query

					t.Run(query.query, func(t *testing.T) {
						t.Parallel()

						result := graphFunc(db, `https://kagucho.net/private`,
							[]string{route.route, query.query})
						if !reflect.DeepEqual(result, query.expected) {
							t.Error(`expected `, query.expected, `, got `, result)
						}
					})
				}
			})
		}

		t.Run(`graphDefault`, func(t *testing.T) {
			t.Parallel()

			expected := graph{
				[]property{
					{`og:description`, `神楽坂一丁目通信局内で利用しているWebサービスです。`},
					{`og:image`, `https://kagucho.net/favicon.ico`},
					{`og:locale`, `ja_JP`},
					{`og:title`, `TsuboneSystem`},
					{`og:type`, `website`},
					{`og:url`, `https://kagucho.net/private`},
				}, `https://kagucho.net/private`,
			}

			result := graphDefault(db, `https://kagucho.net/private`, nil)

			if !reflect.DeepEqual(result, expected) {
				t.Error(`expected `, expected, `, got `, result)
			}
		})

		if err := db.Close(); err != nil {
			t.Fatal(err)
		}

		t.Run(`dbError`, func(t *testing.T) {
			t.Parallel()

			for _, route := range [...]string{
				`club`, `member`, `members`, `officer`,
			} {
				route := route

				t.Run(route, func(t *testing.T) {
					defer func() {
						if recover() == nil {
							t.Error(`expected a panic, which didn't occur`)
						}
					}()

					graphFunc := graphFuncs[route]
					graphFunc(db, `https://kagucho.net/private`, []string{route, ``})
				})
			}
		})
	})
}
