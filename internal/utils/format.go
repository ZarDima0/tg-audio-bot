package utils

import yt "github.com/kkdai/youtube/v2"

func GetAudioFormat(list yt.FormatList) yt.Format {
	list.Sort()
	return list[len(list)-1]
}
