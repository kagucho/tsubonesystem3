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

package jwt

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math"
	mathRand "math/rand"
	"strconv"
	"time"
)

// Issue takes a subject name and returns a new JWT for him.
func (context JWT) Issue(sub string, scope string,
	duration time.Duration, temporary bool) (string, error) {
	estimateEncodedSize := func(bytes int) int {
		return int(math.Ceil(float64(bytes) * 4 / 3))
	}

	header, marshalError := json.Marshal(header{context.authority.Alg()})
	if marshalError != nil {
		return ``, marshalError
	}

	claimStruct := claim{
		Sub: sub, Scope: scope,
		Jti: strconv.Itoa(mathRand.Int()),
	}

	if duration != 0 {
		claimStruct.Exp = time.Now().Add(duration).Unix()
	}

	if temporary {
		claimStruct.Tmp = true
	}

	claim, marshalError := json.Marshal(claimStruct)
	if marshalError != nil {
		return ``, marshalError
	}

	messageBuffer := bytes.NewBuffer(make([]byte, 0,
		estimateEncodedSize(len(header))+
			1+estimateEncodedSize(len(claim))+
			1+estimateEncodedSize(sha256.Size)))

	headerEncoder := base64.NewEncoder(base64.RawURLEncoding, messageBuffer)
	headerEncoder.Write(header)
	headerEncoder.Close()

	messageBuffer.Write([]byte{'.'})

	claimEncoder := base64.NewEncoder(base64.RawURLEncoding, messageBuffer)
	claimEncoder.Write(claim)
	claimEncoder.Close()

	hash := context.authority.Hash()
	hash.Write(messageBuffer.Bytes())
	jwt := hash.Sum(nil)

	messageBuffer.Write([]byte{'.'})

	jwtEncoder := base64.NewEncoder(base64.RawURLEncoding, messageBuffer)
	jwtEncoder.Write(jwt)
	jwtEncoder.Close()

	return messageBuffer.String(), nil
}
