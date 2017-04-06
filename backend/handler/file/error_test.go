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

package file

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type failingResponseRecorder struct {
	*httptest.ResponseRecorder
}

func (writer *failingResponseRecorder) Header() http.Header {
	return writer.ResponseRecorder.Header()
}

func (writer *failingResponseRecorder) WriteHeader(code int) {
	writer.ResponseRecorder.WriteHeader(code)
}

func (writer *failingResponseRecorder) Write(buffer []byte) (int, error) {
	writer.ResponseRecorder.Write(buffer)
	return 0, errors.New(`http.ResponseWriter returned an error`)
}

func TestFileError(t *testing.T) {
	t.Parallel()

	t.Run(`test/301/na`, func(t *testing.T) {
		fileError, err := NewError(`test/301/na`)

		if err == nil {
			t.Error(`expected an error, got nil`)
		}

		if (fileError != FileError{}) {
			t.Error(`expected zero value, got `, fileError)
		}
	})

	t.Run(`test/301/invalid`, func(t *testing.T) {
		var fileError FileError

		if !t.Run(`NewError`, func(t *testing.T) {
			var err error

			fileError, err = NewError(`test/301/invalid`)
			if err != nil {
				t.Error(err)
			}
		}) {
			t.FailNow()
		}

		t.Run(`ServeMovedPermanently`, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			fileError.ServeMovedPermanently(recorder, `https://kagucho.net/`)

			if recorder.Code != http.StatusMovedPermanently {
				t.Error(`invalid status code; expected`, http.StatusMovedPermanently,
					`, got `, recorder.Code)
			}

			if result := recorder.HeaderMap.Get(`Content-Language`); result != `` {
				t.Errorf(`invalid Content-Language field in header; expected "" (empty), got %q`,
					result)
			}

			if result := recorder.HeaderMap.Get(`Content-Length`); result != `` {
				t.Errorf(`invalid Content-Length field in header; expected "" (empty), got %q`,
					result)
			}

			if result := recorder.Body.Len(); result != 0 {
				t.Error(`invalid body length; expected 0, got `, result)
			}
		})
	})

	testFileServeError := func(t *testing.T, fileError FileError, code int,
		expected string) {
		t.Run(`WithContentEncoding`, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			recorder.HeaderMap.Set(`Content-Encoding`, `gzip`)

			if !fileError.ServeError(recorder, code) {
				t.Error(`invalid returned value; expected true, got false`)
			}

			if recorder.Code != code {
				t.Error(`invalid status code; expected `, code,
					`, got `, recorder.Code)
			}

			if result := recorder.HeaderMap.Get(`Content-Language`); result != `ja` {
				t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
					result)
			}

			if result := recorder.HeaderMap.Get(`Content-Length`); result != `` {
				t.Errorf(`invalid Content-Length field in header; expected "" (empty), got %q`,
					result)
			}

			assertBodyWithFile(t, recorder, expected)
		})

		t.Run(`WithoutContentEncoding`, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			if !fileError.ServeError(recorder, code) {
				t.Error(`invalid returned value; expected true, got false`)
			}

			if recorder.Code != code {
				t.Error(`invalid status code; expected `, code, `, got `, recorder.Code)
			}

			if result := recorder.HeaderMap.Get(`Content-Language`); result != `ja` {
				t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
					result)
			}

			expectedLen := strconv.Itoa(recorder.Body.Len())
			if result := recorder.HeaderMap.Get(`Content-Length`); result != expectedLen {
				t.Errorf(`invalid Content-Length field in header; expected %q, got %q`,
					expectedLen, result)
			}

			assertBodyWithFile(t, recorder, expected)
		})
	}

	testDirectoryServeError := func(t *testing.T, fileError FileError, code int) {
		t.Run(`WithContentEncoding`, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			recorder.HeaderMap.Set(`Content-Encoding`, `br`)

			if !fileError.ServeError(recorder, code) {
				t.Error(`invalid returned value; expected true, got false`)
			}

			if recorder.Code != code {
				t.Error(`invalid status code; expected `, code, `, got `, recorder.Code)
			}

			if result := recorder.HeaderMap.Get(`Content-Language`); result != `ja` {
				t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
					result)
			}

			if result := recorder.HeaderMap.Get(`Content-Length`); result != `` {
				t.Errorf(`invalid Content-Length field in header; expected "" (empty), got %q`,
					result)
			}
		})

		t.Run(`WithoutContentEncoding`, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			if !fileError.ServeError(recorder, code) {
				t.Error(`invalid returned value; expected true, got false`)
			}

			if recorder.Code != code {
				t.Error(`invalid status code; expected `, code, `, got `, recorder.Code)
			}

			if result := recorder.HeaderMap.Get(`Content-Language`); result != `ja` {
				t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
					result)
			}

			if result := recorder.HeaderMap.Get(`Content-Length`); result == `` {
				t.Errorf(`invalid Content-Length field in header; expected not empty, got %q`,
					result)
			}
		})
	}

	testUnknownFileServeError := func(t *testing.T,
		fileError FileError) {
		testFileServeError(t, fileError, http.StatusBadRequest,
			`test/unknown/file/error/unknown`)
	}

	testUnknownDirectoryServeError := func(t *testing.T, fileError FileError) {
		testDirectoryServeError(t, fileError, http.StatusBadRequest)
	}

	testUnknownNAServeError := func(t *testing.T, fileError FileError) {
		t.Run(`codeFile`, func(t *testing.T) {
			t.Parallel()
			testFileServeError(t, fileError, http.StatusNotFound,
				`test/unknown/na/error/404`)
		})

		t.Run(`codeDirectory`, func(t *testing.T) {
			t.Parallel()
			testDirectoryServeError(t, fileError, http.StatusUnauthorized)
		})

		t.Run(`codeNA`, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			if fileError.ServeError(recorder, http.StatusBadRequest) {
				t.Error(`invalid returned value; expected false, got true`)
			}

			if recorder.Code != http.StatusBadRequest {
				t.Error(`invalid status code; expected `, http.StatusBadRequest,
					`, got `, recorder.Code)
			}

			if result := recorder.HeaderMap.Get(`Content-Language`); result != `` {
				t.Errorf(`invalid Content-Language field in header; expected "" (empty), got %q`,
					result)
			}

			if result := recorder.HeaderMap.Get(`Content-Length`); result != `` {
				t.Errorf(`invalid Content-Length field in header; expected "" (empty), got %q`,
					result)
			}

			if result := recorder.Body.Len(); result != 0 {
				t.Error(`invalid body length; expected 0, got `,
					result)
			}
		})
	}

	for _, unknown := range [...]struct {
		directory      string
		testServeError func(t *testing.T, fileError FileError)
	}{
		{`test/unknown/file`, testUnknownFileServeError},
		{`test/unknown/directory`, testUnknownDirectoryServeError},
		{`test/unknown/na`, testUnknownNAServeError},
	} {
		unknown := unknown

		t.Run(unknown.directory, func(t *testing.T) {
			t.Parallel()

			var fileError FileError

			if !t.Run(`NewError`, func(t *testing.T) {
				var newError error
				fileError, newError = NewError(unknown.directory)
				if newError != nil {
					t.Error(newError)
				}
			}) {
				t.FailNow()
			}

			t.Run(`ServeError`, func(t *testing.T) {
				t.Parallel()
				unknown.testServeError(t, fileError)
			})

			t.Run(`ServeMovedPermanently`, func(t *testing.T) {
				t.Parallel()

				t.Run(`success`, func(t *testing.T) {
					t.Parallel()

					recorder := httptest.NewRecorder()
					fileError.ServeMovedPermanently(recorder, `https://kagucho.net/`)

					if recorder.Code != http.StatusMovedPermanently {
						t.Error(`invalid status code; expected`, http.StatusMovedPermanently,
							`, got `, recorder.Code)
					}

					if result := recorder.HeaderMap.Get(`Content-Language`); result != `ja` {
						t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
							result)
					}

					if resultLen := recorder.HeaderMap.Get(`Content-Length`); resultLen != `21` {
						t.Errorf(`invalid Content-Length field in header; expected "TODO", got %q`,
							resultLen)
					}

					if result := recorder.HeaderMap.Get(`Location`); result != `https://kagucho.net/` {
						t.Errorf(`invalid Location field in header; expected "https://kagucho.net/", got %q`,
							result)
					}

					if result := recorder.Body.String(); result != "https://kagucho.net/\n" {
						t.Errorf(`invalid body; expected "https://kagucho.net/\n", got %q`,
							result)
					}
				})

				t.Run(`failWrite`, func(t *testing.T) {
					t.Parallel()

					recorder := failingResponseRecorder{httptest.NewRecorder()}
					fileError.ServeMovedPermanently(&recorder, `https://kagucho.net/`)

					if recorder.ResponseRecorder.Code != http.StatusMovedPermanently {
						t.Error(`invalid status code; expected`, http.StatusMovedPermanently,
							`, got `, recorder.ResponseRecorder.Code)
					}

					if result := recorder.ResponseRecorder.HeaderMap.Get(`Content-Language`); result != `ja` {
						t.Errorf(`invalid Content-Language field in header; expected "ja", got %q`,
							result)
					}

					if resultLen := recorder.ResponseRecorder.HeaderMap.Get(`Content-Length`); resultLen != `21` {
						t.Errorf(`invalid Content-Length field in header; expected "TODO", got %q`,
							resultLen)
					}

					if result := recorder.ResponseRecorder.HeaderMap.Get(`Location`); result != `https://kagucho.net/` {
						t.Errorf(`invalid Location field in header; expected "https://kagucho.net/", got %q`,
							result)
					}

					if result := recorder.ResponseRecorder.Body.String(); result != "https://kagucho.net/\n" {
						t.Errorf(`invalid body; expected "https://kagucho.net/\n", got %q`,
							result)
					}
				})
			})
		})
	}
}
