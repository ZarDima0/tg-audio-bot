package utils

import "github.com/kkdai/youtube/v2"

func GetAudioFormat(list youtube.FormatList) youtube.Format {
	list.Sort()
	return list[len(list)-1]
}
