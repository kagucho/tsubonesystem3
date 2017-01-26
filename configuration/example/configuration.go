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

package configuration

/*
	DBDSN is the string of the DSN which refers to the databse for
	TsuboneSystem.

	It should be understood by go-sql-driver. See the following.
	https://github.com/go-sql-driver/mysql#dsn-data-source-name
*/
const DBDSN string = `root@unix(/var/lib/mysql/mysql.sock)/tsubonesystem`

/*
	DBPssswordKey is the string of the key to encrypt password in the
	database. Its length should be 28 and it must be cryptographically
	random.

	Set the following value for testing.
*/
const DBPasswordKey string = `XXXXXXXXXXXXXXXXXXXXXXXXXXXX`

/*
	FcgiListenNet is the string of the network which tsubonesystem_fcgi
	command should listen to.

	It should be understood by net.Listen; see the following.
	https://golang.org/pkg/net/#Listen.
*/
const FcgiListenNet string = `unix`

/*
	FcgiListenAddress is the string of the address which tsubonesystem3_fcgi
	command should listen to.

	It should be understood by net.Listen; see the following.
	https://golang.org/pkg/net/#Listen.
*/
const FcgiListenAddress string = `/var/lib/tsubonesystem3/tsubonesystem3.sock`

/*
	ListenAddress is the string of the address which tsubonesystem3
	command should listen to.

	It should be understood by net.Listen; see the following.
	https://golang.org/pkg/net/#Listen.
*/
const ListenAddress string = `localhost:8000`
