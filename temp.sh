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

# links
Charles 抓包iPhone，获取 微信读书api 接口
https://i.weread.qq.com/shelf/sync
https://weread.qq.com/web/shelf/bookIds

https://i.weread.qq.com/related/bookInfo?bookId=26271721
https://i.weread.qq.com/book/info?bookId=26271721

https://gist.github.com/ericclemmons/b146fe5da72ca1f706b2ef72a20ac39d

书架
https://i.weread.qq.com/shelf/sync

热门划线
https://i.weread.qq.com/book/bestbookmarks?bookId=30638159&count=5000
https://i.weread.qq.com/book/bestbookmarks?bookId=25016201&count=5000

https://i.weread.qq.com/market/category?categoryId=1500000&count=40&gender=1&maxIdx=40&rn=1&teenmode=0

书单榜
https://i.weread.qq.com/market/booklists?count=30&maxIdx=90&rn=1&type=0
https://i.weread.qq.com/booklist/single?booklistId=87560503_77jee5fdK
https://i.weread.qq.com/market/booklists?count=200&type=0
https://i.weread.qq.com/booklist/single?booklistId=3607396_6Z0o0WYhJ

获取书单榜100个书单，书籍的热门划线，
https://i.weread.qq.com/book/bestbookmarks?bookId=917691&count=500
https://weread.qq.com/web/reader/d87320605e00bbd874a12ae
