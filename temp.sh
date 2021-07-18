#!/bin/bash

Cookie="wr_localvid=xxxx" 

# bookId="CB_5W8ABAA9EDMV6R56Q0"
bookId="36511877"

userAgent='user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:88.0) Gecko/20100101 Firefox/88.0'
header="-H '${userAgent}' -H 'accept: application/json' --cookie '$Cookie'"

# -H 'accept: application/json'
echo  "${header}"
# get bookmarklist
# curl "${header}" "https://weread.qq.com/web/book/bookmarklist?bookId=${bookId}&type=1" -vvv | jq .
curl -H "${userAgent}" --cookie "$Cookie" "https://weread.qq.com/web/book/bookmarklist?bookId=${bookId}&type=1" -vvv
# > ${bookId}.json

# get shelf
# curl -H "${userAgent}" --cookie "$Cookie" 'https://weread.qq.com/web/shelf'  > shelf.html
