package main

import "flag"

var (
	FlagURL   string
	FlagQuery string
	FlagProxy string
	FlagSrc   string
	FlagDest  string
)

func InitCmdOptions() {
	flag.StringVar(&FlagURL, "url", "", "Specify new url of the library in case it has changed.")
	flag.StringVar(&FlagProxy, "proxy", "", "Specify proxy path.")
	flag.StringVar(&FlagQuery, "q", "", "Query the library.")
	flag.StringVar(&FlagSrc, "src", "", "bk - Books, art - Articles, img - Images, vid - Videos, adbk - Audiobooks.")
	flag.StringVar(&FlagDest, "target", "", "Specify target directory for the downloads.")
	flag.Parse()
}
