#/bin/sh
# Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

set -e

clean() {
  result=$?
  rm -rf $temp
  exit $result
}

trap clean HUP INT QUIT KILL TERM EXIT

temp=`mktemp -d`
echo "mode: count" > $temp/sum.txt
for package in $(go list ./...); do
  go test -covermode=count -coverprofile=$temp/package.txt $package

  if [ -f $temp/package.txt ]; then
    cat $temp/package.txt | tail -n +2 >> $temp/sum.txt
    rm -f $temp/package.txt
  fi
done

go tool cover -func $temp/sum.txt
