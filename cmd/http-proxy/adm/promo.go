package adm

import (
	"encoding/xml"
)

type Content struct {
	XMLName   xml.Name `xml:"Content"`
	DefaultID int      `xml:"defaultId,attr"`
	Ads       []Ad     `xml:"Ad"`
}

type Ad struct {
	ID       int      `xml:"id,attr"`
	Leading  int      `xml:"leading,attr"`
	ShowFor0 int      `xml:"showFor0,attr"`
	ShowFor1 int      `xml:"showFor1,attr"`
	ShowFor2 int      `xml:"showFor2,attr"`
	ShowFor3 int      `xml:"showFor3,attr"`
	ShowFor4 int      `xml:"showFor4,attr"`
	ShowFor5 int      `xml:"showFor5,attr"`
	ShowFor6 int      `xml:"showFor6,attr"`
	ShowFor7 int      `xml:"showFor7,attr"`
	ShowFor8 int      `xml:"showFor8,attr"`
	Header   LangText `xml:"header"`
	Text     LangText `xml:"text"`
	Images   Images   `xml:"images"`
	URLType  LangText `xml:"urlType"`
	URL      LangURL  `xml:"url"`
}

type LangText struct {
	RU string `xml:"ru"`
	EN string `xml:"en"`
}

type Images struct {
	RU []Image `xml:"ru>image"`
	EN []Image `xml:"en>image"`
}

type Image struct {
	URL  string `xml:"url,attr"`
	Hash string `xml:"hash,attr"`
}

type LangURL struct {
	RU URLValue `xml:"ru"`
	EN URLValue `xml:"en"`
}

type URLValue struct {
	URL string `xml:"url,attr"`
}

func PromoXml() (string, error) {
	content := Content{
		DefaultID: 0,
		Ads: []Ad{
			{
				ID:       0,
				Leading:  0,
				ShowFor2: 1,
				ShowFor3: 1,
				ShowFor4: 1,
				Header:   LangText{RU: "Cat", EN: "Cat"},
				Text:     LangText{RU: "Food ", EN: "Food"},
				Images: Images{
					RU: []Image{{
						URL:  "http://adm.star-conflict.com/image/cat.jpg",
						Hash: "4b9334f97f2b1266f4db6250898a82d9",
					}},
					EN: []Image{{
						URL:  "http://adm.star-conflict.com/image/cat.jpg",
						Hash: "4b9334f97f2b1266f4db6250898a82d9",
					}},
				},
				URLType: LangText{RU: "7", EN: "7"},
				URL: LangURL{
					RU: URLValue{URL: "[ship-Ship_Race3_H_T3]"},
					EN: URLValue{URL: "[ship-Ship_Race3_H_T3]"},
				},
			},
		},
	}

	out, err := xml.MarshalIndent(content, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(out), err
}
